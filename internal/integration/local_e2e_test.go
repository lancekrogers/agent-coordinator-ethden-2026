package integration

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"
)

// envelope mirrors the shared Envelope format used by inference and DeFi agents.
type envelope struct {
	Type        string          `json:"type"`
	Sender      string          `json:"sender"`
	Recipient   string          `json:"recipient,omitempty"`
	SequenceNum uint64          `json:"sequence_num"`
	Timestamp   time.Time       `json:"timestamp"`
	Payload     json.RawMessage `json:"payload"`
}

// memTransport is a thread-safe in-memory pub/sub that simulates HCS topics.
type memTransport struct {
	mu   sync.RWMutex
	subs map[string][]chan []byte // topicID -> subscriber channels
}

func newMemTransport() *memTransport {
	return &memTransport{subs: make(map[string][]chan []byte)}
}

func (m *memTransport) publish(_ context.Context, topicID string, data []byte) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, ch := range m.subs[topicID] {
		select {
		case ch <- append([]byte(nil), data...): // copy to avoid races
		default:
		}
	}
}

func (m *memTransport) subscribe(_ context.Context, topicID string) <-chan []byte {
	ch := make(chan []byte, 16)
	m.mu.Lock()
	m.subs[topicID] = append(m.subs[topicID], ch)
	m.mu.Unlock()
	return ch
}

// marshalEnvelope is a helper that builds a JSON envelope.
func marshalEnvelope(t *testing.T, msgType, sender, recipient string, seq uint64, payload any) []byte {
	t.Helper()
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	env := envelope{
		Type:        msgType,
		Sender:      sender,
		Recipient:   recipient,
		SequenceNum: seq,
		Timestamp:   time.Now().UTC(),
		Payload:     raw,
	}
	data, err := json.Marshal(env)
	if err != nil {
		t.Fatalf("marshal envelope: %v", err)
	}
	return data
}

// receiveEnvelope reads one envelope from ch within the deadline.
func receiveEnvelope(t *testing.T, ch <-chan []byte, timeout time.Duration) envelope {
	t.Helper()
	select {
	case raw := <-ch:
		var env envelope
		if err := json.Unmarshal(raw, &env); err != nil {
			t.Fatalf("unmarshal envelope: %v", err)
		}
		return env
	case <-time.After(timeout):
		t.Fatal("timed out waiting for envelope")
		return envelope{} // unreachable
	}
}

