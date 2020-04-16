package services

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/go-git/go-git/v5"
	"github.com/magbeat/update-tag/services/internal/common"
	utgit "github.com/magbeat/update-tag/services/internal/git"
	"github.com/magbeat/update-tag/services/internal/tag"
	"github.com/manifoldco/promptui"
	"os"
	"os/exec"
	"regexp"
)

func UpdateTag(vPrefix bool, prereleasePrefix string) {
	gitRoot, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	re := regexp.MustCompile(`\r?\n`)
	gitRootString := re.ReplaceAllString(string(gitRoot), "")
	common.CheckIfError(err)

	repo, err := git.PlainOpen(gitRootString)
	common.CheckIfError(err)

	headRef, err := repo.Head()
	common.CheckIfError(err)

	latestVersion, err := utgit.GetLatestTagFromRepository(repo)
	common.CheckIfError(err)

	var masterBranch = regexp.MustCompile(`master$`)
	var developmentBranch = regexp.MustCompile(`develop|feature`)

	if masterBranch.MatchString(headRef.String()) {
		stableTags, err := tag.GetStableTags(latestVersion)
		common.CheckIfError(err)
		runUpdate("Stable", latestVersion, stableTags, repo, vPrefix)
	} else if developmentBranch.MatchString(headRef.String()) {
		developmentTags, err := tag.GetDevelopmentTags(latestVersion, prereleasePrefix)
		common.CheckIfError(err)
		runUpdate("Development or feature branch", latestVersion, developmentTags, repo, vPrefix)
	} else {
		fmt.Println("Tagging only allowed from master, develop or feature branches")
	}
}

func runUpdate(branch string, tag *semver.Version, tags []*semver.Version, repo *git.Repository, vPrefix bool) {
	fmt.Println(fmt.Sprintf("%s branch found", branch))
	fmt.Println(fmt.Sprintf("Current tag: %s", tag.Original()))

	v := ""
	if vPrefix {
		v = "v"
	}

	tagsTemplates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   fmt.Sprintf("\U000025B8 {{\"%s\" | green}}{{ .Original | green }}", v),
		Inactive: fmt.Sprintf("  {{\"%s\" | cyan}}{{ .Original | cyan }}", v),
		Selected: "\U000025B8 {{ .Original | green }}",
	}

	prompt := promptui.Select{
		Label:     "Select next tag",
		Items:     tags,
		Templates: tagsTemplates,
	}
	_, result, err := prompt.Run()
	common.CheckIfError(err)

	head, err := repo.Head()
	common.CheckIfError(err)
	_, err = repo.CreateTag(result, head.Hash(), nil)

	pushTemplates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U000025B8 {{ . | green }}",
		Inactive: "  {{ . | cyan }}",
		Selected: "\U000025B8 {{ . | green }}",
	}

	prompt = promptui.Select{
		Label:     "Push to repo (with tags)?",
		Items:     []string{"Yes", "No"},
		Templates: pushTemplates,
	}

	_, result, err = prompt.Run()
	common.CheckIfError(err)

	if result == "Yes" {
		err = exec.Command("git", "push").Run()
		common.CheckIfError(err)
		err = exec.Command("git", "push", "--tags").Run()
		common.CheckIfError(err)
	}
	os.Exit(0)
}
