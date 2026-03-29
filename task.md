Você é um engenheiro de software senior especializado em Go, design de CLI (DX/UX) e arquitetura limpa.

# CONTEXTO

Estou desenvolvendo um CLI chamado "gitloom", focado em automação inteligente de commits Git.

O CLI já funciona, mas o output atual é técnico demais e pouco eficiente em UX:

- Mensagens redundantes
- Baixa clareza visual
- Feedback pouco acionável
- Falta de assistividade (parece um logger, não um assistente)

Exemplo atual:

◆ commit 1/1 [docs] (qualidade: 85/100)
• tipo: docs
• escopo: gitignore
• descricao: atualizar
› mensagem: docs(gitignore): atualizar
✧ detalhes: - atualiza

# OBJETIVO

Transformar o CLI em uma experiência PREMIUM, com:

- Output limpo, escaneável e profissional
- Feedback inteligente e acionável
- Redução de ruído
- Sensação de “assistente inteligente”

---

# TAREFAS

## 1. REFATORAR OUTPUT (PRINCIPAL)

Criar um novo formato padrão:

◆ docs(config) [82] 🟢

atualizar regras do .gitignore

+4 -1 .gitignore

análise:
⚠ escopo genérico: "gitignore"
sugestões:
→ config
→ repo

---

## 2. REMOVER REDUNDÂNCIA

Eliminar completamente:

- "descricao"
- "detalhes" genéricos
- duplicação entre mensagem e descrição

Manter apenas informação útil.

---

## 3. IMPLEMENTAR SISTEMA DE SUGESTÕES

Criar um módulo:

internal/semantic/suggestions.go

Responsável por:

- detectar mensagens genéricas ("atualizar", "ajustes")
- detectar escopos fracos ("file", "gitignore")
- sugerir melhorias

Exemplo:

ANTES:
escopo: gitignore

DEPOIS:
⚠ escopo genérico
→ sugestões: config, repo, tooling

---

## 4. MELHORAR SCORE DE QUALIDADE

Transformar score em algo mais rico:

qualidade: 82 🟢 bom

critérios:
✔ tamanho adequado
✔ tipo correto
⚠ descrição genérica

---

## 5. IMPLEMENTAR PREVIEW FINAL

Antes de criar commits:

Resumo:

1 commit será criado

docs(config)
→ atualizar regras do .gitignore

---

## 6. ADICIONAR MODO CLEAN E VERBOSE

- default: clean (mínimo necessário)
- flag: --verbose (mostra análise completa)

---

## 7. ADICIONAR MICROINTERAÇÕES

Durante execução:

- "analisando mudanças..."
- "gerando commits..."
- "validando qualidade..."

---

## 8. IMPLEMENTAR RENDERER DESACOPLADO

Criar:

internal/ui/renderer/
renderer.go
commit_view.go
summary_view.go

Separar:

- domínio (commit)
- apresentação (ui)

---

## 9. IMPLEMENTAR COMMIT PREVIEW INTELIGENTE

Mostrar resumo do diff:

preview:

- adiciona regra para ignorar arquivos temporários

* remove regra duplicada

---

## 10. CRIAR ALIAS CURTO PARA O CLI

Objetivo: tornar o comando mais rápido de usar.

Sugestões de alias:

- gl
- gm (git message)
- gcx (git commit extended)
- loom

Escolher o melhor com base em:

- memorabilidade
- ausência de conflito com git

---

## IMPLEMENTAÇÃO DO ALIAS

### 1. Nome alternativo do binário

Permitir execução:

gitloom
gl

### 2. Suporte no Cobra

Ajustar comando root para aceitar alias:

Use: "gitloom"
Aliases: ["gl", "loom"]

---

## 11. MELHORAR OUTPUT FINAL

ANTES:
✔ fluxo finalizado

DEPOIS:

✔ 1 commit criado com sucesso

qualidade média: 82
status: working tree limpa

---

# REQUISITOS TÉCNICOS

- Go idiomático
- Clean Architecture
- Separação de responsabilidades
- Código testável
- Sem lógica de negócio na UI

---

# OUTPUT ESPERADO

1. Refatoração completa da camada de UI
2. Implementação do sistema de sugestões
3. Novo renderer desacoplado
4. Implementação de alias no CLI
5. Código pronto para produção

---

# IMPORTANTE

- Não gerar código genérico
- Pensar como produto real open source
- Priorizar experiência do desenvolvedor (DX)
- Criar algo que pareça uma ferramenta premium
