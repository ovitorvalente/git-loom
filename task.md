Você é um engenheiro de software senior especializado em Go, UX de CLI e arquitetura limpa.

# CONTEXTO

Estou desenvolvendo um CLI chamado "gitloom", responsável por automatizar commits Git com inteligência semântica.

Atualmente o CLI funciona, mas a interface está com problemas:

- Verbosa e redundante
- Baixa hierarquia visual
- Mensagens pouco úteis (ex: "atualizar")
- Feedback pouco acionável
- Experiência mais próxima de logs do que de um assistente

Exemplo atual:

◆ commit 1/1 [docs] (qualidade: 85/100)
• tipo: docs
• escopo: gitignore
• descricao: atualizar
› mensagem: docs(gitignore): atualizar
✧ detalhes: - atualiza

# OBJETIVO

Refatorar COMPLETAMENTE a camada de output/UX do CLI para torná-lo:

- Limpo
- Escaneável
- Profissional
- Orientado a decisão
- Sem redundância
- Com feedback acionável

# REQUISITOS

## 1. NOVO DESIGN DE OUTPUT

Criar um novo formato de exibição com:

- Hierarquia visual clara
- Separação por blocos
- Uso mínimo de texto
- Destaque para informações importantes

Exemplo esperado:

◆ docs(gitignore) [85]

mensagem:
atualizar regras do .gitignore

arquivos:
+4 -1 .gitignore

análise:
⚠ escopo genérico
sugestões:
→ config
→ repo

---

## 2. REMOVER REDUNDÂNCIA

Eliminar:

- descricao duplicada
- detalhes genéricos
- mensagens repetidas

---

## 3. SISTEMA DE RENDERIZAÇÃO

Criar uma camada de UI desacoplada:

internal/ui/
renderer.go
commit_view.go
summary_view.go

Separar:

- domínio (commit)
- apresentação (ui)

---

## 4. MODOS DE EXECUÇÃO

Implementar:

- modo padrão (clean)
- modo verbose (--verbose)

---

## 5. SCORE VISUAL

Transformar:

(qualidade: 85/100)

Em algo visual:

[85] 🟢 bom

---

## 6. FEEDBACK INTELIGENTE

Transformar mensagens genéricas em sugestões acionáveis:

ANTES:
"escopo não identificado"

DEPOIS:
"escopo genérico: gitignore"
"sugestões: config, repo, tooling"

---

## 7. RESUMO FINAL

Adicionar output final mais forte:

✔ 1 commit criado
qualidade média: 85
status: working tree limpa

---

## 8. BOAS PRÁTICAS

- Seguir Clean Architecture
- Código idiomático Go
- Separação clara de responsabilidades
- Sem lógica de negócio na camada de UI
- Código testável

---

# OUTPUT ESPERADO

1. Nova estrutura de arquivos
2. Implementação completa do renderer
3. Refatoração do fluxo atual para usar o novo renderer
4. Sugestões adicionais de melhoria

---

# IMPORTANTE

- NÃO gerar código genérico
- NÃO repetir estrutura atual
- PENSAR como produto real (nível open source premium)
- Código deve ser pronto para produção
