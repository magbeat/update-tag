package gito

import (
	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	. "github.com/magbeat/update-tag/services/internal/common"
	"strings"
)

func GetLatestTagFromRepository(repository *git.Repository) (*semver.Version, error) {
	tagRefs, err := repository.Tags()
	CheckIfError(err)

	latestVersion, _ := semver.NewVersion("0.0.0")
	err = tagRefs.ForEach(func(tagRef *plumbing.Reference) error {
		tagName := strings.Split(tagRef.Name().String(), "/")[2]
		tagVersion, _ := semver.NewVersion(tagName)

		if tagVersion.GreaterThan(latestVersion) {
			latestVersion = tagVersion
		}
		return nil
	})
	CheckIfError(err)

	return latestVersion, nil
}
