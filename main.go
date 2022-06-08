package main

import (
    "flag"
    "fmt"
    "os"
    "os/exec"
    "runtime"
    "strings"
)

var (
    version    = "v0.1.0"
    goVersion  = runtime.Version()
    versionStr = fmt.Sprintf("Smartclone version %v\n%v", version, goVersion)
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
    versionFlag := flag.Bool("v", false, "Displays the version number of Smartclone and Go.")
    flag.Parse()

    if *versionFlag {
        fmt.Println(versionStr)
        os.Exit(0)
    }
    url := flag.Arg(0)

    ctx := context{}
    ctx.fillDefaults()
    repo := scrubUrl(url)
    clonePath, err := createDirPath(repo, ctx)
    if err != nil {
        fmt.Println(err)
    }
    err = cloneRepo(clonePath, url)
    if err != nil {
        fmt.Println(err)
        cleanupPathArtifacts(clonePath)
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
        createAltDir := "n"
        fmt.Printf("Path already exists. Do you want to create '%s' instead? y/N\t", pathWithSource)
        _, err := fmt.Scanf("%s", &createAltDir)
        if err != nil {
            return "", err
        }

        if strings.ToLower(createAltDir) == "n" || createAltDir == "" {
            fmt.Println("Abort cloning..")
            os.Exit(0)
        } else if strings.ToLower(createAltDir) == "y" {
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
        return fmt.Errorf("could not clone repository, error message: %s", err)
    }
    fmt.Printf("Cloned repository '%s' to '%s'", url, clonePath)
    return nil
}

func cleanupPathArtifacts(path string) {
    if _, err := os.Stat(path); !os.IsNotExist(err) {
        fmt.Printf("Removing clone artifacts..\n")
        err := os.RemoveAll(path)
        if err != nil {
            return
        }
    }
}
