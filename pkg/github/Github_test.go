package github

import (
	"strings"
	"testing"
)

const topicToMatch = "TheTopic"

func TestRepoWithGivenTopicIsMatched(t *testing.T) {
	repo := RestRepo{
		Name:     "Nameson",
		Owner:    RepoOwner{"TheOwner"},
		Archived: false,
		Topics:   []string{topicToMatch},
	}
	if !repo.HasTopic(topicToMatch) {
		t.Errorf("repo has topic %s even though it is in denial", topicToMatch)
	}
}

func TestRepoWithoutGivenTopicIsNotMatched(t *testing.T) {
	repo := RestRepo{
		Name:     "Nameson",
		Owner:    RepoOwner{"TheOwner"},
		Archived: false,
		Topics:   []string{"AnotherTopic"},
	}
	if repo.HasTopic(topicToMatch) {
		t.Errorf("repo hasn't got topic %s even though it claims it does", topicToMatch)
	}
}

func TestTopicNamesAreCaseInsensitive(t *testing.T) {
	repo := RestRepo{
		Name:     "Nameson",
		Owner:    RepoOwner{"TheOwner"},
		Archived: false,
		Topics:   []string{strings.ToUpper(topicToMatch)},
	}
	if !repo.HasTopic(strings.ToLower(topicToMatch)) {
		t.Errorf("repo has topic %s even though it is in denial", topicToMatch)
	}
}
