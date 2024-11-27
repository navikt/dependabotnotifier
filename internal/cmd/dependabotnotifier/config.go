package dependabotnotifier

import (
	"context"
	"fmt"

	"github.com/sethvargo/go-envconfig"
)

type config struct {
	GitHubApiToken  string `env:"GITHUB_TOKEN,required"`
	NaisApiToken    string `env:"TEAMS_TOKEN,required"`
	NaisApiEndpoint string `env:"NAIS_API_ENDPOINT,default=https://console.nav.cloud.nais.io/graphql"`
	SlackApiToken   string `env:"SLACK_TOKEN,required"`
	LogFormat       string `env:"LOG_FORMAT,default=json"`
	LogLevel        string `env:"LOG_LEVEL,default=info"`
}

func newConfig(ctx context.Context) (*config, error) {
	cfg := &config{}
	if err := envconfig.Process(ctx, cfg); err != nil {
		return nil, err
	}

	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func validateConfig(cfg *config) error {
	if cfg.GitHubApiToken == "" {
		return fmt.Errorf("missing GitHub API token")
	}

	if cfg.NaisApiToken == "" {
		return fmt.Errorf("missing NAIS API token")
	}

	if cfg.SlackApiToken == "" {
		return fmt.Errorf("missing Slack API token")
	}

	return nil
}
