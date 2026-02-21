package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"

	"github.com/lancekrogers/agent-coordinator-ethden-2026/internal/config"
	"github.com/lancekrogers/agent-coordinator-ethden-2026/internal/coordinator"
	"github.com/lancekrogers/agent-coordinator-ethden-2026/internal/hedera/hcs"
	"github.com/lancekrogers/agent-coordinator-ethden-2026/internal/hedera/hts"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))

	cfg, err := config.Load()
	if err != nil {
		log.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	if err := cfg.Coordinator.Validate(); err != nil {
		log.Error("invalid coordinator config", "error", err)
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Create Hedera client with coordinator credentials.
	hederaClient := hiero.ClientForTestnet()
	hederaClient.SetOperator(cfg.CoordinatorAccountID, cfg.CoordinatorKey)

	// Initialize HCS publisher and subscriber.
	publisher := hcs.NewPublisher(hederaClient, hcs.DefaultPublishConfig())
	subscriber := hcs.NewSubscriber(hederaClient, hcs.DefaultSubscribeConfig())

	// Initialize HTS transfer service.
	transferSvc := hts.NewTransferService(hederaClient)

	// Create coordinator components.
	agentIDs := []string{"inference-001", "defi-001"}
	assigner := coordinator.NewAssigner(publisher, cfg.Coordinator.TaskTopicID, agentIDs)
	monitor := coordinator.NewMonitor(subscriber, cfg.Coordinator.StatusTopicID, nil)
	payment := coordinator.NewPayment(transferSvc, publisher, cfg.Coordinator)

	// Agent ID → Hedera account ID for payments.
	agentAccounts := map[string]string{
		"inference-001": cfg.Agent1AccountID,
		"defi-001":      cfg.Agent2AccountID,
	}

	resultHandler := coordinator.NewResultHandler(coordinator.ResultHandlerConfig{
		Subscriber:    subscriber,
		TopicID:       cfg.Coordinator.StatusTopicID,
		Payment:       payment,
		Config:        cfg.Coordinator,
		Log:           log,
		AgentAccounts: agentAccounts,
	})

	// Start monitor and result handler in background.
	go func() {
		if err := monitor.Start(ctx); err != nil {
			log.Error("monitor stopped", "error", err)
		}
	}()
	go func() {
		if err := resultHandler.Start(ctx); err != nil {
			log.Error("result handler stopped", "error", err)
		}
	}()

	// Build and execute integration plan.
	plan := coordinator.IntegrationCyclePlan("inference-001", "defi-001")
	log.Info("coordinator starting",
		"version", "0.2.0",
		"task_topic", cfg.Coordinator.TaskTopicID,
		"status_topic", cfg.Coordinator.StatusTopicID,
		"tasks", plan.TaskCount())

	assignedIDs, err := assigner.AssignTasks(ctx, plan)
	if err != nil {
		log.Error("failed to assign tasks", "error", err)
		os.Exit(1)
	}
	log.Info("tasks assigned", "count", len(assignedIDs), "task_ids", assignedIDs)

	// Block until shutdown signal.
	<-ctx.Done()
	log.Info("coordinator shutting down")
}
