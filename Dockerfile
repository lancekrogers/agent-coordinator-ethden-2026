FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /src/bin/coordinator ./cmd/coordinator

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /src/bin/coordinator /usr/local/bin/coordinator
ENTRYPOINT ["coordinator"]
