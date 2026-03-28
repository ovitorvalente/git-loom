package commit

import "strings"

func ClassifyCommit(diff string) Type {
	normalizedDiff := strings.ToLower(diff)
	if hasFixSignal(normalizedDiff) {
		return TypeFix
	}
	if hasRefactorSignal(normalizedDiff) {
		return TypeRefactor
	}
	if hasFeatureSignal(normalizedDiff) {
		return TypeFeat
	}

	return TypeChore
}

func hasFixSignal(content string) bool {
	switch {
	case containsAny(content, "fix", "bug", "error", "fail", "hotfix", "regression", "broken", "issue"):
		return true
	case containsAny(content, "remove", "delete", "revert") && containsAny(content, "bug", "error", "fail", "broken"):
		return true
	default:
		return false
	}
}

func hasRefactorSignal(content string) bool {
	return containsAny(content, "refactor", "cleanup", "rename", "extract", "simplify", "move", "reorganize")
}

func hasFeatureSignal(content string) bool {
	return containsAny(content, "feat", "add", "create", "implement", "introduce", "support", "enable")
}

func containsAny(content string, keywords ...string) bool {
	for _, keyword := range keywords {
		if strings.Contains(content, keyword) {
			return true
		}
	}

	return false
}
