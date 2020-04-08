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
		assert.Equal(test.want, got)
	}
}
func TestGetStableTags(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		currentVersion semver.Version
		newMajor       uint64
		newMinor       uint64
		newPatch       uint64
		want           []semver.Version
	}{
		{createVersion("0.0.0"), 1, 1, 1,
			[]semver.Version{createVersion("1.0.0"), createVersion("0.1.0"), createVersion("0.0.1")}},
		{createVersion("1.1.1"), 2, 2, 2,
			[]semver.Version{createVersion("2.0.0"), createVersion("1.2.0"), createVersion("1.1.2")}},
		{createVersion("1.0.0-RC2"), 2, 1, 1,
			[]semver.Version{createVersion("1.0.0")}},
		{createVersion("1.1.0-RC2"), 2, 2, 1,
			[]semver.Version{createVersion("1.1.0")}},
	}

	for _, test := range tests {
		got, _ := GetStableTags(test.currentVersion, test.newMajor, test.newMinor, test.newPatch)
		assert.Equal(test.want, got)
	}
}

func TestCalculateNextTag(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		currentVersion semver.Version
		newMajor       uint64
		newMinor       uint64
		newPatch       uint64
		newPre         semver.PRVersion
	}{
		{createVersion("0.0.0"), 1, 1, 1, semver.PRVersion{VersionStr: "RC0"}},
		{createVersion("1.0.0-RC2"), 2, 1, 1, semver.PRVersion{VersionStr: "RC3"}},
		{createVersion("2.2.2-RC0"), 3, 3, 3, semver.PRVersion{VersionStr: "RC1"}},
	}

	for _, test := range tests {
		newMajor, newMinor, newPatch, newPre, _ := CalculateNextTag(test.currentVersion)
		assert.Equal(test.newMajor, newMajor)
		assert.Equal(test.newMinor, newMinor)
		assert.Equal(test.newPatch, newPatch)
		assert.Equal(test.newPre, newPre)
	}
}

func createVersion(versionString string) semver.Version {
	version, _ := semver.Parse(versionString)
	return version
}
