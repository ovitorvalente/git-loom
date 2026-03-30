package semantic

import (
	"path/filepath"
	"sort"
	"strings"
)

type ScopeSuggestion struct {
	Current      string
	Alternatives []string
	Generic      bool
}

type DescriptionSuggestion struct {
	Reason  string
	Current string
	Generic bool
}

func SuggestScope(scope string, files []ChangedFile) ScopeSuggestion {
	normalized := strings.ToLower(strings.TrimSpace(scope))
	if !IsGenericScope(normalized) {
		return ScopeSuggestion{Current: scope, Generic: false}
	}

	candidates := map[string]bool{}
	for _, file := range files {
		for _, suggestion := range scopeCandidatesFromFilePath(file.Path) {
			if suggestion != normalized {
				candidates[suggestion] = true
			}
		}
	}

	alternatives := make([]string, 0, len(candidates))
	for candidate := range candidates {
		alternatives = append(alternatives, candidate)
	}
	sort.Strings(alternatives)

	if len(alternatives) > 3 {
		alternatives = alternatives[:3]
	}

	return ScopeSuggestion{
		Current:      scope,
		Generic:      true,
		Alternatives: alternatives,
	}
}

func SuggestDescription(desc string) DescriptionSuggestion {
	normalized := strings.ToLower(strings.TrimSpace(desc))
	if !IsGenericDescription(normalized) {
		return DescriptionSuggestion{Current: desc, Generic: false}
	}

	return DescriptionSuggestion{
		Current: desc,
		Generic: true,
		Reason:  "descricao generica detectada",
	}
}

func IsGenericScope(scope string) bool {
	switch strings.ToLower(strings.TrimSpace(scope)) {
	case "", "core", "app", "cli", "ui", "repo", "project", "misc", "general", "gitignore":
		return true
	default:
		return false
	}
}

func IsGenericDescription(desc string) bool {
	normalized := strings.ToLower(strings.TrimSpace(desc))
	genericDescriptions := []string{
		"atualizar projeto", "ajustar projeto",
		"atualizar repositorio", "ajustar repositorio",
		"atualizar cli", "ajustar cli",
		"atualizar core", "ajustar core",
		"refinar core", "refinar app",
		"refinar cli", "refinar ui",
		"refinar config", "refinar commit",
		"atualizar", "ajustar", "modificar", "alterar",
	}

	for _, generic := range genericDescriptions {
		if normalized == generic {
			return true
		}
	}

	return false
}

func scopeCandidatesFromFilePath(path string) []string {
	base := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	base = strings.TrimPrefix(strings.ToLower(base), ".")
	directory := strings.ToLower(filepath.Dir(path))

	candidates := []string{}
	for _, part := range strings.Split(directory, "/") {
		part = strings.TrimSpace(part)
		if part == "" || part == "." || part == "internal" || part == "cmd" {
			continue
		}
		if part == "ui" || part == "cli" || part == "domain" || part == "semantic" || part == "shared" || part == "app" {
			continue
		}
		candidates = append(candidates, part)
	}

	switch {
	case strings.Contains(path, "gitignore"):
		candidates = append(candidates, "config", "repo", "tooling")
	case strings.Contains(path, "go.mod"), strings.Contains(path, "go.sum"), strings.Contains(path, "makefile"):
		candidates = append(candidates, "deps", "build", "tooling")
	case strings.HasSuffix(path, "_test.go"):
		candidates = append(candidates, "test")
	}

	if base != "" && base != "readme" && base != "main" {
		candidates = append(candidates, strings.ReplaceAll(base, "_", "-"))
	}

	return uniqueSuggestions(candidates)
}

func uniqueSuggestions(values []string) []string {
	seen := map[string]bool{}
	result := make([]string, 0, len(values))

	for _, value := range values {
		normalized := strings.TrimSpace(strings.ToLower(value))
		if normalized == "" || seen[normalized] {
			continue
		}
		seen[normalized] = true
		result = append(result, normalized)
	}

	return result
}