func TestLocalThreeAgentCycle(t *testing.T) {
	const (
		taskTopic   = "0.0.1001" // simulated HCS task topic
		resultTopic = "0.0.1002" // simulated HCS result topic

		coordinatorID = "coordinator-0.0.100"
		inferenceID   = "inference-0.0.200"
		defiID        = "defi-0.0.300"

		taskID   = "task-abc-001"
		deadline = 3 * time.Second
	)

	transport := newMemTransport()

	// ── Phase 1: Set up subscribers before any publishes ──────────────
	t.Log("Phase 1: Subscribing to topics")

	ctx := context.Background()

	// Inference agent subscribes to the task topic.
	inferenceTaskCh := transport.subscribe(ctx, taskTopic)

	// DeFi agent subscribes to the task topic.
	defiTaskCh := transport.subscribe(ctx, taskTopic)

	// Coordinator subscribes to the result topic for responses.
	coordResultCh := transport.subscribe(ctx, resultTopic)

	// ── Phase 2: Coordinator publishes a task_assignment ──────────────
	t.Log("Phase 2: Coordinator assigns task to inference agent")

	taskPayload := map[string]string{
		"task_id":   taskID,
		"task_type": "market_sentiment",
		"recipient": inferenceID,
	}
	assignMsg := marshalEnvelope(t, "task_assignment", coordinatorID, inferenceID, 1, taskPayload)
	transport.publish(ctx, taskTopic, assignMsg)

	// ── Phase 3: Inference agent receives the task ───────────────────
	t.Log("Phase 3: Inference agent receives task_assignment")

	infEnv := receiveEnvelope(t, inferenceTaskCh, deadline)
	if infEnv.Type != "task_assignment" {
		t.Fatalf("inference: expected type task_assignment, got %s", infEnv.Type)
	}
	if infEnv.Sender != coordinatorID {
		t.Fatalf("inference: expected sender %s, got %s", coordinatorID, infEnv.Sender)
	}

	var assignedTask map[string]string
	if err := json.Unmarshal(infEnv.Payload, &assignedTask); err != nil {
		t.Fatalf("inference: unmarshal payload: %v", err)
	}
	if assignedTask["task_id"] != taskID {
		t.Fatalf("inference: expected task_id %s, got %s", taskID, assignedTask["task_id"])
	}
	t.Logf("  Inference received task: id=%s type=%s", assignedTask["task_id"], assignedTask["task_type"])

	// DeFi agent also receives the broadcast (it will ignore tasks not for it).
	defiEnv := receiveEnvelope(t, defiTaskCh, deadline)
	if defiEnv.Type != "task_assignment" {
		t.Fatalf("defi: expected type task_assignment, got %s", defiEnv.Type)
	}
	t.Log("  DeFi agent also received broadcast (will filter by recipient)")

	// ── Phase 4: Inference agent publishes task_result ────────────────
	t.Log("Phase 4: Inference agent publishes task_result")

	resultPayload := map[string]any{
		"task_id":     taskID,
		"status":      "completed",
		"tx_hash":     "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
		"duration_ms": 150,
	}
	resultMsg := marshalEnvelope(t, "task_result", inferenceID, coordinatorID, 2, resultPayload)
	transport.publish(ctx, resultTopic, resultMsg)

	// ── Phase 5: DeFi agent publishes pnl_report ─────────────────────
	t.Log("Phase 5: DeFi agent publishes pnl_report")

	pnlPayload := map[string]any{
		"agent_id":          defiID,
		"net_pnl":           1.5,
		"trade_count":       3,
		"is_self_sustaining": true,
	}
	pnlMsg := marshalEnvelope(t, "pnl_report", defiID, coordinatorID, 3, pnlPayload)
	transport.publish(ctx, resultTopic, pnlMsg)

	// ── Phase 6: Coordinator receives both results ───────────────────
	t.Log("Phase 6: Coordinator verifies receipt of both responses")

	received := make(map[string]envelope, 2)
	for range 2 {
		env := receiveEnvelope(t, coordResultCh, deadline)
		received[env.Type] = env
		t.Logf("  Coordinator received: type=%s sender=%s seq=%d", env.Type, env.Sender, env.SequenceNum)
	}

	// Verify task_result from inference agent.
	tr, ok := received["task_result"]
	if !ok {
		t.Fatal("coordinator never received task_result")
	}
	if tr.Sender != inferenceID {
		t.Fatalf("task_result sender: expected %s, got %s", inferenceID, tr.Sender)
	}
	var trPayload map[string]any
	if err := json.Unmarshal(tr.Payload, &trPayload); err != nil {
		t.Fatalf("unmarshal task_result payload: %v", err)
	}
	if trPayload["status"] != "completed" {
		t.Fatalf("task_result status: expected completed, got %v", trPayload["status"])
	}

	// Verify pnl_report from DeFi agent.
	pnl, ok := received["pnl_report"]
	if !ok {
		t.Fatal("coordinator never received pnl_report")
	}
	if pnl.Sender != defiID {
		t.Fatalf("pnl_report sender: expected %s, got %s", defiID, pnl.Sender)
	}
	var pnlResult map[string]any
	if err := json.Unmarshal(pnl.Payload, &pnlResult); err != nil {
		t.Fatalf("unmarshal pnl_report payload: %v", err)
	}
	if pnlResult["agent_id"] != defiID {
		t.Fatalf("pnl_report agent_id: expected %s, got %v", defiID, pnlResult["agent_id"])
	}
	netPnL, _ := pnlResult["net_pnl"].(float64)
	if netPnL != 1.5 {
		t.Fatalf("pnl_report net_pnl: expected 1.5, got %v", netPnL)
	}
	tradeCount, _ := pnlResult["trade_count"].(float64)
	if int(tradeCount) != 3 {
		t.Fatalf("pnl_report trade_count: expected 3, got %v", tradeCount)
	}
	selfSustaining, _ := pnlResult["is_self_sustaining"].(bool)
	if !selfSustaining {
		t.Fatal("pnl_report is_self_sustaining: expected true")
	}

	t.Log("All phases passed: coordinator -> assign -> inference result -> DeFi P&L -> coordinator verified")
}
