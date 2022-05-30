package main

import (
    "fmt"
    "strings"
)

func main() {
    fmt.Println("vim-go")
}

func scrubUrl(url string) (cleanUrl string) {
    if strings.HasPrefix(url, "git@") {
        cleanUrl = scrubSshUrl(url)
    }

    // cleanUrl = "jvanhienen/dotfile"
    return cleanUrl
}

func scrubSshUrl(url string) string {
    splitByColon := strings.SplitAfter(url, ":")
    noPrefixUrl := splitByColon[len(splitByColon)-1]
    noSuffixUrl := strings.TrimSuffix(noPrefixUrl, ".git")

    return noSuffixUrl
}
