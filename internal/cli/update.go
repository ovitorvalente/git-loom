package cli

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ovitorvalente/git-loom/internal/shared"
)

const (
	githubRepo = "ovitorvalente/git-loom"
	updateURL  = "https://raw.githubusercontent.com/%s/main/scripts/install.sh"
	apiURL     = "https://api.github.com/repos/%s/releases/latest"
)

type updateOptions struct {
	check   bool
	force   bool
	json    bool
	version string
}

type updateInfo struct {
	Current string `json:"current"`
	Latest  string `json:"latest"`
	Updated bool   `json:"updated"`
}

func newUpdateCommand() *cobra.Command {
	options := updateOptions{}

	command := &cobra.Command{
		Use:   "update",
		Short: shared.MessageUpdateShort,
		Long: `Atualiza o gitloom para a versao mais recente.
		
Busca a versao mais recente no GitHub e instala automaticamente.
O processo atualiza o binario no diretorio de installacao atual.`,
		Example: `  gitloom update
  gitloom update --check
  gl update --force`,
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate(cmd, options)
		},
	}

	command.Flags().BoolVar(&options.check, "check", false, "verifica se ha atualizacao disponivel sem instalar")
	command.Flags().BoolVar(&options.force, "force", false, "forca a atualizacao mesmo se ja estiver na versao mais recente")
	command.Flags().BoolVar(&options.json, "json", false, "saida em formato JSON")
	command.Flags().StringVar(&options.version, "version", "", "instalar versao especifica")

	return command
}

func runUpdate(cmd *cobra.Command, options updateOptions) error {
	currentVersion := strings.TrimPrefix(Version, "v")

	latestVersion, err := fetchLatestVersion()
	if err != nil {
		if options.json {
			return printJSONError(fmt.Sprintf("falha ao buscar versao: %v", err))
		}
		return fmt.Errorf("falha ao buscar versao mais recente: %w", err)
	}

	needsUpdate := options.force || currentVersion != latestVersion

	if options.check || !needsUpdate {
		info := updateInfo{
			Current: currentVersion,
			Latest:  latestVersion,
			Updated: !needsUpdate,
		}

		if options.json {
			return printJSONUpdate(info)
		}

		if info.Updated {
			fmt.Printf("gitloom %s ja esta atualizado\n", currentVersion)
			return nil
		}

		fmt.Printf("Atualizacao disponivel: %s -> %s\n", currentVersion, latestVersion)
		fmt.Print("Atualizando...\n")
	}

	if err := performUpdate(latestVersion); err != nil {
		if options.json {
			return printJSONError(fmt.Sprintf("falha ao atualizar: %v", err))
		}
		return fmt.Errorf("falha ao atualizar: %w", err)
	}

	info := updateInfo{
		Current: currentVersion,
		Latest:  latestVersion,
		Updated: true,
	}

	if options.json {
		return printJSONUpdate(info)
	}

	fmt.Printf("Atualizado para versao %s\n", latestVersion)
	return nil
}

func fetchLatestVersion() (string, error) {
	url := fmt.Sprintf(apiURL, githubRepo)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status=%d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	tag := extractTag(string(body))
	return strings.TrimPrefix(tag, "v"), nil
}

func extractTag(body string) string {
	marker := `"tag_name":"`
	start := strings.Index(body, marker)
	if start == -1 {
		return ""
	}
	start += len(marker)
	end := strings.Index(body[start:], `"`)
	if end == -1 {
		return ""
	}
	return body[start : start+end]
}

func performUpdate(version string) error {
	tmpDir, err := os.MkdirTemp("", "gitloom-update")
	if err != nil {
		return fmt.Errorf("falha ao criar dir temp: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	scriptPath := tmpDir + "/install.sh"

	url := fmt.Sprintf(updateURL, githubRepo)
	if err := downloadFile(url, scriptPath); err != nil {
		return fmt.Errorf("falha ao baixar script: %w", err)
	}

	if err := os.Chmod(scriptPath, 0o755); err != nil {
		return fmt.Errorf("falha ao definir permissao: %w", err)
	}

	installCmd := exec.Command("bash", scriptPath, "-v", version, "-f")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr

	return installCmd.Run()
}

func downloadFile(url, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status=%d", resp.StatusCode)
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func printJSONUpdate(info updateInfo) error {
	fmt.Fprintf(os.Stdout, "{\"current\":\"%s\",\"latest\":\"%s\",\"updated\":%v}\n",
		info.Current, info.Latest, info.Updated)
	return nil
}

func printJSONError(msg string) error {
	fmt.Fprintf(os.Stderr, "{\"error\":\"%s\"}\n", msg)
	return nil
}
