package app

import (
	"reflect"
	"strings"

	domaincommit "github.com/ovitorvalente/git-loom/internal/domain/commit"
	"github.com/ovitorvalente/git-loom/internal/interfaces"
)

type CommitService struct {
	git interfaces.GitRepository
	ai  interfaces.AIProvider
}

type CommitResult struct {
	Diff    string
	Message string
	Commit  domaincommit.Model
}

func NewCommitService(gitRepository interfaces.GitRepository, aiProvider interfaces.AIProvider) CommitService {
	return CommitService{
		git: gitRepository,
		ai:  aiProvider,
	}
}

func (service CommitService) GenerateCommit() (CommitResult, error) {
	diff, err := service.git.GetDiff()
	if err != nil {
		return CommitResult{}, err
	}

	model := buildCommitModel(diff)
	message, err := service.resolveMessage(diff, model)
	if err != nil {
		return CommitResult{}, err
	}

	return CommitResult{
		Diff:    diff,
		Message: message,
		Commit:  model,
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

func buildCommitModel(diff string) domaincommit.Model {
	return domaincommit.Model{
		Type:        domaincommit.ClassifyCommit(diff),
		Description: buildDescription(diff),
	}
}

func buildDescription(diff string) string {
	firstLine := extractFirstLine(diff)
	if firstLine != "" {
		return firstLine
	}

	return "update repository changes"
}

func extractFirstLine(diff string) string {
	for _, line := range strings.Split(diff, "\n") {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			return trimmedLine
		}
	}

	return ""
}
