package commit

import "strings"

func ClassifyCommit(diff string) Type {
	normalizedDiff := strings.ToLower(diff)
	switch {
	case containsAny(normalizedDiff, "fix", "bug", "error", "fail", "hotfix"):
		return TypeFix
	case containsAny(normalizedDiff, "refactor", "cleanup", "rename", "extract", "simplify"):
		return TypeRefactor
	case containsAny(normalizedDiff, "feat", "add", "create", "implement", "introduce"):
		return TypeFeat
	default:
		return TypeChore
	}
}

func containsAny(content string, keywords ...string) bool {
	for _, keyword := range keywords {
		if strings.Contains(content, keyword) {
			return true
		}
	}

	return false
}
