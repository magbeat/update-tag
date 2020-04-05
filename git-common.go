package main

import (
	"github.com/blang/semver"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"strings"
)

func GetLatestTagFromRepository(repository *git.Repository) (semver.Version, error) {
	tagRefs, err := repository.Tags()
	CheckIfError(err)

	var latestVersion semver.Version
	latestVersion, err = semver.Make("0.0.0")
	err = tagRefs.ForEach(func(tagRef *plumbing.Reference) error {
		tagName := strings.Split(tagRef.Name().String(), "/")[2]

		var version semver.Version
		version, err = semver.Make(tagName)
		if latestVersion.LT(version) {
			latestVersion = version
		}
		return nil
	})
	CheckIfError(err)

	return latestVersion, nil
}
