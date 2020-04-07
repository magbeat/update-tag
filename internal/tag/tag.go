package tag

import (
	"fmt"
	"github.com/blang/semver"
)

func GetDevelopmentTags(tag semver.Version, newMajor uint64, newMinor uint64, newPre semver.PRVersion) ([]semver.Version, error) {
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

func GetStableTags(tag semver.Version, newMajor uint64, newMinor uint64, newPatch uint64) ([]semver.Version, error) {
	var versions []semver.Version
	var err error = nil

	if len(tag.Pre) > 0 {
		newVersion, _ := semver.Make(fmt.Sprintf("%d.%d.%d", tag.Major, tag.Minor, tag.Patch))
		versions = append(versions, newVersion)
	} else {
		newMajorVersion, _ := semver.Make(fmt.Sprintf("%d.%d.%d", newMajor, 0, 0))
		newMinorVersion, _ := semver.Make(fmt.Sprintf("%d.%d.%d", tag.Major, newMinor, 0))
		newPatchVersion, _ := semver.Make(fmt.Sprintf("%d.%d.%d", tag.Major, tag.Minor, newPatch))
		versions = append(versions, newMajorVersion)
		versions = append(versions, newMinorVersion)
		versions = append(versions, newPatchVersion)
	}

	return versions, err
}

func CalculateNextTag(tag semver.Version) (uint64, uint64, uint64, semver.PRVersion, error) {
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
