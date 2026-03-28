package config

type Config struct {
	Commit CommitConfig
	CLI    CLIConfig
}

type CommitConfig struct {
	Scope string
}

type CLIConfig struct {
	AutoConfirm bool
}

func DefaultConfig() Config {
	return Config{}
}
