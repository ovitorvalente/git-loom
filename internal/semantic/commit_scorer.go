package semantic

import "strings"

func ScoreCommit(intent ChangeIntent, context CommitContext) CommitQuality {
	score := 100
	reasons := []string{}
	criteria := []QualityCriteria{}

	fileCount := len(context.Files)
	if fileCount > 4 {
		score -= 40
		reasons = append(reasons, "bloco excede o limite recomendado de 4 arquivos")
		criteria = append(criteria, QualityCriteria{
			Name: "tamanho", Passed: false, Warning: true,
			Message: "excede limite de 4 arquivos",
		})
	} else {
		criteria = append(criteria, QualityCriteria{
			Name: "tamanho", Passed: true,
		})
	}

	if intent.Scope == "" {
		score -= 15
		reasons = append(reasons, "escopo nao foi identificado com clareza")
		criteria = append(criteria, QualityCriteria{
			Name: "escopo", Passed: false, Warning: true,
			Message: "escopo nao identificado",
		})
	} else if IsGenericScope(intent.Scope) {
		criteria = append(criteria, QualityCriteria{
			Name: "escopo", Passed: false, Warning: true,
			Message: "escopo generico: " + intent.Scope,
		})
	} else {
		criteria = append(criteria, QualityCriteria{
			Name: "escopo", Passed: true,
		})
	}

	if isGenericDescription(intent.Description) {
		score -= 25
		reasons = append(reasons, "descricao ainda esta generica")
		criteria = append(criteria, QualityCriteria{
			Name: "descricao", Passed: false, Warning: true,
			Message: "descricao generica",
		})
	} else {
		criteria = append(criteria, QualityCriteria{
			Name: "descricao", Passed: true,
		})
	}

	if hasMixedScopes(context.Files) {
		score -= 15
		reasons = append(reasons, "arquivos misturam contextos diferentes")
		criteria = append(criteria, QualityCriteria{
			Name: "coerencia", Passed: false, Warning: true,
			Message: "contextos mistos",
		})
	} else {
		criteria = append(criteria, QualityCriteria{
			Name: "coerencia", Passed: true,
		})
	}

	if expectedType := inferExpectedType(context); expectedType != "" && expectedType != intent.Type {
		score -= 20
		reasons = append(reasons, "tipo pode nao refletir bem a mudanca")
		criteria = append(criteria, QualityCriteria{
			Name: "tipo", Passed: false, Warning: true,
			Message: "tipo pode nao refletir bem a mudanca",
		})
	} else {
		criteria = append(criteria, QualityCriteria{
			Name: "tipo", Passed: true,
		})
	}

	if score < 0 {
		score = 0
	}

	return CommitQuality{
		Score:    score,
		Reasons:  reasons,
		Criteria: criteria,
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
