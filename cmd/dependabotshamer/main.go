package main

import (
	"dependabotshamer/pkg/github"
	"fmt"
	"log"
	"os"
)

func main() {
	githubToken := envOrDie("GITHUB_TOKEN")
	reposWithoutDepandabotAlerts, err := github.ReposWithDependabotAlertsDisabled("navikt", githubToken)
	if err != nil {
		log.Fatalf("unable to retrieve repos: %v", err)
	}
	for _, repo := range reposWithoutDepandabotAlerts {
		fmt.Printf("%s/%s\n", repo.Owner.Name, repo.Name)
	}
	fmt.Printf("Found %d repos with Depenadbot disabled", len(reposWithoutDepandabotAlerts))
}

func envOrDie(name string) string {
	value, found := os.LookupEnv(name)
	if !found {
		fmt.Printf("unable to find env var '%s', I'm useless without it\n", name)
		os.Exit(1)
	}
	return value
}
