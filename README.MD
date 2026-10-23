# App do Breno MCP

Servidor MCP (stdio) do **App do Breno** — expõe as ferramentas de consulta e
movimentação financeira para clientes de LLM (binário `app-do-breno-mcp`).

As regras de negócio ficam numa camada própria (`internal/breno`); a API é
acessada por um único cliente HTTP (`internal/api/client.go`) com autenticação
e cache de token (`internal/auth`).

## Instalação guiada por IA (para quem não é técnico)

Não sabe usar terminal? Sem problema: deixe o seu assistente de IA (Claude Code,
Cursor, Gemini CLI ou Codex) fazer a instalação por você.

**Antes de começar, tenha em mãos:**

1. O **e-mail e a senha** que você usa para entrar no App do Breno.
2. Um assistente de IA que rode no seu computador (ex.: [Claude Code](https://claude.com/claude-code)).

**Como fazer:**

1. Abra o seu assistente de IA.
2. Copie **todo** o texto do bloco abaixo (no GitHub, use o botão de copiar no
   canto do bloco).
3. Cole na conversa e envie.
4. Responda às perguntas que ele fizer. Ele vai te avisar a cada passo.

```text
Quero instalar e configurar o servidor MCP "app-do-breno-mcp" neste computador.
Sou uma pessoa leiga em tecnologia: explique cada passo em português simples,
um de cada vez, e me peça confirmação antes de qualquer ação que altere o meu
computador. Não presuma que eu sei o que é terminal, variável de ambiente ou JSON.

Informações sobre o projeto:
- Repositório (público): https://github.com/vlrosa-dev/app-do-breno-mcp
- Os binários prontos ficam na página "Releases" (use a versão mais recente).
  Não é preciso conta no GitHub para baixar.
  Nome dos arquivos: app-do-breno-mcp_<sistema>_<arquitetura>.tar.gz
  (no Windows é .zip), onde <sistema> é linux, darwin (macOS) ou windows e
  <arquitetura> é amd64 (Intel/AMD) ou arm64 (Apple Silicon M1/M2/M3/M4,
  Raspberry Pi, etc.).
- O servidor funciona via stdio e só precisa de duas variáveis de ambiente:
  API_USERNAME (e-mail da conta do App do Breno) e API_PASSWORD (senha).
  API_BASE_URL e API_KEY NÃO devem ser configuradas: já têm valor embutido.

Siga estes passos:

1. Descubra sozinho qual é o meu sistema operacional e a arquitetura do
   processador e me diga em uma frase (ex.: "Você usa um Mac com chip Apple").

2. Descubra qual assistente/cliente de IA eu estou usando (Claude Code, Claude
   Desktop, Cursor, Gemini CLI ou Codex CLI). Se não conseguir identificar,
   me pergunte.

3. Baixe o arquivo da release mais recente correspondente ao meu sistema:
   - O link direto da versão mais recente é
     https://github.com/vlrosa-dev/app-do-breno-mcp/releases/latest/download/<nome-do-arquivo>
     — baixe com curl (Linux/macOS) ou Invoke-WebRequest (Windows).
   - Se o download pelo terminal não funcionar, me explique como baixar pelo
     navegador na página https://github.com/vlrosa-dev/app-do-breno-mcp/releases
     e me diga em qual pasta salvar.

4. Extraia o arquivo e coloque o executável em um local fixo e fácil:
   - Linux/macOS: ~/.local/bin/app-do-breno-mcp (crie a pasta se não existir)
     e dê permissão de execução (chmod +x).
   - macOS: remova o bloqueio de segurança do Gatekeeper com
     "xattr -d com.apple.quarantine <caminho>" (explique que é normal, pois o
     programa não é assinado pela Apple).
   - Windows: %LOCALAPPDATA%\app-do-breno-mcp\app-do-breno-mcp.exe
   Guarde o CAMINHO ABSOLUTO completo do executável: ele será usado na
   configuração (não use "~" nem caminhos relativos no arquivo de configuração).

5. Peça meu e-mail e minha senha do App do Breno. Antes de pedir, me avise que:
   - a senha ficará salva em texto no arquivo de configuração do meu cliente de
     IA e em ~/.breno/config.json, apenas no meu computador;
   - se eu preferir não digitar a senha no chat, você pode me passar o comando
     pronto com um espaço reservado para eu mesmo colar a senha.
   Nunca repita a senha de volta nas suas mensagens.

6. Registre o servidor MCP com o nome "app-do-breno-mcp" no meu cliente:
   - Claude Code: rode
     claude mcp add app-do-breno-mcp --scope user \
       --env API_USERNAME=<email> --env API_PASSWORD=<senha> \
       -- <caminho-absoluto-do-executavel>
   - Claude Desktop: edite o claude_desktop_config.json (macOS:
     ~/Library/Application Support/Claude/, Windows: %APPDATA%\Claude\).
   - Cursor: ~/.cursor/mcp.json
   - Gemini CLI: ~/.gemini/settings.json
   - Codex CLI: ~/.codex/config.toml (seção [mcp_servers.app-do-breno-mcp] e
     [mcp_servers.app-do-breno-mcp.env]).
   Nos arquivos JSON use o formato:
   {"mcpServers": {"app-do-breno-mcp": {"command": "<caminho>",
     "env": {"API_USERNAME": "<email>", "API_PASSWORD": "<senha>"}}}}
   Se o arquivo já existir, faça um backup antes e PRESERVE tudo o que já
   estiver nele — apenas acrescente a entrada "app-do-breno-mcp". Garanta que
   o JSON final continue válido (atenção a aspas e caracteres especiais da
   senha, como $, \ e ").

7. Teste se o login funciona antes de terminar: rode o executável com
   API_USERNAME e API_PASSWORD definidos e confirme que ele inicia sem erro de
   autenticação (ele fica aguardando entrada; encerre depois de alguns
   segundos). Se aparecer erro de login, me ajude a conferir e-mail e senha.

8. Me diga para fechar e abrir de novo o meu cliente de IA (ou iniciar uma nova
   sessão) e, depois disso, para pedir algo como:
   "Quais são meus saldos deste mês no App do Breno?"
   Explique como eu vejo se o servidor está conectado (no Claude Code, pelo
   comando /mcp).

Ao final, faça um resumo curto do que foi instalado, onde ficou o executável e
qual arquivo de configuração foi alterado.
```

> **Dica:** se algo der errado, cole a mensagem de erro de volta na conversa e
> peça ao assistente para resolver. Para trocar a senha no futuro, peça:
> *"Atualize a senha do servidor MCP app-do-breno-mcp"* e apague o arquivo
> `~/.breno/config.json` para forçar um novo login.

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

### Exemplos de perguntas

Depois de conectar o servidor ao seu assistente de IA, você pode fazer perguntas
em linguagem natural, por exemplo:

- “Quais são meus saldos deste mês?”
- “Quais foram minhas entradas e saídas em setembro de 2026?”
- “Mostre os totais e indicadores de agosto de 2026.”
- “Como foi minha performance financeira mês a mês em 2026?”
- “Quanto economizei em cada mês de 2026?”
- “Qual foi meu custo de vida mês a mês em 2026?”
- “Qual foi meu gasto diário médio em cada mês de 2026?”
- “Quais tags tenho e quais foram seus totais neste mês?”
- “Quais cartões de crédito estão cadastrados e quando fecham e vencem?”
- “Encontre a movimentação de [descrição] no período de [data inicial] a [data final].”
- “Registre uma saída de R$ 85,50 em [data], com a descrição [descrição].”
- “Crie uma despesa mensal recorrente de R$ 120,00, começando em [data], para [descrição].”

Para alterar ou apagar movimentações, informe qual registro deseja modificar ou
remover. Para atualizar uma movimentação recorrente, o assistente deve primeiro
localizá-la pelo ID. Confira os dados antes de pedir alterações ou exclusões.

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

## Licença

[MIT](LICENSE)
