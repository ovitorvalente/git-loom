package semantic

import "strings"

func ScoreCommit(intent ChangeIntent, context CommitContext) CommitQuality {
	score := 100
	reasons := []string{}

	if len(context.Files) > 4 {
		score -= 40
		reasons = append(reasons, "bloco excede o limite recomendado de 4 arquivos")
	}
	if intent.Scope == "" {
		score -= 15
		reasons = append(reasons, "escopo nao foi identificado com clareza")
	}
	if isGenericDescription(intent.Description) {
		score -= 25
		reasons = append(reasons, "descricao ainda esta generica")
	}
	if hasMixedScopes(context.Files) {
		score -= 15
		reasons = append(reasons, "arquivos misturam contextos diferentes")
	}
	if expectedType := inferExpectedType(context); expectedType != "" && expectedType != intent.Type {
		score -= 20
		reasons = append(reasons, "tipo pode nao refletir bem a mudanca")
	}

	if score < 0 {
		score = 0
	}

	return CommitQuality{
		Score:   score,
		Reasons: reasons,
	}
}

func isGenericDescription(description string) bool {
	normalizedDescription := strings.ToLower(strings.TrimSpace(description))
	genericDescriptions := []string{
		"atualizar projeto",
		"ajustar projeto",
		"atualizar repositorio",
		"ajustar repositorio",
		"atualizar cli",
		"ajustar cli",
		"atualizar core",
		"ajustar core",
		"refinar core",
		"refinar app",
		"refinar cli",
		"refinar ui",
		"refinar config",
		"refinar commit",
	}

	for _, genericDescription := range genericDescriptions {
		if normalizedDescription == genericDescription {
			return true
		}
	}

	return false
}

func hasMixedScopes(files []ChangedFile) bool {
	scopeVotes := map[string]bool{}
	for _, file := range files {
		scopeVotes[NormalizeScope(file.Path)] = true
	}

	return len(scopeVotes) > 2
}

func inferExpectedType(context CommitContext) string {
	if len(context.Files) == 0 {
		return ""
	}

	allTests := true
	allDocs := true
	allDeps := true

	for _, file := range context.Files {
		path := strings.ToLower(file.Path)
		if !strings.HasSuffix(path, "_test.go") {
			allTests = false
		}
		if path != "readme.md" && !strings.HasSuffix(path, ".md") {
			allDocs = false
		}
		if path != "go.mod" && path != "go.sum" && path != "makefile" {
			allDeps = false
		}
	}

	switch {
	case allTests:
		return "test"
	case allDocs:
		return "docs"
	case allDeps:
		return "chore"
	default:
		return ""
	}
}
