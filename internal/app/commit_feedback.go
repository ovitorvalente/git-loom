package app

import (
	"path/filepath"
	"sort"
	"strings"
)

type CommitFeedback struct {
	Highlights  []string
	Suggestions []string
}

func BuildCommitFeedback(plan CommitPlan) CommitFeedback {
	highlights := buildFeedbackHighlights(plan)
	suggestions := buildFeedbackSuggestions(plan)

	return CommitFeedback{
		Highlights:  highlights,
		Suggestions: suggestions,
	}
}

func buildFeedbackHighlights(plan CommitPlan) []string {
	highlights := []string{}
	scope := strings.TrimSpace(plan.Result.Commit.Scope)

	switch {
	case scope == "":
		highlights = append(highlights, "escopo ausente")
	case isGenericScope(scope):
		highlights = append(highlights, "escopo generico: "+scope)
	}

	description := strings.TrimSpace(strings.ToLower(plan.Result.Commit.Description))
	if description == "" {
		highlights = append(highlights, "descricao ausente")
	} else if isGenericFeedbackDescription(description) {
		highlights = append(highlights, "descricao generica")
	}

	for _, reason := range plan.Quality.Reasons {
		if !containsFeedback(highlights, reason) {
			highlights = append(highlights, reason)
		}
	}

	return highlights
}

func isGenericFeedbackDescription(description string) bool {
	switch strings.TrimSpace(strings.ToLower(description)) {
	case "atualizar projeto",
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
		"refinar commit":
		return true
	default:
		return false
	}
}

func buildFeedbackSuggestions(plan CommitPlan) []string {
	candidates := map[string]bool{}
	currentScope := strings.TrimSpace(strings.ToLower(plan.Result.Commit.Scope))
	pathCandidates := []string{}

	for _, file := range plan.Context.Files {
		for _, suggestion := range scopeCandidatesFromPath(file.Path) {
			if suggestion == "" || suggestion == currentScope {
				continue
			}
			candidates[suggestion] = true
			pathCandidates = append(pathCandidates, suggestion)
		}
	}

	if len(pathCandidates) == 0 {
		for _, tag := range plan.Context.Tags {
			for _, suggestion := range scopeCandidatesFromTag(tag) {
				if suggestion == "" || suggestion == currentScope {
					continue
				}
				candidates[suggestion] = true
			}
		}
	}

	suggestions := make([]string, 0, len(candidates))
	for suggestion := range candidates {
		suggestions = append(suggestions, suggestion)
	}
	sort.Strings(suggestions)

	if len(suggestions) > 3 {
		suggestions = suggestions[:3]
	}

	return suggestions
}

func isGenericScope(scope string) bool {
	switch strings.ToLower(strings.TrimSpace(scope)) {
	case "", "core", "app", "cli", "ui", "repo", "project", "misc", "general", "gitignore":
		return true
	default:
		return false
	}
}

func scopeCandidatesFromPath(path string) []string {
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

	return uniqueFeedback(candidates)
}

func scopeCandidatesFromTag(tag string) []string {
	switch strings.TrimSpace(strings.ToLower(tag)) {
	case "config":
		return []string{"config"}
	case "prompt":
		return []string{"ux"}
	case "output", "ui":
		return []string{"ux", "output"}
	case "suggest":
		return []string{"semantic"}
	case "strict", "score":
		return []string{"quality"}
	case "build":
		return []string{"build"}
	case "test":
		return []string{"test"}
	default:
		return nil
	}
}

func uniqueFeedback(values []string) []string {
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

func containsFeedback(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}

	return false
}
