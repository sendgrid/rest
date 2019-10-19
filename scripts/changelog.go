package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const (
	prefix   = "v"
	branch   = "master"
	prTitle  = "Merge pull request #(\\d+) from (\\w*)"
	clFormat = `## [{{ .Version }}] - {{ .Date }}
### Merged
{{- range .PullRequests }}
- Pull #{{ .Number }}: {{ .Commit.Body }}
{{- end }}
{{ if .Authors -}}
- Special thanks for these PRs to:
{{- range .Authors }}
    - [{{ . }}](//github.com/{{ . }})
{{- end }}
{{- else -}}
- No PRs were merged in this release.
{{- end }}
`
)

var (
	prRegexp   = regexp.MustCompile(prTitle)
	clTemplate = template.Must(template.New("changelog").Parse(clFormat))
)

type commit struct {
	Hash  string
	Title string
	Body  string
}

type pullRequest struct {
	Number int
	Author string
	Commit commit
}

func git(args ...string) (io.Reader, error) {
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(out), nil
}

func readLine(r io.Reader) (string, error) {
	br := bufio.NewReader(r)
	line, _, err := br.ReadLine()
	if err != nil {
		return "", err
	}
	return string(line), nil
}

// Gets the latest version tag in master
func latestTag() (string, error) {
	// Run git command
	r, err := git("tag", "--list", "--sort=-version:refname", "--merged", branch, fmt.Sprintf("%s*", prefix))
	if err != nil {
		return "", err
	}
	// Read first line of output
	return readLine(r)
}

// Finds commits in master since the specified revision in chronological order
func commitsSince(rev string) ([]commit, error) {
	// Run git command
	r, err := git("log", "--format=%H %s", "--reverse", fmt.Sprintf("^%s", rev), branch)
	if err != nil {
		return nil, err
	}

	// Read and parse output
	commits := make([]commit, 0)
	s := bufio.NewScanner(r)
	for s.Scan() {
		// Split lines into commit hash and title
		line := s.Text()
		fields := strings.SplitN(line, " ", 2)
		if len(fields) < 2 {
			continue
		}

		c := commit{Hash: fields[0], Title: strings.TrimSpace(fields[1])}

		// Get an additional line from the message of each commit
		r, err := git("show", "--format=%b", c.Hash)
		if err != nil {
			return nil, err
		}

		line, err = readLine(r)
		if err != nil {
			return nil, err
		}
		c.Body = strings.TrimSpace(line)

		commits = append(commits, c)
	}

	return commits, nil
}

// Finds PRs in the list of commits
func parsePullRequests(commits []commit) []pullRequest {
	prs := make([]pullRequest, 0)
	for _, c := range commits {
		// Find pull requests by their commit title
		match := prRegexp.FindSubmatch([]byte(c.Title))
		if match == nil {
			continue
		}

		// Parse attributes in the title
		n, _ := strconv.Atoi(string(match[1]))
		pr := pullRequest{Number: n, Author: string(match[2]), Commit: c}

		prs = append(prs, pr)
	}

	return prs
}

// Collects unique authors in the order of PRs
func collectAuthors(prs []pullRequest) []string {
	seen := make(map[string]bool)
	authors := make([]string, 0)
	for _, pr := range prs {
		if len(pr.Author) == 0 {
			continue
		}
		if _, ok := seen[pr.Author]; !ok {
			authors = append(authors, pr.Author)
		}
	}
	return authors
}

// usage: go run scripts/changelog.go [new version]
func main() {
	args := os.Args[1:]

	version := "TBD"
	if len(args) > 0 {
		version = args[0]
	}

	tag, _ := latestTag()
	commits, _ := commitsSince(tag)
	prs := parsePullRequests(commits)
	authors := collectAuthors(prs)

	clTemplate.Execute(os.Stdout, struct {
		Version, Date string
		PullRequests  []pullRequest
		Authors       []string
	}{version, time.Now().Format("2006-01-02"), prs, authors})
}
