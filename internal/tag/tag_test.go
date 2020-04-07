package tag

import (
	"github.com/blang/semver"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDevelopmentTags(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		currentVersion semver.Version
		newMajor       uint64
		newMinor       uint64
		newPre         semver.PRVersion
		want           []semver.Version
	}{
		{createVersion("0.0.0"), 1, 1, semver.PRVersion{VersionStr: "RC0"},
			[]semver.Version{createVersion("1.0.0-RC0"), createVersion("0.1.0-RC0")}},
		{createVersion("1.0.0-RC0"), 2, 1, semver.PRVersion{VersionStr: "RC1"},
			[]semver.Version{createVersion("1.0.0-RC1"), createVersion("1.1.0-RC0")}},
		{createVersion("1.1.0-RC0"), 2, 2, semver.PRVersion{VersionStr: "RC1"},
			[]semver.Version{createVersion("2.0.0-RC0"), createVersion("1.1.0-RC1")}},
	}

	for _, test := range tests {
		got, _ := GetDevelopmentTags(test.currentVersion, test.newMajor, test.newMinor, test.newPre)
		assert.Equal(got, test.want)
	}
}

func createVersion(versionString string) semver.Version {
	version, _ := semver.Parse(versionString)
	return version
}
