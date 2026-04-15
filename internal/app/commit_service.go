package app

import (
	"errors"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"

	domaincommit "github.com/ovitorvalente/git-loom/internal/domain/commit"
	"github.com/ovitorvalente/git-loom/internal/interfaces"
	"github.com/ovitorvalente/git-loom/internal/semantic"
	"github.com/ovitorvalente/git-loom/internal/shared"
)

var ErrEmptyDiff = errors.New(shared.MessageEmptyDiff)

type CommitService struct {
	git interfaces.GitRepository
	ai  interfaces.AIProvider
}

type GenerateCommitOptions struct {
	Scope             string
	MaxFilesPerCommit int
}

type CommitResult struct {
	Diff    string
	Message string
	Commit  domaincommit.Model
	Paths   []string
}

type CommitPlan struct {
	SemanticGroup string
	Result        CommitResult
	Context       semantic.CommitContext
	Quality       semantic.CommitQuality
	Preview       semantic.CommitPreview
}

type CommitPreview = semantic.CommitPreview
type CommitQuality = semantic.CommitQuality

type CommitSuggestion struct {
	Message        string
	AutoApplicable bool
}

type CommitReview struct {
	Plans       []CommitPlan
	Suggestions []CommitSuggestion
}

func NewCommitService(gitRepository interfaces.GitRepository, aiProvider interfaces.AIProvider) CommitService {
	return CommitService{
		git: gitRepository,
		ai:  aiProvider,
	}
}

func (service CommitService) GenerateCommit(options ...GenerateCommitOptions) (CommitResult, error) {
	diff, err := service.git.GetDiff()
	if err != nil {
		return CommitResult{}, err
	}

	return service.generateFromDiff(diff, nil, firstCommitOptions(options))
}

func (service CommitService) GenerateCommitForPaths(paths []string, options ...GenerateCommitOptions) (CommitResult, error) {
	diff, err := service.git.GetDiff(paths...)
	if err != nil {
		return CommitResult{}, err
	}

	return service.generateFromDiff(diff, paths, firstCommitOptions(options))
}

func (service CommitService) BuildPlanForPaths(paths []string, options ...GenerateCommitOptions) (CommitPlan, error) {
	result, err := service.GenerateCommitForPaths(paths, options...)
	if err != nil {
		return CommitPlan{}, err
	}

	context := semantic.NewCommitContext(result.Diff)
	intent := semantic.DetectIntent(string(result.Commit.Type), context)
	return CommitPlan{
		Result:        result,
		Preview:       semantic.BuildPreview(context),
		Quality:       semantic.ScoreCommit(intent, context),
		Context:       context,
		SemanticGroup: semantic.BuildGroupingKey(string(result.Commit.Type), context),
	}, nil
}

func (service CommitService) PlanCommits(paths []string, options ...GenerateCommitOptions) (CommitReview, error) {
	commitOptions := firstCommitOptions(options)
	chunkSize := commitOptions.MaxFilesPerCommit
	if chunkSize <= 0 {
		chunkSize = 4
	}
	plans := []CommitPlan{}
	for _, group := range planningGroups(paths) {
		groupPlans := []CommitPlan{}
		for _, chunk := range chunkPaths(group, chunkSize) {
			result, err := service.GenerateCommitForPaths(chunk, commitOptions)
			if err != nil {
				return CommitReview{}, err
			}
			context := semantic.NewCommitContext(result.Diff)
			intent := semantic.DetectIntent(string(result.Commit.Type), context)
			groupPlans = append(groupPlans, CommitPlan{
				Result: CommitResult{
					Diff:    result.Diff,
					Message: result.Message,
					Commit:  result.Commit,
					Paths:   result.Paths,
				},
				Preview:       semantic.BuildPreview(context),
				Quality:       semantic.ScoreCommit(intent, context),
				Context:       context,
				SemanticGroup: semantic.BuildGroupingKey(string(result.Commit.Type), context),
			})
		}

		optimizedGroup, err := service.optimizePlans(groupPlans, commitOptions)
		if err != nil {
			return CommitReview{}, err
		}
		plans = append(plans, optimizedGroup...)
	}

	review := CommitReview{
		Plans: plans,
	}
	review.Suggestions = buildSuggestions(review.Plans)

	return review, nil
}

