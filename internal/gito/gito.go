package gito

import (
	"github.com/blang/semver"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	. "gitlab.com/novaloop-oss/update-tag/internal/common"
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
