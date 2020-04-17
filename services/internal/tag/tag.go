package tag

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"regexp"
	"strconv"
)

func GetDevelopmentTags(currentVersion *semver.Version, prereleasePrefix string) (versions []*semver.Version, err error) {

	rc0 := fmt.Sprintf("%s.0", prereleasePrefix)
	rcNext := nextPrereleaseTag(currentVersion.Prerelease(), prereleasePrefix)

	newMajorVersion, err := semver.NewVersion(currentVersion.String())
	newMinorVersion, err := semver.NewVersion(currentVersion.String())

	if len(currentVersion.Prerelease()) > 0 && currentVersion.Minor() == 0 {
		*newMajorVersion, err = newMajorVersion.SetPrerelease(rcNext)
	} else {
		*newMajorVersion = newMajorVersion.IncMajor()
		*newMajorVersion, err = newMajorVersion.SetPrerelease(rc0)
	}
	if len(currentVersion.Prerelease()) > 0 && currentVersion.Minor() != 0 {
		*newMinorVersion, err = newMinorVersion.SetPrerelease(rcNext)
	} else {
		*newMinorVersion = newMinorVersion.IncMinor()
		*newMinorVersion, err = newMinorVersion.SetPrerelease(rc0)
	}

	versions = append(versions, newMajorVersion)
	versions = append(versions, newMinorVersion)

	return versions, err
}

func GetStableTags(curr *semver.Version) (versions []*semver.Version, err error) {
	if len(curr.Prerelease()) > 0 {
		newVersion, _ := semver.NewVersion(curr.String())
		*newVersion, _ = newVersion.SetPrerelease("")
		versions = append(versions, newVersion)
	} else {
		newMajorVersion, _ := semver.NewVersion(fmt.Sprintf("%d.%d.%d", curr.IncMajor().Major(), 0, 0))
		newMinorVersion, _ := semver.NewVersion(fmt.Sprintf("%d.%d.%d", curr.Major(), curr.IncMinor().Minor(), 0))
		newPatchVersion, _ := semver.NewVersion(fmt.Sprintf("%d.%d.%d", curr.Major(), curr.Minor(), curr.IncPatch().Patch()))
		versions = append(versions, newMajorVersion)
		versions = append(versions, newMinorVersion)
		versions = append(versions, newPatchVersion)
	}

	return versions, err
}

func nextPrereleaseTag(prerelease string, prefix string) string {
	numberRegex := regexp.MustCompile(`[0-9]+`)
	versionNum, _ := strconv.Atoi(string(numberRegex.Find([]byte(prerelease))))
	versionNum++
	return fmt.Sprintf("%s.%d", prefix, versionNum)
}