func (service CommitService) ApplySuggestions(review CommitReview, options ...GenerateCommitOptions) (CommitReview, error) {
	commitOptions := firstCommitOptions(options)
	maxFilesPerCommit := commitOptions.MaxFilesPerCommit
	if maxFilesPerCommit <= 0 {
		maxFilesPerCommit = 4
	}

	mergedPaths := mergeSuggestedPlans(review.Plans, maxFilesPerCommit)
	if len(mergedPaths) == 0 {
		return review, nil
	}

	plans := make([]CommitPlan, 0, len(mergedPaths))
	for _, group := range mergedPaths {
		result, err := service.GenerateCommitForPaths(group, commitOptions)
		if err != nil {
			return CommitReview{}, err
		}

		context := semantic.NewCommitContext(result.Diff)
		intent := semantic.DetectIntent(string(result.Commit.Type), context)
		plans = append(plans, CommitPlan{
			Result:        result,
			Preview:       semantic.BuildPreview(context),
			Quality:       semantic.ScoreCommit(intent, context),
			Context:       context,
			SemanticGroup: semantic.BuildGroupingKey(string(result.Commit.Type), context),
		})
	}

	plans, err := service.optimizePlans(plans, commitOptions)
	if err != nil {
		return CommitReview{}, err
	}

	optimizedReview := CommitReview{
		Plans: plans,
	}
	optimizedReview.Suggestions = buildSuggestions(optimizedReview.Plans)
	return optimizedReview, nil
}

func (service CommitService) generateFromDiff(diff string, paths []string, options GenerateCommitOptions) (CommitResult, error) {
	if strings.TrimSpace(diff) == "" {
		return CommitResult{}, ErrEmptyDiff
	}

	model := buildCommitModel(diff, paths, options)
	message, err := service.resolveMessage(diff, model)
	if err != nil {
		return CommitResult{}, err
	}

	return CommitResult{
		Diff:    diff,
		Message: message,
		Commit:  model,
		Paths:   append([]string(nil), paths...),
	}, nil
}

func (service CommitService) resolveMessage(diff string, model domaincommit.Model) (string, error) {
	if isNilAIProvider(service.ai) {
		return domaincommit.GenerateMessage(model)
	}

	message, err := service.ai.GenerateCommit(diff)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(message) != "" {
		return message, nil
	}

	return domaincommit.GenerateMessage(model)
}

func isNilAIProvider(provider interfaces.AIProvider) bool {
	if provider == nil {
		return true
	}

	providerValue := reflect.ValueOf(provider)
	if providerValue.Kind() != reflect.Ptr {
		return false
	}

	return providerValue.IsNil()
}

func buildCommitModel(diff string, paths []string, options GenerateCommitOptions) domaincommit.Model {
	commitType := domaincommit.ClassifyCommit(diff)
	analysis := domaincommit.AnalyzeDiff(diff, commitType)
	context := semantic.NewCommitContext(diff)
	intent := semantic.DetectIntent(string(commitType), context)
	scope := strings.TrimSpace(options.Scope)
	if scope == "" {
		scope = firstNonEmpty(intent.Scope, analysis.Scope)
	}

	return domaincommit.Model{
		Type:        commitType,
		Scope:       scope,
		Intent:      intent.Intent,
		Description: selectDescription(intent.Description, analysis.Description, scope, paths),
		Body:        analysis.Body,
	}
}

func firstCommitOptions(options []GenerateCommitOptions) GenerateCommitOptions {
	if len(options) == 0 {
		return GenerateCommitOptions{}
	}

	return options[0]
}

func planningGroups(paths []string) [][]string {
	dependencyPaths := []string{}
	regularPaths := []string{}

	for _, path := range paths {
		if isDependencyPath(path) {
			dependencyPaths = append(dependencyPaths, path)
			continue
		}

		regularPaths = append(regularPaths, path)
	}

	sort.Strings(dependencyPaths)
	sort.Strings(regularPaths)

	groups := [][]string{}
	if len(dependencyPaths) > 0 {
		groups = append(groups, dependencyPaths)
	}
	if len(regularPaths) > 0 {
		groups = append(groups, regularPaths)
	}

	return groups
}

