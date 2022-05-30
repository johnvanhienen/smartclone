package main

import (
    "testing"
)

func TestScrubUrl(t *testing.T) {
    t.Run("test ssh git url", func(t *testing.T) {
        repoStructure := Repo{
            source: "gitlab.com",
            path:   "user/project",
        }
        url := "git@gitlab.com:user/project.git"
        got := scrubUrl(url)
        want := repoStructure

        if got != want {
            t.Errorf("got %q want %q", got, want)
        }

    })
    t.Run("test https git url", func(t *testing.T) {
        repoStructure := Repo{
            source: "gitlab.com",
            path:   "user/project",
        }
        url := "https://gitlab.com/user/project.git"
        got := scrubUrl(url)
        want := repoStructure

        if got != want {
            t.Errorf("got %q want %q", got, want)
        }
    })
}
