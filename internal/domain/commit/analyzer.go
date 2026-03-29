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
	target := detectTarget(commitType, scope, changes)
	action := detectAction(commitType, changes)

	return strings.TrimSpace(action + " " + target)
}

func detectTarget(commitType Type, scope string, changes []Change) string {
	if len(changes) == 1 {
		return normalizeDescriptionTarget(commitType, changes[0].Path)
	}
	if scope != "" {
		return scope
	}

	return "repositorio"
}

func detectAction(commitType Type, changes []Change) string {
	if hasOnlyStatus(changes, "adicionado") {
		switch commitType {
		case TypeFeat:
			return "adicionar"
		case TypeDocs:
			return "documentar"
		case TypeTest:
			return "adicionar testes para"
		default:
			return "adicionar"
		}
	}

	if hasOnlyStatus(changes, "atualizado") {
		switch commitType {
		case TypeFix:
			return "corrigir"
		case TypeRefactor:
			return "refatorar"
		case TypeDocs:
			return "atualizar"
		case TypeTest:
			return "ajustar testes de"
		default:
			return "atualizar"
		}
	}

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

func hasOnlyStatus(changes []Change, status string) bool {
	if len(changes) == 0 {
		return false
	}

	for _, change := range changes {
		if change.Status != status {
			return false
		}
	}

	return true
}

func normalizeFileTarget(path string) string {
	if target := detectTestTarget(path); target != "" {
		return target
	}

	if target := detectSpecificTarget(path); target != "" {
		return target
	}

	baseName := filepath.Base(path)
	nameWithoutExtension := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	return normalizeName(nameWithoutExtension)
}

func normalizeDescriptionTarget(commitType Type, path string) string {
	if commitType == TypeTest {
		if target := detectTestSubject(path); target != "" {
			return target
		}
	}

	return normalizeFileTarget(path)
}

func detectTestTarget(path string) string {
	subject := detectTestSubject(path)
	if subject == "" {
		return ""
	}

	return "testes de " + subject
}

func detectTestSubject(path string) string {
	if !strings.HasSuffix(path, "_test.go") {
		return ""
	}

	baseName := filepath.Base(path)
	nameWithoutSuffix := strings.TrimSuffix(baseName, "_test.go")
	if target := detectSpecificTarget(nameWithoutSuffix + ".go"); target != "" {
		return target
	}

	return normalizeName(nameWithoutSuffix)
}

func detectSpecificTarget(path string) string {
	switch {
	case strings.Contains(path, "commit_service.go"):
		return "commit service"
	case strings.Contains(path, "branch_service.go"):
		return "branch service"
	case strings.Contains(path, "workflow_service.go"):
		return "workflow service"
	case strings.Contains(path, "commit.go"):
		return "comando commit"
	case strings.Contains(path, "branch.go"):
		return "comando branch"
	case strings.Contains(path, "root.go"):
		return "comando raiz"
	case strings.Contains(path, "loader.go"):
		return "loader de configuracao"
	case strings.Contains(path, "schema.go"):
		return "schema de configuracao"
	case strings.Contains(path, "repository.go"):
		return "repositorio"
	case strings.Contains(path, "prompts.go"):
		return "prompts"
	case strings.Contains(path, "renderer.go"):
		return "renderer"
	case strings.Contains(path, "commit_view.go"):
		return "visao de commit"
	case strings.Contains(path, "summary_view.go"):
		return "resumo"
	case strings.Contains(path, "output.go"):
		return "output"
	case strings.Contains(path, "_test.go"):
		return "testes"
	default:
		return ""
	}
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
		details = append(details, "- "+describeChange(change))
	}

	return strings.Join(details, "\n")
}

func describeChange(change Change) string {
	target := normalizeFileTarget(change.Path)
	context := normalizeChangeContext(change.Path)

	switch change.Status {
	case "adicionado":
		return "adiciona " + target + context
	case "removido":
		return "remove " + target + context
	default:
		return "atualiza " + target + context
	}
}

func normalizeChangeContext(path string) string {
	directory := filepath.Dir(path)
	if directory == "." || directory == "" {
		return ""
	}

	segments := strings.Split(directory, "/")
	if len(segments) > 0 {
		switch segments[0] {
		case "internal", "cmd", "pkg":
			segments = segments[1:]
		}
	}
	if len(segments) > 2 {
		segments = segments[len(segments)-2:]
	}

	return " em " + strings.Join(segments, "/")
}
