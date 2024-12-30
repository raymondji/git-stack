package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/raymondji/git-stack/commitstack"
	"github.com/raymondji/git-stack/githost"
	"github.com/raymondji/git-stack/githost/gitlab"
	"github.com/raymondji/git-stack/libgit"
	"github.com/spf13/cobra"
)

const (
	cacheDuration = 14 * 24 * time.Hour // 14 days
)

var (
	primaryColor   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500")) // Orange
	secondaryColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")) // Green
)

func main() {
	ok, err := isInstalled("glab")
	if err != nil {
		fmt.Println(err.Error())
		return
	} else if !ok {
		fmt.Println("glab CLI must be installed")
		return
	}
	git := libgit.Git{}
	var host githost.Host = gitlab.Gitlab{}
	defaultBranch, err := getDefaultBranchCached(git, host)
	if err != nil {
		fmt.Println("failed to get default branch, are you authenticated to glab?", err)
		return
	}

	var rootCmd = &cobra.Command{
		Use:   "stack",
		Short: "A CLI tool for managing stacked Git branches.",
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the current version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("0.0.1")
		},
	}

	addCmd := newAddCmd(git, defaultBranch)
	switchCmd := newSwitchCmd(git, defaultBranch)
	pushCmd := newPushCmd(git, host, defaultBranch)
	pullCmd := newPullCmd(git, defaultBranch)
	editCmd := newEditCmd(git, defaultBranch)
	fixupCmd := newFixupCmd(git, defaultBranch)
	showCmd := newShowCmd(git, host, defaultBranch)
	listCmd := newListCmd(git, defaultBranch)
	logCmd := newLogCmd(git, defaultBranch)

	rootCmd.SilenceUsage = true
	rootCmd.AddCommand(versionCmd, addCmd, logCmd, editCmd, fixupCmd, listCmd, switchCmd, showCmd, pushCmd, pullCmd)
	rootCmd.Execute()
}

func formatPullRequestDescription(
	currPR githost.PullRequest, prs []githost.PullRequest,
) string {
	var newStackDesc string
	if len(prs) == 1 {
		// (raymond):
	} else {
		var newStackDescParts []string
		currIndex := slices.IndexFunc(prs, func(pr githost.PullRequest) bool {
			return pr.SourceBranch == currPR.SourceBranch
		})

		for i, pr := range prs {
			var prefix string
			if i == currIndex {
				prefix = "Current: "
			} else if i == currIndex-1 {
				prefix = "Next: "
			} else if i == currIndex+1 {
				prefix = "Prev: "
			}
			newStackDescParts = append(newStackDescParts, fmt.Sprintf("- %s%s", prefix, pr.MarkdownWebURL))
		}

		newStackDesc = "Merge request stack:\n" + strings.Join(newStackDescParts, "\n")
	}

	beginMarker := "<!-- DO NOT EDIT: generated by git stack push (start)-->"
	endMarker := "<!-- DO NOT EDIT: generated by git stack push (end) -->"
	newSection := fmt.Sprintf("%s\n%s\n%s", beginMarker, newStackDesc, endMarker)
	sectionPattern := regexp.MustCompile(`(?s)` + regexp.QuoteMeta(beginMarker) + `.*?` + regexp.QuoteMeta(endMarker))

	if sectionPattern.MatchString(currPR.Description) {
		return sectionPattern.ReplaceAllString(currPR.Description, newSection)
	} else {
		return fmt.Sprintf("%s\n\n%s", strings.TrimSpace(currPR.Description), newSection)
	}
}

func printProblem(stack commitstack.Stack) {
	if stack.Error != nil {
		fmt.Println()
		fmt.Println("Problems detected:")
		fmt.Printf("  %s\n", stack.Error.Error())
	}
}

func printProblems(stacks commitstack.Stacks) {
	if len(stacks.Errors) > 0 {
		fmt.Println()
		fmt.Println("Problems detected:")
		for _, err := range stacks.Errors {
			fmt.Printf("  %s\n", err.Error())
		}
	}
}

func isInstalled(file string) (bool, error) {
	_, err := exec.LookPath(file)
	var execErr *exec.Error
	if errors.As(err, &execErr) {
		// Generally returned when file is not a executable
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("error checking if %s is installed, err: %v", file, err)
	}
	return true, nil
}

func getDefaultBranchCached(git libgit.Git, host githost.Host) (string, error) {
	rootDir, err := git.GetRootDir()
	if err != nil {
		return "", nil
	}

	cacheDir := path.Join("/tmp/git-stack", path.Base(rootDir))
	err = os.MkdirAll(cacheDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create cache directory: %v", err)
	}

	// Check if cache file exists
	cacheFilePath := path.Join(cacheDir, "defaultBranch.txt")
	cacheInfo, err := os.Stat(cacheFilePath)
	if err == nil {
		if time.Since(cacheInfo.ModTime()) < cacheDuration {
			data, err := os.ReadFile(cacheFilePath)
			if err == nil {
				return string(data), nil
			}
		}
	}

	// Fetch from GitHost
	repo, err := host.GetRepo()
	if err != nil {
		return "", err
	}

	// Save to cache
	err = os.WriteFile(cacheFilePath, []byte(repo.DefaultBranch), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write cache file, err: %v", err)
	}

	return repo.DefaultBranch, nil
}
