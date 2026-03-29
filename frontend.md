# Frontend Prompt

Use este prompt como base para atualizar o frontend do projeto `gitloom`.

```md
Voce e um designer/frontend engineer senior.

Quero criar ou atualizar o frontend do projeto `gitloom`, um produto focado em:

- revisar commits antes de executar
- agrupar mudancas semanticamente
- gerar mensagens de commit em portugues
- mostrar score, alertas e sugestoes
- expor automacao via JSON
- transmitir confianca, clareza e controle

## Contexto do produto

`gitloom` nao e um app "divertido". Ele e uma ferramenta profissional para developers.
O frontend precisa refletir isso.

Tom da interface:

- claro
- tecnico
- elegante
- direto
- sem hype exagerado
- sem visual genérico de landing page de IA

## Objetivos do frontend

Crie uma experiencia que explique rapidamente:

1. o que o `gitloom` faz
2. por que ele e diferente de um gerador simples de commits
3. como funciona o fluxo de `analyze`, `commit`, `doctor`, `config init` e `version`
4. como instalar, usar e integrar em automacoes
5. como baixar releases multiplataforma

## Direcao visual obrigatoria

- evite roxo como cor dominante
- evite cara de template genérico SaaS
- use uma linguagem visual inspirada em terminal, estrutura, malha, fluxo e composicao
- mantenha leitura excelente em desktop e mobile
- use tipografia forte e intencional
- destaque comandos CLI e outputs com muito cuidado

Paleta base sugerida:

- primario: #0F766E
- primario escuro: #115E59
- destaque: #0EA5E9
- sucesso: #16A34A
- aviso: #D97706
- erro: #DC2626
- fundo: #F8FAFC
- texto forte: #0F172A
- texto secundario: #475569

Fontes sugeridas:

- headlines: IBM Plex Sans, Manrope ou Sora
- body: IBM Plex Sans ou Inter
- monospace: JetBrains Mono ou IBM Plex Mono

## Secoes recomendadas

- Hero
  explicar `gitloom` em uma frase forte
  CTA para instalar e CTA para ver o fluxo

- Product Story
  explicar o problema de commits ruins e historicos confusos

- Workflow
  mostrar visualmente:
  `gitloom analyze`
  `gitloom commit`
  `gitloom doctor`

- Key Features
  agrupar mudancas
  score e feedback
  prompts seguros
  JSON para automacao
  multiplataforma

- Real Output
  mostrar exemplos reais de saida do terminal

- Installation
  Go install
  build local
  releases

- Automation / CI
  destacar `--json`, `doctor --json` e releases por GitHub Actions

- Footer
  GitHub, license, release notes, docs

## Requisitos tecnicos

- componente ou pagina deve ser implementavel de forma real, nao conceitual
- evitar lorem ipsum e textos genéricos
- usar copy coerente com o produto
- incluir estados desktop e mobile
- incluir estrutura responsiva
- se usar React, seguir padroes modernos e evitar overengineering
- se o projeto ja tiver design system, respeitar

## O que eu quero como resposta

- proposta de estrutura da pagina
- direcao visual
- copy principal
- componentes necessarios
- sugestao de layout responsivo
- e, se eu pedir, implementacao pronta em codigo
```
