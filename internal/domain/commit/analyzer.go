package commit

import (
	"path/filepath"
	"regexp"
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

type diffInsights struct {
	Commands  []string
	Flags     []string
	Functions []string
	Topic     string
}

var (
	flagPattern      = regexp.MustCompile(`--[a-z0-9][a-z0-9-]*`)
	cobraFlagPattern = regexp.MustCompile(`(?:BoolVar|StringVar|IntVar|DurationVar)\([^,]+,\s*"([^"]+)"`)
	commandPattern   = regexp.MustCompile(`Use:\s*"([^"]+)"`)
)

func AnalyzeDiff(diff string, commitType Type) Analysis {
	changes := extractChanges(diff)
	insights := extractDiffInsights(diff)
	scope := detectScope(changes)
	description := buildStructuredDescription(commitType, scope, changes, insights)
	body := buildStructuredBody(changes, insights)

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

func buildStructuredDescription(commitType Type, scope string, changes []Change, insights diffInsights) string {
	target := detectTarget(commitType, scope, changes)
	if insightTarget := detectInsightTarget(scope, insights); insightTarget != "" {
		target = insightTarget
	}
	action := detectAction(commitType, changes)

	if commitType == TypeTest && strings.HasPrefix(target, "testes de ") {
		return strings.TrimSpace("ajustar " + target)
	}

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
	case strings.Contains(path, "commit_feedback.go"):
		return "feedback de commit"
	case strings.Contains(path, "branch_service.go"):
		return "branch service"
	case strings.Contains(path, "workflow_service.go"):
		return "workflow service"
	case strings.Contains(path, "commit.go"):
		return "comando commit" //nolint:misspell
	case strings.Contains(path, "branch.go"):
		return "comando branch" //nolint:misspell
	case strings.Contains(path, "root.go"):
		return "comando raiz" //nolint:misspell
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
	case strings.Contains(path, "messages.go"):
		return "mensagens"
	case strings.Contains(path, "intent_detector.go"):
		return "detector de intencao"
	case strings.Contains(path, "scope_normalizer.go"):
		return "normalizador de escopo"
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

func buildStructuredBody(changes []Change, insights diffInsights) string {
	if len(changes) == 0 {
		return ""
	}

	details := make([]string, 0, len(changes))
	for _, change := range changes {
		details = append(details, "- "+describeChange(change))
	}

	if len(insights.Functions) > 0 {
		details = append(details, "- ajusta funcoes: "+strings.Join(limitItems(insights.Functions, 3), ", "))
	}
	if len(insights.Flags) > 0 {
		details = append(details, "- ajusta flags: "+strings.Join(limitItems(insights.Flags, 4), ", "))
	}
	if len(insights.Commands) > 0 {
		details = append(details, "- ajusta comandos: "+strings.Join(limitItems(insights.Commands, 3), ", "))
	}

	return strings.Join(details, "\n")
}

func extractDiffInsights(diff string) diffInsights {
	lines := strings.Split(diff, "\n")
	functions := map[string]bool{}
	flags := map[string]bool{}
	commands := map[string]bool{}
	topicVotes := map[string]int{}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if strings.HasPrefix(trimmed, "+") && !strings.HasPrefix(trimmed, "+++") {
			addedLine := strings.TrimSpace(strings.TrimPrefix(trimmed, "+"))
			if function := extractFunctionName(addedLine); function != "" {
				functions[function] = true
			}
			for _, match := range flagPattern.FindAllString(addedLine, -1) {
				flags[match] = true
			}
			if cobraFlag := extractCobraFlagName(addedLine); cobraFlag != "" {
				flags[cobraFlag] = true
			}
			if command := extractCommandName(addedLine); command != "" {
				commands[command] = true
			}
			collectTopicVotes(strings.ToLower(addedLine), topicVotes)
		}
	}

	return diffInsights{
		Functions: mapKeysSorted(functions),
		Flags:     mapKeysSorted(flags),
		Commands:  mapKeysSorted(commands),
		Topic:     dominantTopic(topicVotes),
	}
}

func extractFunctionName(line string) string {
	if !strings.HasPrefix(line, "func ") {
		return ""
	}

	withoutPrefix := strings.TrimPrefix(line, "func ")
	if strings.HasPrefix(withoutPrefix, "(") {
		receiverEnd := strings.Index(withoutPrefix, ")")
		if receiverEnd < 0 || receiverEnd+1 >= len(withoutPrefix) {
			return ""
		}
		withoutPrefix = strings.TrimSpace(withoutPrefix[receiverEnd+1:])
	}

	nameEnd := strings.Index(withoutPrefix, "(")
	if nameEnd <= 0 {
		return ""
	}

	return strings.TrimSpace(withoutPrefix[:nameEnd])
}

func extractCommandName(line string) string {
	matches := commandPattern.FindStringSubmatch(line)
	if len(matches) < 2 {
		return ""
	}

	command := strings.TrimSpace(matches[1])
	if command == "" {
		return ""
	}

	fields := strings.Fields(command)
	if len(fields) == 0 {
		return ""
	}

	return fields[0]
}

func extractCobraFlagName(line string) string {
	matches := cobraFlagPattern.FindStringSubmatch(line)
	if len(matches) < 2 {
		return ""
	}

	flag := strings.TrimSpace(matches[1])
	if flag == "" {
		return ""
	}

	return "--" + flag
}

func collectTopicVotes(line string, votes map[string]int) {
	switch {
	case containsAny(line, "json", "marshal", "unmarshal"):
		votes["saida json"]++
	case containsAny(line, "prompt", "confirm", "[y/n]", "stdin", "stdout"):
		votes["confirmacao do fluxo"]++
	case containsAny(line, "strict"):
		votes["modo estrito"]++
	case containsAny(line, "preview", "diff impact"):
		votes["preview de mudancas"]++
	case containsAny(line, "optimize", "suggestion", "sugest"):
		votes["sugestoes de agrupamento"]++
	case containsAny(line, "config", "yaml", "loader", "schema"):
		votes["configuracao"]++
	case containsAny(line, "doctor", "check", "diagnostic"):
		votes["diagnostico"]++
	case containsAny(line, "commit", "stage", "staged"):
		votes["fluxo de commit"]++
	case containsAny(line, "analyze", "review", "plan"):
		votes["analise de commits"]++
	}
}

func dominantTopic(votes map[string]int) string {
	bestTopic := ""
	bestVotes := 0
	for topic, value := range votes {
		if value > bestVotes || (value == bestVotes && topic < bestTopic) {
			bestTopic = topic
			bestVotes = value
		}
	}

	return bestTopic
}

func detectInsightTarget(scope string, insights diffInsights) string {
	if insights.Topic == "" {
		return ""
	}
	if scope == "" {
		return insights.Topic
	}
	if scope == "cli" {
		return insights.Topic + " do cli"
	}

	return insights.Topic + " em " + scope
}

func mapKeysSorted(values map[string]bool) []string {
	if len(values) == 0 {
		return nil
	}

	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func limitItems(values []string, limit int) []string {
	if len(values) <= limit {
		return values
	}
	return values[:limit]
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
