package services

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/abiosoft/ishell"
	"github.com/go-git/go-git/v5"
	"github.com/magbeat/update-tag/services/internal/common"
	utgit "github.com/magbeat/update-tag/services/internal/git"
	"github.com/magbeat/update-tag/services/internal/tag"
	"os"
	"os/exec"
	"regexp"
)

func UpdateTag(vPrefix bool, prereleasePrefix string) {
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

	var masterBranch = regexp.MustCompile(`master$`)
	var developmentBranch = regexp.MustCompile(`develop|feature`)

	if masterBranch.MatchString(headRef.String()) {
		stableTags, err := tag.GetStableTags(latestVersion)
		common.CheckIfError(err)
		runUpdate("Stable", latestVersion, stableTags, repo, vPrefix)
	} else if developmentBranch.MatchString(headRef.String()) {
		developmentTags, err := tag.GetDevelopmentTags(latestVersion, prereleasePrefix)
		common.CheckIfError(err)
		runUpdate("Development or feature branch", latestVersion, developmentTags, repo, vPrefix)
	} else {
		fmt.Println("Tagging only allowed from master, develop or feature branches")
	}
}

func runUpdate(branch string, tag *semver.Version, tags []*semver.Version, repo *git.Repository, vPrefix bool) {
	shell := ishell.New()

	infoString := fmt.Sprintf("branch: %s\ncurrent tag: %s\nnext tag:", branch, tag.Original())

	v := ""
	if vPrefix {
		v = "v"
	}

	var suggestedTags []string
	for _, tag := range tags {
		suggestedTags = append(suggestedTags, fmt.Sprintf("%s%s", v, tag.String()))
	}
	newTagChoice := shell.MultiChoice(suggestedTags, infoString)
	newTag := suggestedTags[newTagChoice]

	head, err := repo.Head()
	common.CheckIfError(err)
	_, err = repo.CreateTag(newTag, head.Hash(), nil)

	pushChoice := shell.MultiChoice([]string{"Yes", "No"}, "Push to repo (with tags)")

	if pushChoice == 0 {
		err = exec.Command("git", "push").Run()
		common.CheckIfError(err)
		err = exec.Command("git", "push", "--tags").Run()
		common.CheckIfError(err)
	}
	os.Exit(0)
}
