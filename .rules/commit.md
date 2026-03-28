# 📦 commit.md

## 🎯 Objetivo

Este documento define as regras oficiais de commits do projeto **gitloom**.

O objetivo é garantir:

- Histórico limpo e legível
- Commits pequenos e atômicos
- Facilidade de revisão
- Compatibilidade com automação e ferramentas

---

# 🧠 Princípios Fundamentais

## 1. Commits devem ser atômicos

Cada commit deve representar **uma única intenção clara**.

✅ Correto:

- adicionar endpoint de login
- corrigir validação de token

❌ Incorreto:

- adicionar login + corrigir bug + refatorar serviço

---

## 2. Máximo de 4 arquivos por commit

Cada commit deve alterar **no máximo 4 arquivos**.

### Motivo:

- Facilita code review
- Reduz risco de erro
- Melhora rastreabilidade

---

## 3. Commits devem ser pequenos

- Evite grandes blocos de mudança
- Prefira múltiplos commits pequenos

---

## 4. Clareza > Complexidade

Se um commit precisa de explicação longa, ele está errado.

---

# 🧾 Padrão de Commit

Seguimos o padrão Conventional Commits

---

## 📌 Estrutura

```
type(scope): description

- detail 1
- detail 2
```

---

## 📚 Tipos permitidos

| Tipo     | Uso                                       |
| -------- | ----------------------------------------- |
| feat     | Nova funcionalidade                       |
| fix      | Correção de bug                           |
| refactor | Mudança interna sem alterar comportamento |
| chore    | Tarefas técnicas                          |
| docs     | Documentação                              |
| test     | Testes                                    |

---

## ✍️ Regras de escrita

- Usar verbo no imperativo:
  - ✅ add login endpoint
  - ❌ added login endpoint

- Máximo de 72 caracteres na linha principal

- Tudo em minúsculo

- Sem ponto final

---

## 📌 Exemplos

### ✔️ Bom commit

```
feat(auth): add jwt authentication

- implement token validation
- protect private routes
```

---

### ✔️ Commit pequeno e focado

```
fix(auth): handle expired token error
```

---

### ❌ Commit ruim

```
update system
```

---

### ❌ Commit muito grande

```
feat: add auth, fix bug, refactor service, update docs
```

---

# 🧩 Regras de Divisão de Commits

Se houver múltiplas mudanças:

### ❌ Errado:

1 commit com tudo

### ✅ Correto:

```
feat(auth): add login endpoint
fix(auth): fix token validation
refactor(auth): simplify middleware
```

---

# 🔍 Estratégia de Agrupamento

Agrupe commits por:

- Contexto (auth, db, api)
- Tipo de mudança
- Impacto funcional

---

# ⚠️ Regras Estritas

- ❌ Nunca misturar:
  - feature + fix
  - refactor + feature

- ❌ Nunca commitar código quebrado

- ❌ Nunca usar mensagens genéricas:
  - "update"
  - "fix stuff"
  - "changes"

---

# 🔁 Fluxo recomendado

1. Faça mudanças pequenas
2. Revise arquivos modificados
3. Separe commits por intenção
4. Gere mensagem seguindo padrão
5. Valide antes de commitar

---

# 🤖 Integração com gitloom

O gitloom deve:

- Respeitar limite de 4 arquivos
- Dividir commits automaticamente
- Gerar mensagens no padrão correto
- Sugerir melhorias quando necessário

---

# 🧠 Regra Final

> Se o commit não for claro sozinho, ele está errado.

---
