package app

import (
	"errors"
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
	Scope string
}

type CommitResult struct {
	Diff    string
	Message string
	Commit  domaincommit.Model
	Paths   []string
}

type CommitPlan struct {
	Result        CommitResult
	Preview       semantic.CommitPreview
	Quality       semantic.CommitQuality
	Context       semantic.CommitContext
	SemanticGroup string
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

func (service CommitService) PlanCommits(paths []string, options ...GenerateCommitOptions) (CommitReview, error) {
	commitOptions := firstCommitOptions(options)
	groupedPaths := map[string][]string{}

	for _, path := range paths {
		result, err := service.GenerateCommitForPaths([]string{path}, commitOptions)
		if err != nil {
			return CommitReview{}, err
		}

		context := semantic.NewCommitContext(result.Diff)
		groupKey := buildGroupKey(string(result.Commit.Type), result.Commit.Scope, semantic.BuildGroupingKey(string(result.Commit.Type), context))
		groupedPaths[groupKey] = append(groupedPaths[groupKey], path)
	}

	plans := []CommitPlan{}
	for _, group := range stableGroups(groupedPaths) {
		for _, chunk := range chunkPaths(group, 4) {
			result, err := service.GenerateCommitForPaths(chunk, commitOptions)
			if err != nil {
				return CommitReview{}, err
			}
			context := semantic.NewCommitContext(result.Diff)
			intent := semantic.DetectIntent(string(result.Commit.Type), context)
			plans = append(plans, CommitPlan{
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
	}

	review := CommitReview{
		Plans: plans,
	}
	review.Suggestions = buildSuggestions(review.Plans)

	return review, nil
}

func (service CommitService) ApplySuggestions(review CommitReview, options ...GenerateCommitOptions) (CommitReview, error) {
	mergedPaths := mergeSuggestedPlans(review.Plans)
	if len(mergedPaths) == 0 {
		return review, nil
	}

	commitOptions := firstCommitOptions(options)
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

func buildGroupKey(commitType string, scope string, semanticGroup string) string {
	return commitType + "|" + scope + "|" + semanticGroup
}

func stableGroups(groupedPaths map[string][]string) [][]string {
	keys := make([]string, 0, len(groupedPaths))
	for key := range groupedPaths {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	result := make([][]string, 0, len(keys))
	for _, key := range keys {
		group := append([]string(nil), groupedPaths[key]...)
		sort.Strings(group)
		result = append(result, group)
	}

	return result
}

func chunkPaths(paths []string, chunkSize int) [][]string {
	if len(paths) == 0 {
		return nil
	}

	chunks := [][]string{}
	for start := 0; start < len(paths); start += chunkSize {
		end := start + chunkSize
		if end > len(paths) {
			end = len(paths)
		}
		chunks = append(chunks, append([]string(nil), paths[start:end]...))
	}

	return chunks
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

	reducedPlans := mergeSuggestedPlans(plans)
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

func mergeSuggestedPlans(plans []CommitPlan) [][]string {
	if len(plans) == 0 {
		return nil
	}

	mergedPaths := [][]string{}
	currentGroup := append([]string(nil), plans[0].Result.Paths...)
	currentPlan := plans[0]

	for index := 1; index < len(plans); index++ {
		nextPlan := plans[index]
		if canMergePlans(currentPlan, nextPlan) && len(currentGroup)+len(nextPlan.Result.Paths) <= 4 {
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
