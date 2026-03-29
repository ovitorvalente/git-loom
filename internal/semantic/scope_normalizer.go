package semantic

import (
	"path/filepath"
	"strings"
)

func NormalizeScopeFromFiles(files []ChangedFile) string {
	scopeVotes := map[string]int{}
	for _, file := range files {
		scope := NormalizeScope(file.Path)
		if scope == "" {
			continue
		}
		scopeVotes[scope]++
	}

	return mostCommon(scopeVotes)
}

func NormalizeScope(path string) string {
	normalizedPath := strings.TrimSpace(path)
	switch {
	case normalizedPath == "":
		return ""
	case normalizedPath == "README.md":
		return "readme"
	case normalizedPath == "go.mod", normalizedPath == "go.sum":
		return "deps"
	case normalizedPath == "Makefile":
		return "build"
	case normalizedPath == ".gitloom.yaml":
		return "config"
	case strings.HasPrefix(normalizedPath, "cmd/"), strings.HasPrefix(normalizedPath, "internal/cli/"):
		return "cli"
	case strings.HasPrefix(normalizedPath, "internal/ui/"):
		return "ui"
	case strings.HasPrefix(normalizedPath, "internal/infra/git/"):
		return "git"
	case strings.HasPrefix(normalizedPath, "internal/infra/config/"):
		return "config"
	case strings.HasPrefix(normalizedPath, "internal/domain/commit/"):
		return "commit"
	case strings.HasPrefix(normalizedPath, "internal/app/"):
		return "app"
	case strings.HasPrefix(normalizedPath, "internal/domain/"):
		return "core"
	case strings.HasPrefix(normalizedPath, "internal/infra/"):
		return "infra"
	default:
		baseName := strings.TrimSuffix(filepath.Base(normalizedPath), filepath.Ext(normalizedPath))
		return strings.ToLower(strings.ReplaceAll(baseName, "_", "-"))
	}
}
