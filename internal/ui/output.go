package ui

import "github.com/ovitorvalente/git-loom/internal/app"

type OutputRenderer interface {
	CommitPlan(index, total int, plan app.CommitPlan) string
	ChangedFiles(staged, changed []string) string
	CommitSummary(summary CommitSummary) string
	Suggestions(suggestions []app.CommitSuggestion) string
	FinalPreview(plans []app.CommitPlan) string
}
