package main

import (
    "testing"
)

func TestScrubUrl(t *testing.T) {
    t.Run("test ssh git url", func(t *testing.T) {
        url := "git@gitlab.com:jvanhienen/dotfiles.git"
        got := scrubUrl(url)
        want := "jvanhienen/dotfiles"

        if got != want {
            t.Errorf("got %q want %q", got, want)
        }

    })
    // https://gitlab.com/jvanhienen/dotfiles.git
}
