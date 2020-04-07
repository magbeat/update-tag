package main

import (
	"fmt"
	"github.com/blang/semver"
	"github.com/go-git/go-git/v5"
	"github.com/manifoldco/promptui"
	"gitlab.com/novaloop-oss/update-tag/internal/common"
	"gitlab.com/novaloop-oss/update-tag/internal/gito"
	"gitlab.com/novaloop-oss/update-tag/internal/tag"
	"os"
	"os/exec"
	"regexp"
)

func main() {
	currentDir, err := os.Getwd()
	common.CheckIfError(err)

	repo, err := git.PlainOpen(currentDir)
	common.CheckIfError(err)

	headRef, err := repo.Head()
	common.CheckIfError(err)

	latestVersion, err := gito.GetLatestTagFromRepository(repo)
	common.CheckIfError(err)

	major, minor, patch, pre, err := tag.CalculateNextTag(latestVersion)

	var masterBranch = regexp.MustCompile(`master$`)
	var developmentBranch = regexp.MustCompile(`develop|feature`)

	if masterBranch.MatchString(headRef.String()) {
		stableTags, err := tag.GetStableTags(latestVersion, major, minor, patch)
		common.CheckIfError(err)
		updateTag("Stable", latestVersion, stableTags, repo)
	} else if developmentBranch.MatchString(headRef.String()) {
		developmentTags, err := tag.GetDevelopmentTags(latestVersion, major, minor, pre)
		common.CheckIfError(err)
		updateTag("Development or feature branch", latestVersion, developmentTags, repo)
	} else {
		fmt.Println("Tagging only allowed from master, develop or feature branches")
	}
}

func updateTag(branch string, tag semver.Version, tags []semver.Version, repo *git.Repository) {
	fmt.Println(fmt.Sprintf("%s branch found", branch))
	fmt.Println(fmt.Sprintf("Current tag: %s", tag))

	prompt := promptui.Select{
		Label: "Select next tag",
		Items: tags,
	}
	_, result, err := prompt.Run()
	common.CheckIfError(err)

	head, err := repo.Head()
	common.CheckIfError(err)
	_, err = repo.CreateTag(result, head.Hash(), nil)

	prompt = promptui.Select{
		Label: "Push to repo (with tags)?",
		Items: []string{"Yes", "No"},
	}

	_, result, err = prompt.Run()
	common.CheckIfError(err)

	if result == "Yes" {
		push := exec.Command("git", "push")
		pushTags := exec.Command("git", "push", "--tags")
		err = push.Run()
		common.CheckIfError(err)
		err = pushTags.Run()
		common.CheckIfError(err)
	}
	os.Exit(0)
}
