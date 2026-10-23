# App do Breno

Use um assistente de IA para consultar e organizar suas finanças no App do Breno.

## Pré-requisitos

Para usar o MCP, você precisa de:

- Uma conta ativa no App do Breno e o e-mail e a senha usados para acessá-la.
- Um assistente de IA compatível com MCP, instalado no computador, como Claude
  Code, Claude Desktop, Cursor, Gemini CLI ou Codex.
- Acesso à internet durante o uso.
- Permissão para instalar e configurar ferramentas no computador. Se não tiver,
  peça ajuda à pessoa responsável pelo dispositivo.

## Instalação guiada por IA (para quem não é técnico)

Não sabe usar ferramentas de configuração? Peça ao seu assistente de IA (Claude Code,
Cursor, Gemini CLI ou Codex) para fazer a instalação com você.

**Antes de começar, tenha em mãos:**

1. O **e-mail e a senha** que você usa para entrar no App do Breno.
2. Um assistente de IA que possa ser configurado para usar ferramentas no seu computador.

**Como fazer:**

1. Abra o seu assistente de IA.
2. Copie todo o texto do bloco abaixo.
3. Cole na conversa e envie.
4. Responda às perguntas que ele fizer. Ele vai explicar cada etapa.

```text
Quero instalar e configurar o App do Breno para usar com este assistente de IA.
Sou uma pessoa leiga em tecnologia: explique cada passo em português simples,
um de cada vez. Antes de fazer mudanças no meu computador, explique o que fará
e peça minha confirmação. Não presuma que eu saiba usar ferramentas técnicas.

As instruções e o instalador estão disponíveis no projeto:
https://github.com/vlrosa-dev/app-do-breno-mcp

Descubra qual sistema e assistente estou usando e siga as instruções adequadas.
Use a versão mais recente disponível. Avise-me antes de qualquer etapa que possa
salvar minha senha no computador. Nunca repita minha senha nas mensagens; se eu
preferir, explique como informá-la sem digitá-la na conversa.

Ao terminar, confirme que a conexão funcionou e explique como começar a usar o
App do Breno com perguntas simples.
```

> **Dica:** se algo der errado, copie a mensagem de erro e peça ao assistente para resolver. Para trocar a senha, peça ajuda ao assistente e siga as instruções com atenção.

## Exemplos de perguntas

Depois de conectar o App do Breno ao seu assistente, você pode fazer perguntas
em linguagem natural, por exemplo:

- “Quais são meus saldos deste mês?”
- “Quais foram minhas entradas e saídas em setembro?”
- “Mostre os totais e indicadores deste mês.”
- “Como foi minha vida financeira mês a mês neste ano?”
- “Quanto economizei em cada mês?”
- “Qual foi meu custo de vida mês a mês?”
- “Qual foi meu gasto diário médio em cada mês?”
- “Quais categorias de gastos tenho e quais foram seus totais neste mês?”
- “Quais cartões de crédito estão cadastrados e quando vencem?”
- “Encontre a movimentação de [descrição] neste período.”
- “Registre uma saída de R$ 85,50 em [data], com a descrição [descrição].”
- “Crie uma despesa mensal de R$ 120,00, começando em [data], para [descrição].”

Para alterar ou apagar movimentações, diga ao assistente qual registro deseja
modificar ou remover. Confira os dados apresentados antes de confirmar alterações
ou exclusões.

## Para mantenedores

Esta seção reúne informações técnicas para desenvolvimento e configuração manual.

### Sobre o projeto

O App do Breno MCP é um servidor MCP escrito em Go. Ele se comunica com os
clientes de IA pelo padrão **stdio**. O fluxo principal é:

`cliente de IA → servidor MCP → regras de negócio → API do App do Breno`

As ferramentas e suas regras ficam no código-fonte. O executável é
`app-do-breno-mcp`. O servidor precisa das credenciais da conta:

- `API_USERNAME`: e-mail da conta do App do Breno.
- `API_PASSWORD`: senha da conta.

`API_BASE_URL` e `API_KEY` já têm valores padrão no projeto. As credenciais e o
token de acesso são guardados localmente em `~/.breno/config.json`; proteja esse
conteúdo como uma senha.

### Configuração nos clientes de IA

Use o caminho completo do executável instalado no campo `command`. Nos exemplos,
substitua `/caminho/para/app-do-breno-mcp` pelo caminho correto no computador.

**Claude Code**

```bash
claude mcp add app-do-breno-mcp \
  --env API_USERNAME=seu-email@exemplo.com \
  --env API_PASSWORD=sua-senha \
  -- /caminho/para/app-do-breno-mcp
```

**Claude Desktop, Cursor e Gemini CLI**

Adicione uma entrada `app-do-breno-mcp` à configuração MCP do cliente:

```json
{
  "mcpServers": {
    "app-do-breno-mcp": {
      "command": "/caminho/para/app-do-breno-mcp",
      "env": {
        "API_USERNAME": "seu-email@exemplo.com",
        "API_PASSWORD": "sua-senha"
      }
    }
  }
}
```

No Claude Desktop, essa configuração fica nas preferências do aplicativo.
No Cursor e no Gemini CLI, use a área de configuração MCP do cliente.

**Codex CLI**

Adicione ao `config.toml`:

```toml
[mcp_servers.app-do-breno-mcp]
command = "/caminho/para/app-do-breno-mcp"
args = []

[mcp_servers.app-do-breno-mcp.env]
API_USERNAME = "seu-email@exemplo.com"
API_PASSWORD = "sua-senha"
```

Não substitua configurações existentes do cliente ao adicionar o servidor. Os
exemplos contêm credenciais ilustrativas; não compartilhe sua senha e restrinja
o acesso às configurações que a guardam.

### Desenvolvimento e publicação

Com Go instalado, os principais comandos são:

```bash
make build       # compila em bin/app-do-breno-mcp
make test        # roda os testes unitários
make test-race   # roda os testes com detector de condições de corrida
make cover       # calcula a cobertura dos testes
make lint        # executa go vet e golangci-lint, se disponível
make run-mcp     # inicia o servidor local via stdio
```

Os testes de integração usam `API_USERNAME` e `API_PASSWORD`. Testes de escrita
podem alterar dados reais da conta e só devem ser executados quando isso for
intencional. A publicação das versões é feita pelo GoReleaser ao criar e enviar
uma tag no formato `v*`.


## Aviso

Este é um projeto independente, criado pela comunidade. Não é oficial, não é
mantido e não possui vínculo com o App do Breno ou seus responsáveis.

## Sugestões de melhorias

Tem uma ideia ou encontrou algo que pode melhorar? Compartilhe sua sugestão na
comunidade da Escola do Breno para conversar com outras pessoas e contribuir
com o projeto.

## Segurança

O servidor MCP é um programa de código aberto que roda no seu computador e usa
sua própria conta para acessar os serviços do App do Breno. As credenciais e o
token de acesso ficam guardados localmente, por isso mantenha seu computador e
as configurações do assistente protegidos. Nunca compartilhe sua senha.

O código pode ser consultado por qualquer pessoa, mas isso não elimina todos os
riscos. Use uma versão obtida pelos canais do projeto, confira os dados antes
de confirmar alterações financeiras e avise a comunidade se encontrar algum
problema de segurança.

## Licença

[MIT](LICENSE)
