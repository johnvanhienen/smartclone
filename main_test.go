package main

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestScrubUrl(t *testing.T) {
	t.Run("ssh git url", func(t *testing.T) {
		url := "git@gitlab.com:org/project.git"
		desiredRepoStructure := Repo{
			url:      url,
			savePath: "gitlab.com/org/project",
		}

		repo := Repo{}
		repo.url = url
		repo.savePath = scrubUrl(repo.url)

		got := repo
		want := desiredRepoStructure

		if got != want {
			t.Errorf("got %q want %q", got, want)
		}

	})
	t.Run("https git url", func(t *testing.T) {
		url := "https://gitlab.com/org/project.git"
		desiredRepoStructure := Repo{
			url:      url,
			savePath: "gitlab.com/org/project",
		}

		repo := Repo{}
		repo.url = url
		repo.savePath = scrubUrl(repo.url)

		got := repo
		want := desiredRepoStructure
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}
	})
}

func TestCleanupPathArtifacts(t *testing.T) {
	t.Run("Cleanup artifacts", func(t *testing.T) {
		path := t.TempDir()
		err := cleanupPathArtifacts(path)
		if err != nil {
			t.Errorf("failed with error: %s", err)
		}
	})
}

func TestCreateDirPath(t *testing.T) {
	t.Run("Clone https repository", func(t *testing.T) {
		tmpdir := t.TempDir()
		repo := Repo{}
		repo.fillDefaults()
		repo.url = "https://github.com/github/gitignore"
		repo.cloneDir = tmpdir
		repo.savePath = scrubUrl(repo.url)
		clonePath, err := createDirPath(repo)
		defer cleanupPathArtifacts(repo.cloneDir)

		if err != nil {
			t.Errorf("failed with error: %s", err)
		}

		got := clonePath
		want := fmt.Sprintf("%s/github.com/github/gitignore", tmpdir)
		if got != want {
			t.Errorf("got %q want %q", got, want)
		}

	})

}

func TestCloneRepo(t *testing.T) {
	t.Run("Clone https repository", func(t *testing.T) {
		tmpdir := t.TempDir()
		repo := Repo{}
		repo.fillDefaults()
		repo.url = "https://github.com/github/gitignore"
		repo.cloneDir = tmpdir
		repo.savePath = scrubUrl(repo.url)
		clonePath, _ := createDirPath(repo)
		defer cleanupPathArtifacts(repo.cloneDir)
		fmt.Println()
		err := cloneRepo(clonePath, repo.url)
		if err != nil {
			t.Errorf("failed with error: %s", err)
		}
		gotPath := fmt.Sprintf("%s/github.com/github/gitignore/Go.gitignore", tmpdir)
		if _, err := os.Stat(gotPath); errors.Is(err, os.ErrNotExist) {
			t.Errorf("File does not exist. Detailed error: %s", err)
		}
	})
}
