package main

import (
    "testing"
)

func TestScrubUrl(t *testing.T) {
    t.Run("test ssh git url", func(t *testing.T) {
        url := "git@gitlab.com:org/project.git"
        desiredRepoStructure := Repo{
            url:       url,
            gitDomain: "gitlab.com",
            savePath:  "org/project",
        }

        repo := Repo{}
        repo.url = url
        repo.gitDomain, repo.savePath = scrubUrl(repo.url)

        got := repo
        want := desiredRepoStructure

        if got != want {
            t.Errorf("got %q want %q", got, want)
        }

    })
    t.Run("test https git url", func(t *testing.T) {
        url := "https://gitlab.com/org/project.git"
        desiredRepoStructure := Repo{
            url:       url,
            gitDomain: "gitlab.com",
            savePath:  "org/project",
        }

        repo := Repo{}
        repo.url = url
        repo.gitDomain, repo.savePath = scrubUrl(repo.url)

        got := repo
        want := desiredRepoStructure
        if got != want {
            t.Errorf("got %q want %q", got, want)
        }
    })
}
