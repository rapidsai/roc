package github

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/apex/log"
	"github.com/google/go-github/v42/github"
)

// the original Ops-Bot typescript code for formatting a pull request commit message from:
// https://github.com/rapidsai/ops-bot/blob/0b8e429e02ee46291c7fb579e529d3d32f307e2e/src/plugins/AutoMerger/auto_merger.ts#L279
// transliterated to Go
func FormatCommitMessageForPR(ghcli *github.Client, ctx context.Context, org, repo string, prNumber int) string {
	pr, _, err := ghcli.PullRequests.Get(ctx, org, repo, prNumber)
	if err != nil {
		log.Fatalf("error when looking up pr #%d on repo '%s/%s': '%s'", prNumber, org, repo, err.Error())
	}

	log.Debugf("formatting commit message for pr: %s", *pr.URL)

	prBody := *pr.Body
	prBody = strings.TrimSpace(RemoveHTMLComments(prBody))

	prBody = fmt.Sprintf("%s\n\n", prBody)
	prBody = fmt.Sprintf("%sAuthors:\n", prBody)

	listOpts := &github.ListOptions{}

	commits, _, err := ghcli.PullRequests.ListCommits(ctx, *pr.Base.Repo.Owner.Login, repo, prNumber, listOpts)
	if err != nil {
		log.Fatalf("Error when fetching commits for pr: %s", err.Error())
	}

	uniqueAuthors := make(map[string]bool)
	for _, commit := range commits {
		commitAuthor := commit.Author

		// I don't know if this nil check works
		if commitAuthor == nil {
			fmt.Printf("nil commit author! %+v\n", commit)
			continue
		}

		uniqueAuthors[*commitAuthor.Login] = true
	}

	formattedAuthors := FormatUsers(ghcli, ctx, uniqueAuthors)
	for _, formattedAuthor := range formattedAuthors {
		prBody = fmt.Sprintf("%s   - %s\n", prBody, formattedAuthor)
	}

	prBody = fmt.Sprintf("%s\nApprovers:\n", prBody)

	reviews, _, err := ghcli.PullRequests.ListReviews(ctx, *pr.Base.Repo.Owner.Login, repo, prNumber, listOpts)
	if err != nil {
		log.Fatalf("Error when fetching reviews for pr: %s", err.Error())
	}

	uniqueApprovers := make(map[string]bool)
	for _, review := range reviews {
		if *review.State != "APPROVED" {
			continue
		}

		approver := review.User

		// I don't know if this nil check works
		if approver == nil {
			fmt.Printf("nil reviewer! %+v\n", approver)
			continue
		}

		uniqueApprovers[*approver.Login] = true
	}

	formattedApprovers := FormatUsers(ghcli, ctx, uniqueApprovers)
	for _, formattedApprover := range formattedApprovers {
		prBody = fmt.Sprintf("%s   - %s\n", prBody, formattedApprover)
	}

	prBody = strings.TrimSuffix(prBody, "\n")

	return prBody
}

// TODO: unit test this
// cribbed from https://siongui.github.io/2016/03/10/go-minify-html/
func RemoveHTMLComments(content string) string {
	// https://www.google.com/search?q=regex+html+comments
	// http://stackoverflow.com/a/1084759
	htmlcmt := regexp.MustCompile(`<!--[^>]*-->`)
	return string(htmlcmt.ReplaceAll([]byte(content), []byte("")))
}

func FormatUsers(ghcli *github.Client, ctx context.Context, uniqueUsernames map[string]bool) []string {
	formattedUsers := []string{}

	for username := range uniqueUsernames {
		user, _, err := ghcli.Users.Get(ctx, username)
		if err != nil {
			log.Fatalf("Error fetching user '%s': %s", user, err.Error())
		}
		if user.Name != nil {
			formattedUsers = append(formattedUsers, fmt.Sprintf("%s (%s)", *user.Name, *user.HTMLURL))
		} else {
			formattedUsers = append(formattedUsers, *user.HTMLURL)
		}
	}

	return formattedUsers
}
