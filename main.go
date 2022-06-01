package main

import (
    "fmt"
    "os"
    "os/exec"
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
    url := "git@gitlab.com:jvanhienen/dotfiles.git"
    repo := scrubUrl(url)
    clonePath, err := createDirPath(repo, ctx)
    if err != nil {
        fmt.Println(err)
    }
    err = cloneRepo(clonePath, url)
    if err != nil {
        fmt.Println(err)
    }

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

func createDirPath(r Repo, ctx context) (path string, err error) {
    path = fmt.Sprintf("%s/%s", ctx.defaultPath, r.path)

    if _, err := os.Stat(path); !os.IsNotExist(err) {
        pathWithSource := fmt.Sprintf("%s-%s", path, r.source)
        defaultAnswer := "n"
        fmt.Printf("Path already exists. Do you want to create '%s' instead? y/N\t", pathWithSource)
        fmt.Scanf("%s", &defaultAnswer)

        if strings.ToLower(defaultAnswer) == "n" || defaultAnswer == "" {
            fmt.Println("Abort cloning..")
            os.Exit(0)
        } else if strings.ToLower(defaultAnswer) == "y" {
            path = pathWithSource
        }
    }
    err = os.MkdirAll(path, 0755)
    if err != nil {
        return "", err
    }
    return path, nil
}

func cloneRepo(clonePath string, url string) (err error) {
    cmd := exec.Command("git", "clone", url, clonePath)
    err = cmd.Run()
    if err != nil {
        return err
    }
    fmt.Printf("Cloned repository '%s' to '%s'", url, clonePath)
    return nil
}