func isDependencyPath(path string) bool {
	normalized := strings.ToLower(filepath.Base(strings.TrimSpace(path)))
	switch normalized {
	case "go.mod", "go.sum", "package.json", "package-lock.json", "pnpm-lock.yaml", "yarn.lock", "cargo.toml", "cargo.lock", "composer.json", "composer.lock", "requirements.txt", "poetry.lock", "pyproject.toml":
		return true
	default:
		return false
	}
}

func planningArea(path string) string {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return ""
	}

	directory := filepath.Dir(trimmed)
	if directory == "." || directory == "" {
		return "root"
	}

	segments := strings.Split(directory, "/")
	if len(segments) >= 2 {
		return segments[0] + "/" + segments[1]
	}

	return segments[0]
}

func chunkPaths(paths []string, chunkSize int) [][]string {
	if len(paths) == 0 {
		return nil
	}
	if chunkSize <= 1 {
		return [][]string{append([]string(nil), paths...)}
	}

	chunks := [][]string{}
	current := []string{}
	for _, path := range paths {
		if len(current) == 0 {
			current = append(current, path)
			continue
		}

		// Keep chunks cohesive: avoid forcing a mixed-area 4th file when a 3-file
		// chunk is already coherent.
		if len(current) >= chunkSize-1 && planningArea(current[0]) != planningArea(path) {
			chunks = append(chunks, append([]string(nil), current...))
			current = []string{path}
			continue
		}

		current = append(current, path)
		if len(current) == chunkSize {
			chunks = append(chunks, append([]string(nil), current...))
			current = current[:0]
		}
	}
	if len(current) > 0 {
		chunks = append(chunks, append([]string(nil), current...))
	}

	chunks = rebalanceSingleFileTail(chunks)
	return chunks
}

func rebalanceSingleFileTail(chunks [][]string) [][]string {
	if len(chunks) < 2 {
		return chunks
	}

	lastIndex := len(chunks) - 1
	previousIndex := lastIndex - 1
	if len(chunks[lastIndex]) != 1 {
		return chunks
	}
	if len(chunks[previousIndex]) <= 2 {
		return chunks
	}
	if !samePlanningArea(chunks[previousIndex], chunks[lastIndex]) {
		return chunks
	}

	moved := chunks[previousIndex][len(chunks[previousIndex])-1]
	chunks[previousIndex] = chunks[previousIndex][:len(chunks[previousIndex])-1]
	chunks[lastIndex] = append([]string{moved}, chunks[lastIndex]...)

	return chunks
}

func samePlanningArea(left []string, right []string) bool {
	if len(left) == 0 || len(right) == 0 {
		return false
	}

	return planningArea(left[0]) == planningArea(right[0])
}

func buildSuggestions(plans []CommitPlan) []CommitSuggestion {
	suggestions := []CommitSuggestion{}

	for index, plan := range plans {
		if plan.Quality.Score < 85 {
			suggestions = append(suggestions, CommitSuggestion{
				Message: "melhorar descricao do commit " + ordinal(index+1),
			})
		}
	}

	reducedPlans := mergeSuggestedPlans(plans, 4)
	if len(reducedPlans) > 0 && len(reducedPlans) < len(plans) {
		suggestions = append(suggestions, CommitSuggestion{
			Message:        "reduzir commits de " + ordinal(len(plans)) + " para " + ordinal(len(reducedPlans)),
			AutoApplicable: true,
		})
	}

	for index := 0; index < len(plans)-1; index++ {
		currentPlan := plans[index]
		nextPlan := plans[index+1]
		if canMergePlans(currentPlan, nextPlan) {
			suggestions = append(suggestions, CommitSuggestion{
				Message:        "agrupar commits " + ordinal(index+1) + " e " + ordinal(index+2),
				AutoApplicable: true,
			})
		}
	}

	return deduplicateSuggestions(suggestions)
}

