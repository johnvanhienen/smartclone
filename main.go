package main

import (
    "fmt"
    "os"
    "strings"
)

type Repo struct {
    source string
    path   string
}

type context struct {
    defaultPath string
}

func (c *context) fillDefaults() {
    c.defaultPath = fmt.Sprintf("%s/git", os.Getenv("HOME"))
}

func main() {
    ctx := context{}
    ctx.fillDefaults()
    fmt.Println(scrubUrl("https://gitlab.com/jvanhienen/dotfiles.git"))
    fmt.Println(scrubUrl("git@gitlab.com:jvanhienen/dotfiles.git"))
}

func scrubUrl(url string) (repo Repo) {
    if strings.HasPrefix(url, "git@") {
        repo = scrubSshUrl(url)
    } else if strings.HasPrefix(url, "https://") {
        repo = scrubHttpsUrl(url)
    } else {
        fmt.Println("Please provide a url that starts with 'git@' or 'https://'")
        os.Exit(1)
    }
    return repo
}

func scrubSshUrl(url string) (r Repo) {
    splitUrl := strings.Split(url, ":")
    r.source = strings.TrimPrefix(splitUrl[0], "git@")
    r.path = strings.TrimSuffix(splitUrl[1], ".git")
    return r
}

func scrubHttpsUrl(url string) (r Repo) {
    r.source = strings.Split(url, "/")[2]

    prefix := fmt.Sprintf("https://%s/", r.source)
    pathNoPrefix := strings.TrimPrefix(url, prefix)
    r.path = strings.TrimSuffix(pathNoPrefix, ".git")
    
    return r
}
