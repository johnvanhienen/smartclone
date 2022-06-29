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
    url       string
    gitDomain string
    savePath  string
    cloneDir  string
}

func (r *Repo) fillDefaults() {
    r.cloneDir = fmt.Sprintf("%s/git", os.Getenv("HOME"))
}

func main() {
    versionFlag := flag.Bool("v", false, "Displays the version number of Smartclone and Go.")
    flag.Parse()
    if *versionFlag {
        fmt.Println(versionStr)
        os.Exit(0)
    }

    repo := Repo{}
    repo.fillDefaults()
    repo.url = flag.Arg(0)

    repo.gitDomain, repo.savePath = scrubUrl(repo.url)
    clonePath, err := createDirPath(repo)
    if err != nil {
        fmt.Println(err)
    }
    err = cloneRepo(clonePath, repo.url)
    if err != nil {
        fmt.Errorf("%s", err)
        err := cleanupPathArtifacts(clonePath)
        if err != nil {
            fmt.Errorf("%s", err)
        }
    }

}

func scrubUrl(url string) (gitDomain string, savePath string) {
    if strings.HasPrefix(url, "git@") {
        return scrubSshUrl(url)
    } else if strings.HasPrefix(url, "https://") {
        return scrubHttpsUrl(url)
    } else {
        fmt.Println("Please provide a url that starts with 'git@' or 'https://'")
        os.Exit(1)
    }
    return
}

func scrubSshUrl(originalUrl string) (gitDomain string, savePath string) {
    splittedUrl := strings.Split(originalUrl, ":")
    gitDomain = strings.TrimPrefix(splittedUrl[0], "git@")
    savePath = strings.TrimSuffix(splittedUrl[1], ".git")
    return gitDomain, savePath
}

func scrubHttpsUrl(originalUrl string) (gitDomain string, savePath string) {
    gitDomain = strings.Split(originalUrl, "/")[2]
    prefix := fmt.Sprintf("https://%s/", gitDomain)
    pathNoPrefix := strings.TrimPrefix(originalUrl, prefix)
    savePath = strings.TrimSuffix(pathNoPrefix, ".git")
    return gitDomain, savePath
}

func createDirPath(r Repo) (path string, err error) {
    path = fmt.Sprintf("%s/%s", r.cloneDir, r.savePath)

    if _, err := os.Stat(path); !os.IsNotExist(err) {
        pathWithSource := fmt.Sprintf("%s-%s", path, r.gitDomain)
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
    fmt.Printf("Cloned repository '%s' to '%s'\n", url, clonePath)
    return nil
}

func cleanupPathArtifacts(path string) (err error) {
    if _, err := os.Stat(path); !os.IsNotExist(err) {
        fmt.Printf("Removing clone artifacts..\n")
        err := os.RemoveAll(path)
        if err != nil {
            return err
        }
    }
    return
}
