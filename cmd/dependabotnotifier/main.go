package main

import (
	"dependabotshamer/pkg/slack"
	"fmt"
	"os"
)

func main() {
	//githubToken := envOrDie("GITHUB_TOKEN")
	//teamsToken := envOrDie("TEAMS_TOKEN")
	slackToken := envOrDie("SLACK_TOKEN")

	//reposWithoutDepandabotAlerts, err := github.ReposWithDependabotAlertsDisabled("navikt", githubToken)
	//if err != nil {
	//	log.Fatalf("unable to retrieve repos: %v", err)
	//}
	//fmt.Printf("Found %d repos with Depenadbot disabled", len(reposWithoutDepandabotAlerts))
	//
	//repoOwners := make(map[string][]teams.Team)
	//for _, repo := range reposWithoutDepandabotAlerts {
	//	repoFullName := fmt.Sprintf("%s/%s", repo.Owner.Name, repo.Name)
	//	teamsWithAdmin, err := teams.AdminsFor(repoFullName, teamsToken)
	//	if err != nil {
	//		fmt.Printf("%v", err)
	//		os.Exit(1)
	//	}
	//	repoOwners[repoFullName] = teamsWithAdmin
	//}
	//fmt.Printf("%v", repoOwners)

	err := slack.SendMessage("C05G9377NRL", "Hei 👋\nDette er en melding over flere \nlinjer.", slackToken)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	println("Sent!")
}

func envOrDie(name string) string {
	value, found := os.LookupEnv(name)
	if !found {
		fmt.Printf("unable to find env var '%s', I'm useless without it\n", name)
		os.Exit(1)
	}
	return value
}