func mergeSuggestedPlans(plans []CommitPlan, maxFilesPerCommit int) [][]string {
	if len(plans) == 0 {
		return nil
	}
	if maxFilesPerCommit <= 0 {
		maxFilesPerCommit = 4
	}

	mergedPaths := [][]string{}
	currentGroup := append([]string(nil), plans[0].Result.Paths...)
	currentPlan := plans[0]

	for index := 1; index < len(plans); index++ {
		nextPlan := plans[index]
		if canMergePlans(currentPlan, nextPlan) && len(currentGroup)+len(nextPlan.Result.Paths) <= maxFilesPerCommit {
			currentGroup = append(currentGroup, nextPlan.Result.Paths...)
			currentPlan = nextPlan
			continue
		}

		mergedPaths = append(mergedPaths, append([]string(nil), currentGroup...))
		currentGroup = append([]string(nil), nextPlan.Result.Paths...)
		currentPlan = nextPlan
	}

	mergedPaths = append(mergedPaths, append([]string(nil), currentGroup...))
	return mergedPaths
}

func canMergePlans(left CommitPlan, right CommitPlan) bool {
	if left.Result.Commit.Type != right.Result.Commit.Type {
		return false
	}
	if left.Result.Commit.Scope != right.Result.Commit.Scope {
		return false
	}

	return semantic.NormalizeScopeFromFiles(left.Context.Files) == semantic.NormalizeScopeFromFiles(right.Context.Files)
}

func (service CommitService) optimizePlans(plans []CommitPlan, options GenerateCommitOptions) ([]CommitPlan, error) {
	if len(plans) < 2 {
		return plans, nil
	}
	maxFilesPerCommit := options.MaxFilesPerCommit
	if maxFilesPerCommit <= 0 {
		maxFilesPerCommit = 4
	}

	groupedPaths := make([][]string, 0, len(plans))
	groupAnchors := make([]CommitPlan, 0, len(plans))
	for _, plan := range plans {
		if len(groupedPaths) == 0 {
			groupedPaths = append(groupedPaths, append([]string(nil), plan.Result.Paths...))
			groupAnchors = append(groupAnchors, plan)
			continue
		}

		lastIndex := len(groupedPaths) - 1
		lastGroup := groupedPaths[lastIndex]
		if shouldAttachSupportPlan(lastGroup, plan, groupAnchors[lastIndex], maxFilesPerCommit) {
			lastGroup = append(lastGroup, plan.Result.Paths...)
			sort.Strings(lastGroup)
			groupedPaths[lastIndex] = lastGroup
			continue
		}
		if shouldAttachSupportPlan(append([]string(nil), plan.Result.Paths...), groupAnchors[lastIndex], plan, maxFilesPerCommit) {
			mergedGroup := append(append([]string(nil), plan.Result.Paths...), lastGroup...)
			sort.Strings(mergedGroup)
			groupedPaths[lastIndex] = mergedGroup
			groupAnchors[lastIndex] = plan
			continue
		}

		groupedPaths = append(groupedPaths, append([]string(nil), plan.Result.Paths...))
		groupAnchors = append(groupAnchors, plan)
	}

	optimized := make([]CommitPlan, 0, len(groupedPaths))
	for _, group := range groupedPaths {
		result, err := service.GenerateCommitForPaths(group, options)
		if err != nil {
			return nil, err
		}

		context := semantic.NewCommitContext(result.Diff)
		intent := semantic.DetectIntent(string(result.Commit.Type), context)
		optimized = append(optimized, CommitPlan{
			Result:        result,
			Preview:       semantic.BuildPreview(context),
			Quality:       semantic.ScoreCommit(intent, context),
			Context:       context,
			SemanticGroup: semantic.BuildGroupingKey(string(result.Commit.Type), context),
		})
	}

	return optimized, nil
}

func shouldAttachSupportPlan(currentGroup []string, support CommitPlan, primary CommitPlan, maxFilesPerCommit int) bool {
	if len(currentGroup)+len(support.Result.Paths) > maxFilesPerCommit {
		return false
	}
	if !isSupportPlan(support) {
		return false
	}
	if isSupportPlan(primary) && !isPrimaryPlan(support) {
		return false
	}

	return samePlanArea(primary, support)
}

func isSupportPlan(plan CommitPlan) bool {
	if len(plan.Result.Paths) != 1 {
		return false
	}

	switch plan.Result.Commit.Type {
	case domaincommit.TypeTest, domaincommit.TypeDocs, domaincommit.TypeChore:
		return true
	default:
		return false
	}
}

