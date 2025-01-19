package ui

import (
	"fmt"

	"github.com/raymondji/git-stack-cli/config"
	"github.com/raymondji/git-stack-cli/githost"
)

func PrintBranchesInStack(
	branches []string,
	totalOrder bool,
	currBranch string,
	theme config.Theme,
	prsBySourceBranch map[string]githost.PullRequest,
	showPRs bool,
) {
	for i, branch := range branches {
		// if len(c.LocalBranches) == 0 && c.Hash == currCommit {
		// 	fmt.Println("* " + theme.PrimaryColor.Render(fmt.Sprintf("(HEAD detached at %s)", c.Hash)))
		// 	continue
		// } else if len(c.LocalBranches) == 0 {
		// 	continue
		// }
		var hereMarker, branchesSegment, suffix string
		if totalOrder && i == 0 {
			suffix = fmt.Sprintf(" (%s)", theme.TertiaryColor.Render("top"))
		}
		if branch == currBranch {
			hereMarker = "*"
			branchesSegment = theme.PrimaryColor.Render(branch)
		} else {
			hereMarker = " "
			if totalOrder && i == 0 && branch != currBranch {
				branchesSegment = theme.TertiaryColor.Render(branch)
			} else {
				branchesSegment = branch
			}
		}

		fmt.Printf("%s %s%s\n", hereMarker, branchesSegment, suffix)
		if showPRs {
			if pr, ok := prsBySourceBranch[branch]; ok {
				fmt.Printf("  └── %s\n", pr.WebURL)
			} else {
				fmt.Printf("  └── No pull request\n")
			}

			if i != len(branches)-1 {
				fmt.Println()
			}
		}
	}

	var missingPRs bool
	for _, b := range branches {
		_, ok := prsBySourceBranch[b]
		if !ok {
			missingPRs = true
		}
	}
	if showPRs && missingPRs {
		fmt.Println()
		fmt.Println(`some branches don't have pull requests yet (use "git stack push --create-prs" to open)`)
	}
}
