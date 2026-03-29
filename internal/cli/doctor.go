package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	infraai "github.com/ovitorvalente/git-loom/internal/infra/ai"
	infraconfig "github.com/ovitorvalente/git-loom/internal/infra/config"
	infragit "github.com/ovitorvalente/git-loom/internal/infra/git"
	"github.com/ovitorvalente/git-loom/internal/shared"
)

type doctorOptions struct {
	json bool
}

type doctorCheck struct {
	Name    string   `json:"name"`
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

type doctorReport struct {
	Status string        `json:"status"`
	Checks []doctorCheck `json:"checks"`
}

func newDoctorCommand() *cobra.Command {
	configuration, err := infraconfig.Load(".gitloom.yaml")
	if err != nil {
		configuration = infraconfig.DefaultConfig()
	}

	return newDoctorCommandWithDependencies(commitDependencies{
		gitRepository: infragit.NewRepository(),
		aiProvider:    infraai.NewNoopProvider(),
		config: commitConfig{
			DefaultScope: configuration.Commit.Scope,
			AutoConfirm:  configuration.CLI.AutoConfirm,
		},
	})
}

func newDoctorCommandWithDependencies(dependencies commitDependencies) *cobra.Command {
	options := doctorOptions{}
	command := &cobra.Command{
		Use:           "doctor",
		Short:         shared.MessageDoctorShort,
		Long:          doctorHelpText(),
		Example:       doctorExamples(),
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDoctorCommand(cmd, dependencies, options)
		},
	}

	command.Flags().BoolVar(&options.json, "json", false, shared.MessageJSONFlag)
	return command
}

func runDoctorCommand(command *cobra.Command, dependencies commitDependencies, options doctorOptions) error {
	report, err := buildDoctorReport(dependencies)
	if err != nil {
		return err
	}

	if options.json {
		payload, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(command.OutOrStdout(), string(payload))
		return err
	}

	lines := []string{"◆ doctor"}
	for _, check := range report.Checks {
		icon := "•"
		switch check.Status {
		case "ok":
			icon = "✔"
		case "warn":
			icon = "⚠"
		case "fail":
			icon = "✖"
		}
		lines = append(lines, fmt.Sprintf("%s %s: %s", icon, check.Name, check.Message))
		for _, detail := range check.Details {
			lines = append(lines, "  - "+detail)
		}
	}
	lines = append(lines, fmt.Sprintf("status geral: %s", report.Status))

	_, err = fmt.Fprintln(command.OutOrStdout(), strings.Join(lines, "\n"))
	return err
}

func buildDoctorReport(dependencies commitDependencies) (doctorReport, error) {
	checks := []doctorCheck{}

	isRepository, err := dependencies.gitRepository.IsRepository()
	if err != nil {
		return doctorReport{}, err
	}
	if !isRepository {
		checks = append(checks, doctorCheck{
			Name:    "repositorio",
			Status:  "fail",
			Message: "diretorio atual nao esta dentro de um repositorio git",
		})
		return doctorReport{Status: "fail", Checks: checks}, nil
	}
	checks = append(checks, doctorCheck{
		Name:    "repositorio",
		Status:  "ok",
		Message: "repositorio git detectado",
	})

	checks = append(checks, inspectConfigCheck())

	stagedPaths, err := dependencies.gitRepository.ListStagedFiles()
	if err != nil {
		return doctorReport{}, err
	}
	changedPaths, err := dependencies.gitRepository.ListChangedFiles()
	if err != nil {
		return doctorReport{}, err
	}
	partialPaths := partiallyStagedPaths(stagedPaths, changedPaths)

	checks = append(checks, buildWorkingTreeCheck(stagedPaths, changedPaths, partialPaths))
	checks = append(checks, buildPlanningCheck(stagedPaths, partialPaths))
	return doctorReport{
		Status: summarizeDoctorStatus(checks),
		Checks: checks,
	}, nil
}

