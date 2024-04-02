package git

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/king8fisher/git_puller/config"
	"github.com/king8fisher/git_puller/output"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/zzwx/jsonwalk"
)

func LatestCommit() (string, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.github.com/repos/"+config.Owner+"/"+config.Repo+"/commits?per_page=1", nil)
	req.Header.Set("Authorization", "Bearer "+config.Token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	var sb bytes.Buffer
	_, _ = io.Copy(&sb, res.Body)
	var f interface{}
	err = json.Unmarshal(sb.Bytes(), &f)
	if err != nil {
		return "", fmt.Errorf("can't decode json %v", string(sb.String()))
	}
	foundCommit := ""
	jsonwalk.Walk(&f, jsonwalk.Callback(func(path jsonwalk.WalkPath, key interface{}, value interface{}, vType jsonwalk.NodeValueType) {
		if path.Path() == "[0].commit.url" && vType == jsonwalk.String {
			// https://api.github.com/repos/<owner>/<repo>/git/commits/35d80512
			url := value.(string)
			parts := strings.Split(url, "/")
			if len(parts) > 0 {
				foundCommit = parts[len(parts)-1]
			}
		}
	}))
	return foundCommit, nil
}

func Clone(toPath string) {
	if stat, err := os.Stat(filepath.Join(toPath, ".git")); err == nil && stat.IsDir() {
		output.Info("git clone", "seems to be already cloned")
		return
	}

	cmd := exec.Command("git", "clone", config.RepoURL, toPath)
	out, err := cmd.CombinedOutput()
	output.Error("git clone", err)
	output.Info("git clone", string(out))
}

var lastPull = time.Time{}
var lastPullMutex = sync.Mutex{}
var lastCommit string

// PullCapped calls Pull in a capped manner.
func PullCapped(path string) string {
	if !lastPullMutex.TryLock() {
		return lastCommit
	}
	defer lastPullMutex.Unlock()

	if time.Since(lastPull).Seconds() >= config.CapSeconds {
		newLatestCommit, err := LatestCommit()
		doPull := false
		if err == nil {
			if lastCommit == "" {
				doPull = true
			} else {
				if lastCommit != newLatestCommit {
					doPull = true
				}
			}
		} else {
			doPull = true
		}
		lastCommit = newLatestCommit
		if doPull {
			_ = Pull(path)
			lastPull = time.Now()
			return lastCommit
		}
	}
	return lastCommit
}

var gitPullMutex = sync.Mutex{}

// Pull returns the output, running git pull within a Mutex.
func Pull(path string) string {
	gitPullMutex.Lock()
	defer gitPullMutex.Unlock()

	cmd := exec.Command("git", "pull", "--force")

	cmd.Dir = path

	o, err := cmd.CombinedOutput()
	output.Error("git pull", err)
	return string(o)
}
