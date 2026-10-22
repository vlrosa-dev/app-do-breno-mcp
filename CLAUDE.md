# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## O que é

Servidor MCP (stdio) em Go para o App do Breno (finanças pessoais sobre Supabase). Binário único: `breno-mcp` (`cmd/breno-mcp`). Código, comentários, mensagens de erro e commits são em português.

## Comandos

```bash
make build              # bin/breno-mcp (CGO_ENABLED=0)
make test               # go test ./... -count=1 (unitários, sem rede)
make test-race          # com race detector (requer cgo)
make lint               # go vet (+ golangci-lint se instalado)
make run-mcp            # roda o servidor via stdio
make test-integration   # rotas reais da API (tag `integration`)

# um teste específico
go test ./internal/breno -run TestCriarMovimentacaoUnica -count=1 -v
go test -tags integration ./internal/integration -run TestRotasDeConsulta -count=1 -v
```

- Go `1.26.2` (ver `go.mod`).
- CI (`.github/workflows/ci.yml`): vet + testes com `-race` + golangci-lint em push na `main` e PRs. Release via GoReleaser ao empurrar tag `v*`.
- Testes de integração precisam de `API_USERNAME`/`API_PASSWORD` (ambiente ou `.env` na raiz) e pulam sem eles. Os de **escrita** só rodam com `BRENO_IT_WRITE=1` e alteram dados reais da conta (criam movimentações de R$ 0,01 marcadas `[it-mcp]` e as apagam) — não rodar sem o usuário pedir.

## Arquitetura

Fluxo de uma chamada: tool MCP → `internal/mcpserver` → `breno.Financas` → `internal/api` → Supabase.

- `cmd/breno-mcp/main.go` — composição: `.env` (opcional) → `config.Load()` → `auth.Manager` → `api.Client` → `breno.Service` → `mcpserver.NewMcpServer`. **stdout é do protocolo MCP; todo log vai para stderr.**
- `internal/mcpserver` — registra as 17 tools (`server.go`) e os handlers (`tools.go`). Handlers só mapeiam entrada → domínio; nenhuma regra de negócio aqui. Erros do serviço viram resultado em texto (`"Unable to ...: <err>"`), não erro de protocolo. Listas são embrulhadas em `{"items": [...]}` no structuredContent (o MCP exige objeto).
- `internal/breno` — todas as regras de negócio e validação (`validate.go`: normaliza tipo/frequência, datas `YYYY-MM-DD`). Interface `Financas` em `contract.go`; `Service` depende de interfaces locais (`httpDoer`, `Identity`) para permitir fakes nos testes. `CriarMovimentacao` cria movimentação única ou série (`Recorrencia != nil`). Read models são aliases de `internal/structs`.
- `internal/api` — `Client.REST` (PostgREST `/rest/v1/...`) e `Client.RPC` (`POST /rest/v1/rpc/<fn>`), `APIError` para status ≥ 300. `transport.go` injeta `apikey` + `Bearer` e, num 401, força relogin e repete **uma vez** (só se o corpo for rebobinável).
- `internal/auth` — login sempre por senha (`grant_type=password`); `refresh_token` é ignorado de propósito. Token + credenciais persistidos em `~/.breno/config.json` (0600) e reaproveitados até expirar. `gotrue.go` usa um `http.Client` próprio, nunca o `api.Client`, para evitar recursão no login. `Manager.EnsureUserID` fornece o `user_id` usado em toda chamada.
- `internal/structs` — structs de wire (campos `p_*` das RPCs) que **também** são os tipos de entrada das tools MCP, ou seja, definem o schema publicado.
- `internal/config` — `API_USERNAME`, `API_PASSWORD`; `API_BASE_URL` e `API_KEY` têm default embutido (chave anon pública).

## Pontos de atenção

- **Schema das tools**: o SDK gera o schema a partir dos structs; campo sem `omitempty` vira obrigatório. `inputSchema`/`camposOpcionais` em `internal/mcpserver/server.go` tiram da lista de obrigatórios os campos que o servidor preenche ou ignora (ids de usuário, `Order`, `Select`, `Limit`). Ao adicionar campo de usuário a um struct de entrada, inclua-o lá. Nomes e schemas das tools são consumidos por agentes externos — mudanças quebram integrações.
- **Nunca confiar no `user_id` vindo da entrada da tool**: o `Service` sempre usa `EnsureUserID()` do login.
- Os 5 tipos de movimentação: `Diario`, `Saida`, `Cartao`, `Entrada`, `Economia`. Frequências de série: `mensalmente`, `semanalmente`, `diariamente`, `parcelado`.
- `.env`: valores com `$` precisam de aspas simples — o godotenv expande variáveis dentro de aspas duplas (já causou falha de login).
