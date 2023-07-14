package main

import (
	"dependabotnotifier/pkg/github"
	"dependabotnotifier/pkg/slack"
	"dependabotnotifier/pkg/teams"
	"fmt"
	"log"
	"os"
)

func main() {
	githubToken := envOrDie("GITHUB_TOKEN")
	teamsToken := envOrDie("TEAMS_TOKEN")
	slackToken := envOrDie("SLACK_TOKEN")

	reposWithoutDependabotAlerts, err := github.ReposWithDependabotAlertsDisabled("navikt", githubToken)
	if err != nil {
		log.Fatalf("unable to retrieve repos: %v", err)
	}
	fmt.Printf("Found %d repos with Depenadbot disabled\n", len(reposWithoutDependabotAlerts))

	repoOwners := make(map[string][]teams.Team)
	var reposWithoutOwners []string
	for _, repo := range reposWithoutDependabotAlerts {
		repoFullName := fmt.Sprintf("%s/%s", repo.Owner.Name, repo.Name)
		admins, err := teams.AdminsFor(repoFullName, teamsToken)
		if err != nil {
			fmt.Printf("%v", err)
			os.Exit(1)
		}
		repoOwners[repoFullName] = admins
		if len(admins) == 0 {
			reposWithoutOwners = append(reposWithoutOwners, repo.Name)
		}
	}

	for repo, owners := range repoOwners {
		for _, owner := range owners {
			fmt.Printf("Notifying %s about %s in %s: ", owner.Slug, repo, owner.SlackChannel)
			heading := fmt.Sprintf(`:wave: *Hei, %s* :github2:`, owner.Slug)
			msg := fmt.Sprintf(`Dere er admins i GitHub-repoet *%s*. Dette repoet har ikke Dependabot alerts aktivert. Dependabot hjelper deg å oppdage biblioteker med kjente sårbarheter i appene dine. Du kan sjekke status og enable Dependabot <https://github.com/%s/security|her>. Hvis repoet ikke er i bruk, vurder å arkivere det. Det kan gjøres nederst på <https://github.com/%s/settings|denne siden>.`, repo, repo, repo)
			response := slack.SendMessage("kanalsomikkefinnes", heading, msg, slackToken)
			if response.Ok {
				fmt.Println("OK")
			} else {
				fmt.Printf("%s\n", response.ErrorMsg)
			}
		}
	}

	fmt.Printf("Repos without owners: %v\n", reposWithoutOwners)

	println("Done!")
}

func envOrDie(name string) string {
	value, found := os.LookupEnv(name)
	if !found {
		fmt.Printf("unable to find env var '%s', I'm useless without it\n", name)
		os.Exit(1)
	}
	return value
}
