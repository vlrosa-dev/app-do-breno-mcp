.DEFAULT_GOAL := help

BIN_DIR   := bin
MCP_BIN   := $(BIN_DIR)/breno-mcp
MCP_MAIN  := ./cmd/breno-mcp

.PHONY: help
help: ## Lista os targets disponíveis
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Compila o servidor MCP (breno-mcp)
	CGO_ENABLED=0 go build -o $(MCP_BIN) $(MCP_MAIN)

.PHONY: test
test: ## Roda todos os testes
	go test ./... -count=1

.PHONY: test-integration
test-integration: ## Testa as rotas reais da API (requer API_USERNAME/API_PASSWORD; escrita com BRENO_IT_WRITE=1)
	go test -tags integration ./internal/integration/ -count=1 -v

.PHONY: test-race
test-race: ## Roda os testes com race detector (requer cgo)
	CGO_ENABLED=1 go test ./... -race -count=1

.PHONY: cover
cover: ## Roda os testes gerando perfil de cobertura
	go test ./... -covermode=atomic -coverprofile=coverage.out
	go tool cover -func=coverage.out | tail -n 1

.PHONY: lint
lint: ## Roda go vet (e golangci-lint se instalado)
	go vet ./...
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run ./... || \
		echo "golangci-lint não instalado — pulando"

.PHONY: tidy
tidy: ## go mod tidy
	go mod tidy

.PHONY: run-mcp
run-mcp: ## Executa o servidor MCP localmente (stdio)
	go run $(MCP_MAIN)

.PHONY: clean
clean: ## Remove artefatos de build
	rm -rf $(BIN_DIR) coverage.out