func inspectConfigCheck() doctorCheck {
	if _, err := os.Stat(".gitloom.yaml"); err != nil {
		if os.IsNotExist(err) {
			return doctorCheck{
				Name:    "configuracao",
				Status:  "warn",
				Message: "arquivo .gitloom.yaml nao encontrado; o CLI usara defaults",
			}
		}

		return doctorCheck{
			Name:    "configuracao",
			Status:  "fail",
			Message: "nao foi possivel inspecionar .gitloom.yaml",
			Details: []string{err.Error()},
		}
	}

	configuration, err := infraconfig.Load(".gitloom.yaml")
	if err != nil {
		return doctorCheck{
			Name:    "configuracao",
			Status:  "fail",
			Message: "falha ao carregar .gitloom.yaml",
			Details: []string{err.Error()},
		}
	}

	details := []string{}
	if configuration.Commit.Scope != "" {
		details = append(details, "scope padrao: "+configuration.Commit.Scope)
	}
	if configuration.CLI.AutoConfirm {
		details = append(details, "auto_confirm: true")
	}

	message := "arquivo .gitloom.yaml carregado com sucesso"
	if len(details) == 0 {
		message = "arquivo .gitloom.yaml presente sem overrides ativos"
	}

	return doctorCheck{
		Name:    "configuracao",
		Status:  "ok",
		Message: message,
		Details: details,
	}
}

func buildWorkingTreeCheck(stagedPaths []string, changedPaths []string, partialPaths []string) doctorCheck {
	status := "ok"
	message := "working tree pronta para revisao"
	details := []string{
		fmt.Sprintf("staged: %d", len(stagedPaths)),
		fmt.Sprintf("changes: %d", len(changedPaths)),
	}

	if len(changedPaths) > 0 {
		status = "warn"
		message = "existem mudancas fora do stage que podem alterar o plano"
	}
	if len(partialPaths) > 0 {
		status = "fail"
		message = "existem arquivos parcialmente staged; o fluxo automatico nao suporta esse estado"
		details = append(details, "parcialmente staged: "+strings.Join(partialPaths, ", "))
	}

	return doctorCheck{
		Name:    "working-tree",
		Status:  status,
		Message: message,
		Details: details,
	}
}

func buildPlanningCheck(stagedPaths []string, partialPaths []string) doctorCheck {
	if len(partialPaths) > 0 {
		return doctorCheck{
			Name:    "planejamento",
			Status:  "fail",
			Message: "corrija arquivos parcialmente staged antes de usar gitloom commit",
		}
	}
	if len(stagedPaths) == 0 {
		return doctorCheck{
			Name:    "planejamento",
			Status:  "warn",
			Message: "nenhum arquivo staged; rode git add antes de gitloom commit",
		}
	}

	return doctorCheck{
		Name:    "planejamento",
		Status:  "ok",
		Message: fmt.Sprintf("%d arquivo(s) staged prontos para analise", len(stagedPaths)),
	}
}

func summarizeDoctorStatus(checks []doctorCheck) string {
	status := "ok"
	for _, check := range checks {
		if check.Status == "fail" {
			return "fail"
		}
		if check.Status == "warn" {
			status = "warn"
		}
	}
	return status
}

func partiallyStagedPaths(stagedPaths []string, changedPaths []string) []string {
	stagedSet := map[string]bool{}
	for _, path := range stagedPaths {
		stagedSet[path] = true
	}

	result := []string{}
	for _, path := range changedPaths {
		if stagedSet[path] {
			result = append(result, path)
		}
	}

	return result
}

func doctorHelpText() string {
	return `Valida se o repositorio e o ambiente atual estao prontos para usar o gitloom.

O comando inspeciona:
  - repositorio git
  - configuracao local
  - estado do stage e do working tree
  - bloqueios conhecidos do fluxo automatico`
}

func doctorExamples() string {
	return `  gitloom doctor
  gitloom doctor --json`
}
