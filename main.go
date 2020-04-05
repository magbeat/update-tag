package main

import (
	"fmt"
	"github.com/blang/semver"
	"gopkg.in/src-d/go-git.v4"
	"os"
	"strings"
)

func main() {
	currentDir, err := os.Getwd()
	CheckIfError(err)

	repo, err := git.PlainOpen(currentDir)
	CheckIfError(err)

	headRef, err := repo.Head()
	CheckIfError(err)

	latestTag, err := GetLatestTagFromRepository(repo)
	CheckIfError(err)

	currentTag, err := semver.Make("0.0.0-RC0")
	CheckIfError(err)

	if latestTag != "" {
		latestTag = strings.Split(latestTag, "/")[2]
		currentTag, err = semver.Parse(latestTag)
		CheckIfError(err)
	}

	major, minor, patch, pre, err := calculateNextTag(currentTag)

	if strings.HasSuffix(headRef.String(), "master") {
		stableTags, err := getStableTags(currentTag, major, minor, patch)
		CheckIfError(err)
		updateTag("Stable", currentTag, stableTags, repo)
	} else if strings.HasSuffix(headRef.String(), "develop") {
		developmentTags, err := getDevelopmentTags(currentTag, major, minor, pre)
		CheckIfError(err)
		updateTag("Development", currentTag, developmentTags, repo)
	} else {
		fmt.Println("Tagging only allowed from stable or development branch")
	}
}

func updateTag(branch string, tag semver.Version, tags []semver.Version, repo *git.Repository) {
	fmt.Println(fmt.Sprintf("%s branch found", branch))
	fmt.Println(fmt.Sprintf("Current tag: %s", tag))
	fmt.Println("\nPossible new Tags:")

	for index, newTag := range tags {
		fmt.Println(fmt.Sprintf(" %d) %s", index+1, newTag))
	}

	fmt.Print("Please choose a tag: ")
	var tagIndex int
	_, err := fmt.Scanln(&tagIndex)
	CheckIfError(err)
	if tagIndex > 0 && tagIndex <= len(tags) {
		// fmt.Println(tags[tagIndex - 1])
		head, err := repo.Head()
		CheckIfError(err)
		_, err = repo.CreateTag(tags[tagIndex-1].String(), head.Hash(), nil)
		CheckIfError(err)
	} else {
		fmt.Println("Index out of range")
		os.Exit(0)
	}

	fmt.Print("Push to repo (y/n): ")
	var push string
	_, err = fmt.Scanln(&push)
	CheckIfError(err)

	if push == "y" {
		fmt.Println(push)
		err = repo.Push(&git.PushOptions{
			RemoteName: "origin",
		})
		CheckIfError(err)
	}

	os.Exit(0)
}

func getDevelopmentTags(tag semver.Version, newMajor uint64, newMinor uint64, newPre semver.PRVersion) ([]semver.Version, error) {
	var versions []semver.Version
	var err error = nil

	var newMajorVersion semver.Version
	newMajorVersion, err = semver.Make(fmt.Sprintf("%d.%d.%d-%s", newMajor, tag.Minor, tag.Patch, newPre))
	var newMinorVersion semver.Version
	newMinorVersion, err = semver.Make(fmt.Sprintf("%d.%d.%d-%s", tag.Major, newMinor, tag.Patch, newPre))

	versions = append(versions, newMajorVersion)
	versions = append(versions, newMinorVersion)

	return versions, err
}

func getStableTags(tag semver.Version, newMajor uint64, newMinor uint64, newPatch uint64) ([]semver.Version, error) {
	var versions []semver.Version
	var err error = nil

	var newMajorVersion semver.Version
	newMajorVersion, err = semver.Make(fmt.Sprintf("%d.%d.%d", newMajor, tag.Minor, tag.Patch))
	var newMinorVersion semver.Version
	newMinorVersion, err = semver.Make(fmt.Sprintf("%d.%d.%d", tag.Major, newMinor, tag.Patch))
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
