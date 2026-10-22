# App do Breno MCP

Servidor MCP (stdio) do **App do Breno** — expõe as ferramentas de consulta e
movimentação financeira para clientes de LLM (binário `app-do-breno-mcp`).

As regras de negócio ficam numa camada própria (`internal/breno`); a API é
acessada por um único cliente HTTP (`internal/api/client.go`) com autenticação
e cache de token (`internal/auth`).

## Arquitetura

```
cmd/app-do-breno-mcp
       |
internal/mcpserver   ← só mapeia entrada das tools -> domínio
       |
internal/breno       ← interface Financas + TODAS as regras de negócio
       |
internal/api         ← client.go (REST/RPC, APIError, relogin no 401)
       |
internal/auth        ← login por senha + token em ~/.breno/config.json
```

## Autenticação

O login usa **e-mail + senha** (`grant_type=password`). O token é salvo em
`~/.breno/config.json` (texto plano, permissão `0600`) junto das credenciais, e
reaproveitado até expirar; ao expirar, um novo login é feito automaticamente pela
mesma rota — sem prompt.

```jsonc
// ~/.breno/config.json
{
  "email": "voce@exemplo.com",
  "senha": "...",
  "token": { "access_token": "...", "expires_at": "...", "user_id": "..." }
}
```

Precedência de configuração: **variável de ambiente > `~/.breno/config.json` (credenciais) > default embutido**.

| Variável | Obrigatória | Descrição |
|---|---|---|
| `API_USERNAME` | sim (1º login) | E-mail da conta |
| `API_PASSWORD` | sim (1º login) | Senha da conta |
| `API_KEY` | não | Default embutido (chave anon) |
| `API_BASE_URL` | não | Default embutido |

## Servidor MCP (`app-do-breno-mcp`)

Roda via **stdio**. Configuração igual em todos os clientes: comando (caminho do
binário) + variáveis de ambiente.

### Claude Code

```bash
claude mcp add app-do-breno-mcp \
  --env API_USERNAME=seu-email@exemplo.com \
  --env API_PASSWORD=sua-senha \
  -- /caminho/para/app-do-breno-mcp
```

Ou no `.mcp.json` / `~/.claude.json`:

```json
{
  "mcpServers": {
    "app-do-breno-mcp": {
      "command": "/caminho/para/app-do-breno-mcp",
      "env": { "API_USERNAME": "seu-email@exemplo.com", "API_PASSWORD": "sua-senha" }
    }
  }
}
```

### Cursor / Gemini CLI

`~/.cursor/mcp.json` ou `~/.gemini/settings.json`:

```json
{
  "mcpServers": {
    "app-do-breno-mcp": {
      "command": "/caminho/para/app-do-breno-mcp",
      "env": { "API_USERNAME": "seu-email@exemplo.com", "API_PASSWORD": "sua-senha" }
    }
  }
}
```

### Codex CLI

`~/.codex/config.toml`:

```toml
[mcp_servers.app-do-breno-mcp]
command = "/caminho/para/app-do-breno-mcp"
args = []

[mcp_servers.app-do-breno-mcp.env]
API_USERNAME = "seu-email@exemplo.com"
API_PASSWORD = "sua-senha"
```

## Build

```bash
make build        # bin/app-do-breno-mcp
make test         # testes
make test-race    # testes + race detector (requer cgo)
make cover        # cobertura
make lint         # go vet (+ golangci-lint se instalado)
```

Ou direto:

```bash
go install github.com/vlrosa-dev/app-do-breno-mcp/cmd/app-do-breno-mcp@latest
```

## Release (mantenedores)

`.github/workflows/release.yml` roda o GoReleaser ao empurrar uma tag `v*` e
publica o binário `app-do-breno-mcp_*` para linux/windows/darwin
em amd64/arm64:

```bash
git tag v0.1.0
git push origin v0.1.0
```
