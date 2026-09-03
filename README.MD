# app-do-breno-agent

Ferramentas para o **App do Breno** compartilhando uma única camada de regras de
negócio (`internal/breno`):

| Binário | O que é |
|---|---|
| `breno-mcp` | Servidor MCP (stdio) — expõe as ferramentas de consulta e movimentação para clientes de LLM |
| `breno` | CLI para operar as mesmas rotas direto do terminal / de scripts / de agentes |

Ambos falam com a API via o mesmo cliente HTTP (`internal/api/client.go`) e o
mesmo fluxo de autenticação com cache de token (`internal/auth`).

## Arquitetura

```
cmd/breno-mcp        cmd/breno-cli
       \                  /
        internal/mcpserver / internal/cli   (só mapeiam entrada -> domínio)
                     |
             internal/breno   ← interface Financas + TODAS as regras de negócio
                     |
             internal/api     ← client.go (REST/RPC, APIError, auto-refresh no 401)
                     |
             internal/auth    ← login por senha + token em ~/.breno/config.json
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
  "base_url": "",   // vazio => default embutido
  "api_key": "",    // vazio => default embutido
  "token": { "access_token": "...", "expires_at": "...", "user_id": "..." }
}
```

Precedência de configuração: **flag > variável de ambiente > `~/.breno/config.json` > default embutido**.

| Variável | Obrigatória | Descrição |
|---|---|---|
| `API_USERNAME` | sim (1º login) | E-mail da conta |
| `API_PASSWORD` | sim (1º login) | Senha da conta |
| `API_KEY` | não | Default embutido (chave anon) |
| `API_BASE_URL` | não | Default embutido |

## CLI (`breno`)

```bash
# autenticação
breno auth login --email voce@exemplo.com   # pede a senha por prompt oculto
breno auth status
breno auth logout

# movimentações (atalhos por tipo)
breno diario   --valor 1.00 --descricao Padaria --recorrencia 0
breno saida    --valor 1500 --descricao Aluguel --recorrencia 12 --frequencia mensalmente
breno cartao   --valor 250  --descricao "Mercado" --recorrencia 0
breno entrada  --valor 8000 --descricao Salário --recorrencia 0
breno economia --valor 500  --descricao "Reserva" --recorrencia 0

# consulta / edição
breno mov listar    --tipo diario --de 2026-09-01 --ate 2026-09-30
breno mov obter     --id <uuid>
breno mov atualizar --id <uuid> --tipo diario --valor 2.00 --descricao Padaria --data 2026-09-03
breno mov deletar   --id <uuid>

# séries recorrentes
breno serie atualizar    --id <recorrencia> --de 2026-09-01 --nova-data 2026-10-01 --tipo saida --valor 1600 --descricao Aluguel
breno serie deletar      --id <recorrencia> --de 2026-10-01
breno serie deletar-tudo --id <recorrencia>

# indicadores
breno saldos  --mes 2026-09-01 --meses 6
breno tags    --mes 2026-09-01 --busca mercado
breno cartoes
breno totais mes          --mes 2026-09-01
breno totais performance  --ano 2026
breno totais custo-vida   --ano 2026
breno totais economia     --ano 2026
breno totais diario-medio --ano 2026
```

- `--recorrencia 0` cria uma movimentação única; `--recorrencia N` (> 0) cria uma
  série com `N` repetições na frequência de `--frequencia`. `--infinita` limita a 24.
- Flag global `--output table|json` (default `table`). Use `json` em scripts/agentes.
- `--base-url` / `--api-key` / `--config` sobrescrevem a config.

## Servidor MCP (`breno-mcp`)

Roda via **stdio**. Configuração igual em todos os clientes: comando (caminho do
binário) + variáveis de ambiente.

### Claude Code

```bash
claude mcp add app-do-breno \
  --env API_USERNAME=seu-email@exemplo.com \
  --env API_PASSWORD=sua-senha \
  -- /caminho/para/breno-mcp
```

Ou no `.mcp.json` / `~/.claude.json`:

```json
{
  "mcpServers": {
    "app-do-breno": {
      "command": "/caminho/para/breno-mcp",
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
    "app-do-breno": {
      "command": "/caminho/para/breno-mcp",
      "env": { "API_USERNAME": "seu-email@exemplo.com", "API_PASSWORD": "sua-senha" }
    }
  }
}
```

### Codex CLI

`~/.codex/config.toml`:

```toml
[mcp_servers.app-do-breno]
command = "/caminho/para/breno-mcp"
args = []

[mcp_servers.app-do-breno.env]
API_USERNAME = "seu-email@exemplo.com"
API_PASSWORD = "sua-senha"
```

## Build

```bash
make build        # bin/breno-mcp + bin/breno
make test         # testes
make test-race    # testes + race detector (requer cgo)
make cover        # cobertura
make lint         # go vet (+ golangci-lint se instalado)
```

Ou direto (o pacote `cmd/breno-cli` gera um executável chamado `breno-cli` — renomeie
para `breno` se quiser):

```bash
go install github.com/vlrosa-dev/app-do-breno-agent/cmd/breno-cli@latest
go install github.com/vlrosa-dev/app-do-breno-agent/cmd/breno-mcp@latest
```

## Release (mantenedores)

`.github/workflows/release.yml` roda o GoReleaser ao empurrar uma tag `v*` e
publica os dois binários (`breno-mcp_*`, `breno_*`) para linux/windows/darwin
em amd64/arm64:

```bash
git tag v0.1.0
git push origin v0.1.0
```
