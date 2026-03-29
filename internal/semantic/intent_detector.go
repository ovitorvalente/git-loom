package semantic

import (
	"path/filepath"
	"strings"
)

func DetectIntent(commitType string, context CommitContext) ChangeIntent {
	scope := NormalizeScopeFromFiles(context.Files)
	target := detectIntentTarget(scope, context)
	intent := detectIntentReason(scope, target, context)

	return ChangeIntent{
		Type:        commitType,
		Scope:       scope,
		Intent:      intent,
		Description: buildIntentDescription(commitType, scope, target, context),
	}
}

func buildIntentDescription(commitType string, scope string, target string, context CommitContext) string {
	if commitType == "test" && strings.TrimSpace(target) != "" {
		return "ajustar " + target
	}

	switch scope {
	case "readme":
		if containsTag(context.Tags, "commit") || containsTag(context.Tags, "cli") {
			return "atualizar instrucoes de uso do cli commit"
		}
		return "atualizar documentacao principal do projeto"
	case "deps":
		return "atualizar dependencias do projeto"
	case "build":
		return "ajustar comandos de desenvolvimento"
	case "config":
		if containsTag(context.Tags, "strict") {
			return "ajustar configuracao do modo estrito"
		}
		return "ajustar configuracao do projeto"
	}

	verb := detectIntentVerb(commitType, context)
	if scope == "cli" && target == "fluxo de commit" {
		return strings.TrimSpace(verb + " fluxo de commit")
	}
	if scope == "ui" && target == "layout do cli" {
		return strings.TrimSpace(verb + " layout do cli")
	}
	if target == "" {
		return strings.TrimSpace(verb + " " + scope)
	}

	return strings.TrimSpace(verb + " " + target)
}

func detectIntentVerb(commitType string, context CommitContext) string {
	switch commitType {
	case "feat":
		if hasOnlyStatus(context.Files, "adicionado") {
			return "adicionar"
		}
		return "evoluir"
	case "fix":
		return "corrigir"
	case "refactor":
		return "refinar"
	case "docs":
		return "atualizar"
	case "test":
		return "ajustar testes de"
	default:
		return "ajustar"
	}
}

func detectIntentTarget(scope string, context CommitContext) string {
	if len(context.Files) == 0 {
		return "repositorio"
	}

	if len(context.Files) == 1 {
		return normalizeIntentTarget(scope, context.Files[0].Path)
	}

	switch scope {
	case "cli":
		return "fluxo de commit"
	case "ui":
		if containsTag(context.Tags, "renderer") || containsTag(context.Tags, "output") || containsTag(context.Tags, "prompt") {
			return "camada de apresentacao do cli"
		}
		return "layout do cli"
	case "git":
		return "repositorio git"
	case "commit":
		return "motor de commits"
	case "app":
		return "planejamento de commits"
	default:
		return scope
	}
}

func normalizeIntentTarget(scope string, path string) string {
	baseName := filepath.Base(path)
	switch {
	case path == "README.md":
		return "documentacao principal do projeto"
	case path == "go.mod" || path == "go.sum":
		return "dependencias do projeto"
	case path == "Makefile":
		return "comandos de desenvolvimento"
	case path == ".gitloom.yaml":
		return "configuracao do projeto"
	case strings.HasSuffix(path, "main.go"):
		return "inicializacao do cli"
	case strings.HasSuffix(path, "root.go"):
		return "comando raiz"
	case strings.HasSuffix(path, "commit.go"):
		return "fluxo de commit"
	case strings.HasSuffix(path, "commit_service.go"):
		return "planejamento de commits"
	case strings.HasSuffix(path, "commit_feedback.go"):
		return "feedback semantico de commit"
	case strings.HasSuffix(path, "prompts.go"):
		return "prompts do cli"
	case strings.HasSuffix(path, "renderer.go"):
		return "renderer do cli"
	case strings.HasSuffix(path, "commit_view.go"):
		return "visao de commit"
	case strings.HasSuffix(path, "summary_view.go"):
		return "resumo do fluxo de commit"
	case strings.HasSuffix(path, "messages.go"):
		return "mensagens do cli"
	case strings.HasSuffix(path, "intent_detector.go"):
		return "heuristicas de intencao semantica"
	case strings.HasSuffix(path, "scope_normalizer.go"):
		return "normalizacao de escopo"
	case strings.HasSuffix(path, "output.go"):
		return "layout do cli"
	case strings.HasSuffix(path, "repository.go"):
		return "repositorio git"
	case strings.HasSuffix(path, "analyzer.go"):
		return "analise semantica de commits"
	case strings.HasSuffix(path, "classifier.go"):
		return "classificacao de commits"
	case strings.HasSuffix(path, "generator.go"):
		return "geracao de mensagens de commit"
	case strings.HasSuffix(path, "_test.go"):
		subject := strings.TrimSuffix(baseName, "_test.go") + ".go"
		if target := normalizeIntentTarget(scope, subject); target != "" {
			return "testes de " + target
		}
		name := strings.TrimSuffix(baseName, "_test.go")
		return "testes de " + strings.ReplaceAll(name, "_", " ")
	case scope != "":
		return scope
	default:
		name := strings.TrimSuffix(baseName, filepath.Ext(baseName))
		return strings.ReplaceAll(name, "_", " ")
	}
}

func detectIntentReason(scope string, target string, context CommitContext) string {
	switch {
	case scope == "readme":
		return "melhorar a clareza de uso para quem adota o projeto"
	case scope == "deps":
		return "manter as dependencias alinhadas com a versao atual do projeto"
	case scope == "build":
		return "simplificar o fluxo de desenvolvimento local"
	case scope == "config":
		return "padronizar o comportamento do cli"
	case containsTag(context.Tags, "preview"):
		return "aumentar a visibilidade do impacto antes do commit"
	case containsTag(context.Tags, "strict"):
		return "reforcar validacoes de qualidade antes de commitar"
	case containsTag(context.Tags, "score"):
		return "medir melhor a qualidade semantica dos commits"
	case containsTag(context.Tags, "suggest"):
		return "orientar agrupamentos e descricoes melhores"
	case containsTag(context.Tags, "prompt") || containsTag(context.Tags, "output"):
		return "melhorar a experiencia visual do cli"
	case target != "":
		return "deixar o fluxo mais claro e previsivel"
	default:
		return "melhorar a organizacao das mudancas"
	}
}

func hasOnlyStatus(files []ChangedFile, status string) bool {
	if len(files) == 0 {
		return false
	}

	for _, file := range files {
		if file.Status != status {
			return false
		}
	}

	return true
}

func containsTag(tags []string, expected string) bool {
	for _, tag := range tags {
		if tag == expected {
			return true
		}
	}

	return false
}
