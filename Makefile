.PHONY: proto
proto:
	@protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	./pkg/proto/comm_agent.proto

.PHONY: clean
clean:
	rm -r ./tmp && mkdir ./tmp

.PHONY: run-agent
run-agent: clean build-agent
	@export CONNECTION_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJydW5uZXJJZCI6N30.oLni-njsf32Rcc_O8erONKsHYrHvnuhGBrJO-KN2cK0 && \
	export SERVER_ADDRESS=localhost:8091 && \
	./bin/agent

.PHONY: build-agent
build-agent:
	@go build -o ./bin/agent ./cmd/agent/agent.go

.PHONY: run-controlplane
run-controlplane: build-controlplane
	@export CONN_STRING=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable && \
	./bin/controlplane

.PHONY: build-controlplane
build-controlplane:
	@go build -o ./bin/controlplane ./cmd/controlplane/controlplane.go