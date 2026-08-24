# app-do-breno-agent

Servidor MCP (Model Context Protocol) para o App do Breno — expõe ferramentas de consulta e
movimentação financeira (diário, saída fixa, cartão, entrada, economia) para uso em clientes de
LLM compatíveis com MCP.

## Variáveis de ambiente

| Variável | Obrigatória | Descrição |
|---|---|---|
| `API_USERNAME` | sim | E-mail da sua conta no App do Breno |
| `API_PASSWORD` | sim | Senha da sua conta no App do Breno |
| `API_KEY` | não | Chave de API. Já vem com um valor padrão (chave pública/anon) embutido no binário — só defina se precisar usar uma chave própria |
| `API_BASE_URL` | não | URL base da API. Já vem com um valor padrão embutido no binário — só defina se precisar apontar para outro ambiente |

Só `API_USERNAME` e `API_PASSWORD` são realmente pessoais e obrigatórias; o binário é o mesmo para todo mundo.

## Instalação

1. Baixe o binário para o seu sistema operacional na página de [Releases](https://github.com/vlrosa-dev/app-do-breno-agent/releases) (`app-do-breno-mcp_<os>_<arch>.tar.gz` ou `.zip` no Windows).

   > Releases automatizados cobrem **Linux**, **Windows** e **macOS (darwin)**, todos em amd64/arm64.

2. Extraia o arquivo e opcionalmente mova o binário para um diretório no seu `PATH`:

   ```bash
   tar -xzf app-do-breno-mcp_linux_amd64.tar.gz
   chmod +x app-do-breno-mcp
   sudo mv app-do-breno-mcp /usr/local/bin/
   ```

3. Anote o caminho completo do binário (ex.: `/usr/local/bin/app-do-breno-mcp`) — ele será usado na configuração do cliente MCP.

### Alternativa: compilar do código-fonte

Requer Go 1.26+:

```bash
go install github.com/vlrosa-dev/app-do-breno-agent/cmd/mcp@latest
```

O binário fica em `$(go env GOPATH)/bin/mcp` (renomeie se quiser).

## Configurando nos clientes de LLM

O servidor roda via **stdio**, então a configuração em todos os clientes segue o mesmo padrão:
comando (caminho do binário) + variáveis de ambiente. Substitua `/caminho/para/app-do-breno-mcp`
e as credenciais pelos seus valores reais. Só `API_USERNAME` e `API_PASSWORD` são necessárias;
`API_KEY` e `API_BASE_URL` são opcionais e usam os valores padrão embutidos no binário caso não
sejam definidas.

### Claude Code

Via CLI:

```bash
claude mcp add app-do-breno \
  --env API_USERNAME=seu-email@exemplo.com \
  --env API_PASSWORD=sua-senha \
  -- /caminho/para/app-do-breno-mcp
```

Ou manualmente no `.mcp.json` (na raiz do projeto) ou `~/.claude.json` (escopo global):

```json
{
  "mcpServers": {
    "app-do-breno": {
      "command": "/caminho/para/app-do-breno-mcp",
      "env": {
        "API_USERNAME": "seu-email@exemplo.com",
        "API_PASSWORD": "sua-senha"
      }
    }
  }
}
```

### Cursor

Edite (ou crie) `~/.cursor/mcp.json` (global) ou `.cursor/mcp.json` (por projeto):

```json
{
  "mcpServers": {
    "app-do-breno": {
      "command": "/caminho/para/app-do-breno-mcp",
      "env": {
        "API_USERNAME": "seu-email@exemplo.com",
        "API_PASSWORD": "sua-senha"
      }
    }
  }
}
```

Depois: `Settings → MCP` e confirme que o servidor `app-do-breno` aparece como conectado.

### Codex CLI

Edite `~/.codex/config.toml` e adicione:

```toml
[mcp_servers.app-do-breno]
command = "/caminho/para/app-do-breno-mcp"
args = []

[mcp_servers.app-do-breno.env]
API_USERNAME = "seu-email@exemplo.com"
API_PASSWORD = "sua-senha"
```

### Gemini CLI

Edite `~/.gemini/settings.json`:

```json
{
  "mcpServers": {
    "app-do-breno": {
      "command": "/caminho/para/app-do-breno-mcp",
      "env": {
        "API_USERNAME": "seu-email@exemplo.com",
        "API_PASSWORD": "sua-senha"
      }
    }
  }
}
```

## Publicando uma nova versão (mantenedores)

O pipeline em `.github/workflows/release.yml` usa [GoReleaser](https://goreleaser.com) e roda
automaticamente ao empurrar uma tag `v*`:

```bash
git tag v0.1.0
git push origin v0.1.0
```

Isso builda os binários, empacota e publica uma Release no GitHub com os artefatos e checksums.