func isPrimaryPlan(plan CommitPlan) bool {
	return !isSupportPlan(plan)
}

func samePlanArea(left CommitPlan, right CommitPlan) bool {
	return dominantArea(left.Result.Paths) != "" && dominantArea(left.Result.Paths) == dominantArea(right.Result.Paths)
}

func dominantArea(paths []string) string {
	if len(paths) == 0 {
		return ""
	}

	votes := map[string]int{}
	for _, path := range paths {
		votes[normalizeAreaPath(path)]++
	}

	selected := ""
	selectedVotes := 0
	for area, count := range votes {
		if count > selectedVotes || (count == selectedVotes && area < selected) {
			selected = area
			selectedVotes = count
		}
	}

	return selected
}

func normalizeAreaPath(path string) string {
	directory := filepath.Dir(path)
	base := filepath.Base(path)
	extension := filepath.Ext(base)
	name := strings.TrimSuffix(base, extension)
	name = strings.TrimSuffix(name, "_test")

	if directory == "." || directory == "" {
		return name
	}

	switch {
	case strings.HasPrefix(path, "internal/ui/"):
		return "internal/ui"
	case strings.HasPrefix(path, "internal/cli/"):
		return "internal/cli/" + name
	case strings.HasPrefix(path, "internal/shared/"):
		return "internal/shared"
	case strings.HasPrefix(path, "internal/semantic/"):
		return "internal/semantic"
	case strings.HasPrefix(path, "internal/domain/commit/"):
		return "internal/domain/commit/" + name
	default:
		return directory + "/" + name
	}
}

func deduplicateSuggestions(suggestions []CommitSuggestion) []CommitSuggestion {
	seenSuggestions := map[string]bool{}
	result := make([]CommitSuggestion, 0, len(suggestions))

	for _, suggestion := range suggestions {
		if suggestion.Message == "" || seenSuggestions[suggestion.Message] {
			continue
		}

		seenSuggestions[suggestion.Message] = true
		result = append(result, suggestion)
	}

	return result
}

func ordinal(value int) string {
	return strconv.Itoa(value)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}

	return ""
}

func describeFallback(paths []string) string {
	if len(paths) == 0 {
		return ""
	}

	return semantic.NormalizeScope(paths[0])
}

func selectDescription(intentDescription string, analysisDescription string, scope string, paths []string) string {
	intentDescription = strings.TrimSpace(intentDescription)
	analysisDescription = strings.TrimSpace(analysisDescription)

	switch {
	case intentDescription == "":
		return firstNonEmpty(analysisDescription, describeFallback(paths))
	case analysisDescription == "":
		return firstNonEmpty(intentDescription, describeFallback(paths))
	case isWeakDescription(intentDescription, scope):
		return analysisDescription
	case descriptionSpecificity(analysisDescription) > descriptionSpecificity(intentDescription)+1:
		return analysisDescription
	default:
		return intentDescription
	}
}

func isWeakDescription(description string, scope string) bool {
	normalizedDescription := strings.ToLower(strings.TrimSpace(description))
	normalizedScope := strings.ToLower(strings.TrimSpace(scope))

	if normalizedDescription == "" {
		return true
	}

	weakDescriptions := []string{
		"corrigir " + normalizedScope,
		"refinar " + normalizedScope,
		"atualizar " + normalizedScope,
		"ajustar " + normalizedScope,
		"adicionar " + normalizedScope,
		"cobrir " + normalizedScope,
		"ajustar testes de " + normalizedScope,
		"testes de " + normalizedScope,
		"ajustar testes de testes de " + normalizedScope,
	}

	for _, weakDescription := range weakDescriptions {
		if normalizedDescription == strings.TrimSpace(weakDescription) {
			return true
		}
	}

	return strings.Contains(normalizedDescription, "testes de testes de")
}

func descriptionSpecificity(description string) int {
	score := 0
	for _, token := range strings.Fields(strings.ToLower(strings.TrimSpace(description))) {
		if len(token) >= 4 {
			score++
		}
	}

	return score
}
