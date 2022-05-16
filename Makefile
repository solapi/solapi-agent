tidy:
	go mod tidy

download:
	go mod download

build:
	GOOS=freebsd GOARCH=386 go build -o build/agent-freebsd-386 ./cmd/agent/agent.go
	GOOS=freebsd GOARCH=amd64 go build -o build/agent-freebsd-amd64 ./cmd/agent/agent.go
	GOOS=linux GOARCH=386 go build -o build/agent-linux-386 ./cmd/agent/agent.go
	GOOS=linux GOARCH=amd64 go build -o build/agent-linux-amd64 ./cmd/agent/agent.go

all: tidy download build