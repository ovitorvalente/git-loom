package commit

import "strings"

func ClassifyCommit(diff string) Type {
	normalizedDiff := strings.ToLower(diff)
	changes := extractChanges(diff)

	if hasDocsFileSignal(changes) || hasDocsSignal(normalizedDiff) {
		return TypeDocs
	}
	if hasTestFileSignal(changes) || hasTestSignal(normalizedDiff) {
		return TypeTest
	}
	if hasChoreFileSignal(changes) {
		return TypeChore
	}
	if hasFixSignal(normalizedDiff) {
		return TypeFix
	}
	if hasFeatureFileSignal(normalizedDiff) {
		return TypeFeat
	}
	if hasRefactorFileSignal(normalizedDiff) {
		return TypeRefactor
	}
	if hasDocsSignal(normalizedDiff) {
		return TypeDocs
	}
	if hasTestSignal(normalizedDiff) {
		return TypeTest
	}
	if hasFeatureFileSignal(normalizedDiff) {
		return TypeFeat
	}
	if hasRefactorFileSignal(normalizedDiff) {
		return TypeRefactor
	}
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

func hasDocsFileSignal(changes []Change) bool {
	if len(changes) == 0 {
		return false
	}

	for _, change := range changes {
		path := strings.ToLower(change.Path)
		if path == "readme.md" || strings.HasSuffix(path, ".md") {
			return true
		}
	}

	return false
}

func hasTestFileSignal(changes []Change) bool {
	if len(changes) == 0 {
		return false
	}

	for _, change := range changes {
		if strings.HasSuffix(strings.ToLower(change.Path), "_test.go") {
			return true
		}
	}

	return false
}

func hasChoreFileSignal(changes []Change) bool {
	if len(changes) == 0 {
		return false
	}

	for _, change := range changes {
		path := strings.ToLower(change.Path)
		if path == "go.mod" || path == "go.sum" || path == "makefile" || path == ".gitloom.yaml" {
			return true
		}
	}

	return false
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

func hasFeatureFileSignal(content string) bool {
	return strings.Contains(content, "new file mode") && containsAny(content, ".go", ".ts", ".tsx", ".js", ".jsx")
}

func hasRefactorFileSignal(content string) bool {
	return strings.Contains(content, "diff --git") && containsAny(content, ".go", ".ts", ".tsx", ".js", ".jsx")
}

func hasDocsSignal(content string) bool {
	return containsAny(content, ".md", "readme", "docs", "document", "documentation")
}

func hasTestSignal(content string) bool {
	return containsAny(content, "_test.go", "test", "spec", "coverage")
}

func containsAny(content string, keywords ...string) bool {
	for _, keyword := range keywords {
		if strings.Contains(content, keyword) {
			return true
		}
	}

	return false
}
