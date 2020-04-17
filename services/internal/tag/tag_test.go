package tag

import (
	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDevelopmentTags(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		currentVersion   *semver.Version
		prereleasePrefix string
		want             []*semver.Version
	}{
		{createVersion("0.0.0"), "RC",
			[]*semver.Version{createVersion("1.0.0-RC.0"), createVersion("0.1.0-RC.0")}},
		{createVersion("1.0.0-RC.123"), "RC",
			[]*semver.Version{createVersion("1.0.0-RC.124"), createVersion("1.1.0-RC.0")}},
		{createVersion("1.0.0-RC.0"), "RC",
			[]*semver.Version{createVersion("1.0.0-RC.1"), createVersion("1.1.0-RC.0")}},
		{createVersion("1.1.0-RC.0"), "RC",
			[]*semver.Version{createVersion("2.0.0-RC.0"), createVersion("1.1.0-RC.1")}},
		{createVersion("1.0.0-TEST.11"), "TEST",
			[]*semver.Version{createVersion("1.0.0-TEST.12"), createVersion("1.1.0-TEST.0")}},
	}

	for _, test := range tests {
		got, _ := GetDevelopmentTags(test.currentVersion, test.prereleasePrefix)
		for index, tag := range got {
			assert.Equal(test.want[index].String(), tag.String())
		}
	}
}

func TestGetStableTags(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		currentVersion *semver.Version
		want           []*semver.Version
	}{
		{createVersion("v0.0.0"),
			[]*semver.Version{createVersion("1.0.0"), createVersion("0.1.0"), createVersion("0.0.1")}},
		{createVersion("v1.1.1"),
			[]*semver.Version{createVersion("2.0.0"), createVersion("1.2.0"), createVersion("1.1.2")}},
		{createVersion("v1.0.0-RC.2"),
			[]*semver.Version{createVersion("1.0.0")}},
		{createVersion("v1.1.0-RC.2"),
			[]*semver.Version{createVersion("1.1.0")}},
	}

	for _, test := range tests {
		got, _ := GetStableTags(test.currentVersion)
		for index, tag := range got {
			assert.Equal(test.want[index].String(), tag.String())
		}
	}
}

func createVersion(versionString string) *semver.Version {
	version, _ := semver.NewVersion(versionString)
	return version
}
