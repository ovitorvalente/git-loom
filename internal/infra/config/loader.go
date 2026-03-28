package config

import (
	"os"
	"strings"
)

func Load(path string) (Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}

		return Config{}, err
	}

	return parseConfig(string(content)), nil
}

func parseConfig(content string) Config {
	configuration := DefaultConfig()
	section := ""

	for _, rawLine := range strings.Split(content, "\n") {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasSuffix(line, ":") && !strings.Contains(line, " ") {
			section = strings.TrimSuffix(line, ":")
			continue
		}

		key, value, ok := splitKeyValue(line)
		if !ok {
			continue
		}

		switch section {
		case "commit":
			if key == "scope" {
				configuration.Commit.Scope = value
			}
		case "cli":
			if key == "auto_confirm" {
				configuration.CLI.AutoConfirm = value == "true"
			}
		}
	}

	return configuration
}

func splitKeyValue(line string) (string, string, bool) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}

	key := strings.TrimSpace(parts[0])
	value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
	if key == "" {
		return "", "", false
	}

	return key, value, true
}
