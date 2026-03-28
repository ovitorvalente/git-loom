package commit

import (
	"path/filepath"
	"sort"
	"strings"
)

type Analysis struct {
	Scope       string
	Description string
	Body        string
}

type Change struct {
	Path   string
	Status string
}

func AnalyzeDiff(diff string, commitType Type) Analysis {
	changes := extractChanges(diff)
	scope := detectScope(changes)
	description := buildStructuredDescription(commitType, scope, changes)
	body := buildStructuredBody(changes)

	return Analysis{
		Scope:       scope,
		Description: description,
		Body:        body,
	}
}

func extractChanges(diff string) []Change {
	lines := strings.Split(diff, "\n")
	changes := []Change{}

	for index := 0; index < len(lines); index++ {
		line := strings.TrimSpace(lines[index])
		if !strings.HasPrefix(line, "diff --git ") {
			continue
		}

		path := extractPath(line)
		status := detectStatus(lines, index)
		changes = append(changes, Change{
			Path:   path,
			Status: status,
		})
	}

	return deduplicateChanges(changes)
}

func extractPath(line string) string {
	fields := strings.Fields(line)
	if len(fields) < 4 {
		return ""
	}

	return strings.TrimPrefix(fields[3], "b/")
}

func detectStatus(lines []string, start int) string {
	if start+1 >= len(lines) {
		return "modificado"
	}

	nextLine := strings.TrimSpace(lines[start+1])
	switch {
	case strings.HasPrefix(nextLine, "new file mode"):
		return "adicionado"
	case strings.HasPrefix(nextLine, "deleted file mode"):
		return "removido"
	default:
		return "atualizado"
	}
}

func deduplicateChanges(changes []Change) []Change {
	uniqueChanges := map[string]Change{}

	for _, change := range changes {
		if change.Path == "" {
			continue
		}
		uniqueChanges[change.Path] = change
	}

	paths := make([]string, 0, len(uniqueChanges))
	for path := range uniqueChanges {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	result := make([]Change, 0, len(paths))
	for _, path := range paths {
		result = append(result, uniqueChanges[path])
	}

	return result
}

func detectScope(changes []Change) string {
	scopeVotes := map[string]int{}

	for _, change := range changes {
		scope := detectPathScope(change.Path)
		if scope == "" {
			continue
		}
		scopeVotes[scope]++
	}

	return mostCommonScope(scopeVotes)
}

func detectPathScope(path string) string {
	segments := strings.Split(path, "/")
	if len(segments) == 0 {
		return ""
	}

	switch {
	case len(segments) >= 3 && segments[0] == "internal" && segments[1] == "domain":
		return segments[2]
	case len(segments) >= 3 && segments[0] == "internal" && segments[1] == "infra":
		return segments[2]
	case len(segments) >= 2 && segments[0] == "internal":
		return segments[1]
	case len(segments) >= 2 && segments[0] == "cmd":
		return segments[1]
	case len(segments) >= 2 && segments[0] == "pkg":
		return segments[1]
	case len(segments) >= 1:
		return normalizeName(segments[0])
	default:
		return ""
	}
}

func mostCommonScope(scopeVotes map[string]int) string {
	selectedScope := ""
	selectedVotes := 0

	for scope, votes := range scopeVotes {
		if votes > selectedVotes || (votes == selectedVotes && scope < selectedScope) {
			selectedScope = scope
			selectedVotes = votes
		}
	}

	return selectedScope
}

func buildStructuredDescription(commitType Type, scope string, changes []Change) string {
	target := detectTarget(scope, changes)
	action := detectAction(commitType)

	return strings.TrimSpace(action + " " + target)
}

func detectTarget(scope string, changes []Change) string {
	if len(changes) == 1 {
		return normalizeFileTarget(changes[0].Path)
	}
	if scope != "" {
		return scope
	}

	return "repositorio"
}

func detectAction(commitType Type) string {
	switch commitType {
	case TypeFeat:
		return "adicionar"
	case TypeFix:
		return "corrigir"
	case TypeRefactor:
		return "refatorar"
	case TypeDocs:
		return "documentar"
	case TypeTest:
		return "cobrir"
	default:
		return "atualizar"
	}
}

func normalizeFileTarget(path string) string {
	baseName := filepath.Base(path)
	nameWithoutExtension := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	return normalizeName(nameWithoutExtension)
}

func normalizeName(name string) string {
	replacer := strings.NewReplacer("_", " ", "-", " ", ".", " ")
	normalizedName := replacer.Replace(name)
	return strings.TrimSpace(normalizedName)
}

func buildStructuredBody(changes []Change) string {
	if len(changes) == 0 {
		return ""
	}

	details := make([]string, 0, len(changes))
	for _, change := range changes {
		target := normalizeFileTarget(change.Path)
		details = append(details, "- "+change.Status+" "+target)
	}

	return strings.Join(details, "\n")
}
