package services

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/magbeat/update-tag/services/internal/common"
	utgit "github.com/magbeat/update-tag/services/internal/git"
	"github.com/magbeat/update-tag/services/internal/tag"
	"os"
	"os/exec"
	"regexp"
)

func UpdateTag(vPrefix bool, prereleasePrefix string, forceProdTags bool, forceDevTags bool) {
	gitRoot, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	re := regexp.MustCompile(`\r?\n`)
	gitRootString := re.ReplaceAllString(string(gitRoot), "")
	common.CheckIfError(err)

	repo, err := git.PlainOpen(gitRootString)
	common.CheckIfError(err)

	headRef, err := repo.Head()
	common.CheckIfError(err)

	latestVersion, err := utgit.GetLatestTagFromRepository(repo)
	common.CheckIfError(err)

	masterBranch := regexp.MustCompile(`master|main`)
	developmentBranch := regexp.MustCompile(`develop|feature`)
	isProdBranch := masterBranch.MatchString(headRef.String())
	isDevBranch := developmentBranch.MatchString(headRef.String())

	if (isProdBranch && !forceDevTags) || forceProdTags {
		stableTags, err := tag.GetStableTags(latestVersion)
		common.CheckIfError(err)
		runUpdate("Stable", latestVersion, stableTags, repo, vPrefix)
	} else if (isDevBranch && !forceProdTags) || forceDevTags {
		developmentTags, err := tag.GetDevelopmentTags(latestVersion, prereleasePrefix)
		common.CheckIfError(err)
		runUpdate("Development or feature branch", latestVersion, developmentTags, repo, vPrefix)
	} else {
		fmt.Println("Tagging only allowed from master, develop or feature branches")
	}
}

func runUpdate(branch string, tag *semver.Version, tags []*semver.Version, repo *git.Repository, vPrefix bool) {
	fmt.Printf("\nCurrent Branch: %s\n", branch)
	fmt.Printf("Current Tag   : %s\n\n", tag.Original())

	v := ""
	if vPrefix {
		v = "v"
	}

	var suggestedTags []string
	for _, tag := range tags {
		suggestedTags = append(suggestedTags, fmt.Sprintf("%s%s", v, tag.String()))
	}

	var qs = []*survey.Question{
		{
			Name:   "tag",
			Prompt: &survey.Select{Message: "Next tag:", Options: suggestedTags},
		},
		{
			Name:   "push",
			Prompt: &survey.Select{Message: "Push to remote (with tags):", Options: []string{"Yes", "No"}, Default: "Yes"},
		},
	}

	answers := struct {
		Tag  string
		Push string
	}{}

	err := survey.Ask(qs, &answers)
	common.CheckIfError(err)

	head, err := repo.Head()
	common.CheckIfError(err)
	_, err = repo.CreateTag(answers.Tag, head.Hash(), nil)

	if answers.Push == "Yes" {
		err = exec.Command("git", "push").Run()
		common.CheckIfError(err)
		err = exec.Command("git", "push", "--tags").Run()
		common.CheckIfError(err)
	}
	os.Exit(0)
}
