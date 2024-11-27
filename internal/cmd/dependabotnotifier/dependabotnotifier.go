package dependabotnotifier

import (
	"context"
	"fmt"
	"os"

	"github.com/navikt/dependabotnotifier/internal/github"
	"github.com/navikt/dependabotnotifier/internal/naisapi"
	"github.com/navikt/dependabotnotifier/internal/slack"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

const (
	exitCodeSuccess = iota
	exitCodeEnvFileError
	exitCodeConfigError
	exitCodeLoggerError
	exitCodeRunError
)

func Run(ctx context.Context) {
	log := logrus.StandardLogger()
	log.SetFormatter(&logrus.JSONFormatter{})

	if err := loadEnvFile(log); err != nil {
		log.WithError(err).Errorf("error loading .env file")
		os.Exit(exitCodeEnvFileError)
	}

	cfg, err := newConfig(ctx)
	if err != nil {
		log.WithError(err).Errorf("error when loading config")
		os.Exit(exitCodeConfigError)
	}

	appLogger, err := newLogger(cfg.LogFormat, cfg.LogLevel)
	if err != nil {
		log.WithError(err).Errorf("creating application logger")
		os.Exit(exitCodeLoggerError)
	}

	if err := run(ctx, cfg, appLogger); err != nil {
		appLogger.WithError(err).Errorf("error in run()")
		os.Exit(exitCodeRunError)
	}

	os.Exit(exitCodeSuccess)
}

func run(ctx context.Context, cfg *config, log logrus.FieldLogger) error {
	eg, egCtx := errgroup.WithContext(ctx)
	var gitHubRepos []string
	eg.Go(func() error {
		var err error
		gitHubRepos, err = github.
			NewClient(cfg.GitHubApiToken, log.WithField("client", "GitHub")).
			ReposWithDependabotAlertsDisabled(egCtx)
		if err != nil {
			return fmt.Errorf("fetch GitHub repositories: %w", err)
		}
		return nil
	})

	var teamsForRepos naisapi.RepoTeams
	eg.Go(func() error {
		var err error
		teamsForRepos, err = naisapi.
			NewClient(cfg.NaisApiEndpoint, cfg.NaisApiToken, log.WithField("client", "NAIS API")).
			TeamsForRepos(egCtx)
		if err != nil {
			return fmt.Errorf("fetch NAIS teams: %w", err)
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	numRepos := 0
	numNotifications := 0

	log.Debugf("start sending notifications to Slack")
	slackClient := slack.NewClient(cfg.SlackApiToken, log.WithField("client", "Slack"))
	for _, repoName := range gitHubRepos {
		log := log.WithField("repo_name", repoName)

		teams, exists := teamsForRepos[repoName]
		if !exists || len(teams) == 0 {
			log.Warnf("no NAIS team found for repository, unable to notify")
			continue
		}

		numRepos++
		for _, team := range teams {
			log := log.WithFields(logrus.Fields{
				"team_slug":     team.Slug,
				"slack_channel": team.SlackChannel,
			})
			log.Infof("send Slack notification")

			if err := slackClient.SendMessage(ctx, team.SlackChannel, team.Slug, repoName); err != nil {
				log.WithError(err).Errorf("failed to send Slack notification")
			}
			numNotifications++
		}
	}

	log.WithFields(logrus.Fields{
		"num_repos":              numRepos,
		"num_notifications_sent": numNotifications,
	}).Infof("done sending notifications")
	return nil
}
