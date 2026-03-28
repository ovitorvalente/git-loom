package app

import (
	"errors"
	"reflect"
	"sort"
	"strings"

	domaincommit "github.com/ovitorvalente/git-loom/internal/domain/commit"
	"github.com/ovitorvalente/git-loom/internal/interfaces"
)

var ErrEmptyDiff = errors.New("nenhuma mudanca staged encontrada; execute git add antes de gitloom commit")

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
	Result CommitResult
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

func (service CommitService) PlanCommits(paths []string, options ...GenerateCommitOptions) ([]CommitPlan, error) {
	commitOptions := firstCommitOptions(options)
	groupedPaths := map[string][]string{}

	for _, path := range paths {
		result, err := service.GenerateCommitForPaths([]string{path}, commitOptions)
		if err != nil {
			return nil, err
		}

		groupKey := buildGroupKey(result.Commit.Type, result.Commit.Scope)
		groupedPaths[groupKey] = append(groupedPaths[groupKey], path)
	}

	plans := []CommitPlan{}
	for _, group := range stableGroups(groupedPaths) {
		for _, chunk := range chunkPaths(group, 4) {
			result, err := service.GenerateCommitForPaths(chunk, commitOptions)
			if err != nil {
				return nil, err
			}
			plans = append(plans, CommitPlan{
				Result: result,
			})
		}
	}

	return plans, nil
}

func (service CommitService) generateFromDiff(diff string, paths []string, options GenerateCommitOptions) (CommitResult, error) {
	if strings.TrimSpace(diff) == "" {
		return CommitResult{}, ErrEmptyDiff
	}

	model := buildCommitModel(diff, options)
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

func buildCommitModel(diff string, options GenerateCommitOptions) domaincommit.Model {
	analysis := domaincommit.AnalyzeDiff(diff, domaincommit.ClassifyCommit(diff))
	scope := strings.TrimSpace(options.Scope)
	if scope == "" {
		scope = analysis.Scope
	}

	return domaincommit.Model{
		Type:        domaincommit.ClassifyCommit(diff),
		Scope:       scope,
		Description: analysis.Description,
		Body:        analysis.Body,
	}
}

func firstCommitOptions(options []GenerateCommitOptions) GenerateCommitOptions {
	if len(options) == 0 {
		return GenerateCommitOptions{}
	}

	return options[0]
}

func buildGroupKey(commitType domaincommit.Type, scope string) string {
	return string(commitType) + "|" + scope
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
