# Security Policy

## Supported Versions

Como o projeto ainda esta em evolucao, a politica atual e simples:

- a branch principal `main` recebe correcoes e melhorias de seguranca
- tags de release mais recentes sao consideradas o canal suportado

## Reporting a Vulnerability

Se voce encontrar uma vulnerabilidade de seguranca:

1. nao abra uma issue publica
2. entre em contato com o maintainer do projeto em canal privado
3. inclua contexto suficiente para reproducao

Ao reportar, inclua:

- versao do `gitloom`
- sistema operacional
- versao do Go, se relevante
- descricao do impacto
- passos para reproducao
- proof of concept minima, se possivel

## Escopo de seguranca

Areas mais sensiveis do projeto:

- execucao de comandos Git
- serializacao JSON
- parser de configuracao
- fluxo de release e distribuicao

## Expectativa de resposta

O projeto e mantido como open source e pode nao ter SLA formal.
Mesmo assim, reports legitimos de seguranca devem receber prioridade acima de melhorias comuns.
