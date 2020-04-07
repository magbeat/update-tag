package main

import (
	"fmt"
	"github.com/blang/semver"
	"github.com/go-git/go-git/v5"
	"github.com/manifoldco/promptui"
	"os"
	"os/exec"
	"regexp"
)

func main() {
	currentDir, err := os.Getwd()
	CheckIfError(err)

	repo, err := git.PlainOpen(currentDir)
	CheckIfError(err)

	headRef, err := repo.Head()
	CheckIfError(err)

	latestVersion, err := GetLatestTagFromRepository(repo)
	CheckIfError(err)

	major, minor, patch, pre, err := calculateNextTag(latestVersion)

	var masterBranch = regexp.MustCompile(`master$`)
	var developmentBranch = regexp.MustCompile(`develop|feature`)

	if masterBranch.MatchString(headRef.String()) {
		stableTags, err := getStableTags(latestVersion, major, minor, patch)
		CheckIfError(err)
		updateTag("Stable", latestVersion, stableTags, repo)
	} else if developmentBranch.MatchString(headRef.String()) {
		developmentTags, err := getDevelopmentTags(latestVersion, major, minor, pre)
		CheckIfError(err)
		updateTag("Development or feature branch", latestVersion, developmentTags, repo)
	} else {
		fmt.Println("Tagging only allowed from master, develop or feature branches")
	}
}

func updateTag(branch string, tag semver.Version, tags []semver.Version, repo *git.Repository) {
	fmt.Println(fmt.Sprintf("%s branch found", branch))
	fmt.Println(fmt.Sprintf("Current tag: %s", tag))

	prompt := promptui.Select{
		Label: "Select Tag",
		Items: tags,
	}
	_, result, err := prompt.Run()
	CheckIfError(err)

	head, err := repo.Head()
	CheckIfError(err)
	_, err = repo.CreateTag(result, head.Hash(), nil)

	prompt = promptui.Select{
		Label: "Push to repo?",
		Items: []string{"Yes", "No"},
	}

	_, result, err = prompt.Run()
	CheckIfError(err)

	if result == "Yes" {
		push := exec.Command("git", "push")
		pushTags := exec.Command("git", "push", "--tags")
		err = push.Run()
		CheckIfError(err)
		err = pushTags.Run()
		CheckIfError(err)
	}
	os.Exit(0)
}

func getDevelopmentTags(tag semver.Version, newMajor uint64, newMinor uint64, newPre semver.PRVersion) ([]semver.Version, error) {
	var versions []semver.Version
	var err error = nil
	rc0 := "RC0"

	var newMajorVersion semver.Version
	if len(tag.Pre) > 0 && tag.Minor == 0 {
		newMajorVersion, err = semver.Make(fmt.Sprintf("%d.%d.%d-%s", tag.Major, 0, 0, newPre))
	} else {
		newMajorVersion, err = semver.Make(fmt.Sprintf("%d.%d.%d-%s", newMajor, 0, 0, rc0))
	}
	var newMinorVersion semver.Version
	if len(tag.Pre) > 0 && tag.Minor != 0 {
		newMinorVersion, err = semver.Make(fmt.Sprintf("%d.%d.%d-%s", tag.Major, tag.Minor, tag.Patch, newPre))
	} else {
		newMinorVersion, err = semver.Make(fmt.Sprintf("%d.%d.%d-%s", tag.Major, newMinor, tag.Patch, rc0))
	}

	versions = append(versions, newMajorVersion)
	versions = append(versions, newMinorVersion)

	return versions, err
}

func getStableTags(tag semver.Version, newMajor uint64, newMinor uint64, newPatch uint64) ([]semver.Version, error) {
	var versions []semver.Version
	var err error = nil

	var newMajorVersion semver.Version
	newMajorVersion, err = semver.Make(fmt.Sprintf("%d.%d.%d", newMajor, 0, 0))
	var newMinorVersion semver.Version
	newMinorVersion, err = semver.Make(fmt.Sprintf("%d.%d.%d", tag.Major, newMinor, 0))
	var newPatchVersion semver.Version
	newPatchVersion, err = semver.Make(fmt.Sprintf("%d.%d.%d", tag.Major, tag.Minor, newPatch))

	versions = append(versions, newMajorVersion)
	versions = append(versions, newMinorVersion)
	versions = append(versions, newPatchVersion)

	return versions, err
}

func calculateNextTag(tag semver.Version) (uint64, uint64, uint64, semver.PRVersion, error) {
	nextMajor := tag.Major + 1
	nextMinor := tag.Minor + 1
	nextPatch := tag.Patch + 1
	var newSuffix semver.PRVersion
	var err error = nil

	if len(tag.Pre) > 0 {
		nextNum := tag.Pre[0].VersionNum + 1
		newSuffix, err = semver.NewPRVersion(fmt.Sprintf("RC%d", nextNum))
	} else {
		newSuffix, err = semver.NewPRVersion("RC0")
	}

	return nextMajor, nextMinor, nextPatch, newSuffix, err
}
