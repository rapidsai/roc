package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/apex/log"
	"github.com/rapidsai/roc/pkg/github"
	"github.com/spf13/cobra"
)

var PRCommitMsgCmd = &cobra.Command{
	Use:   "prcommitmsg rapidsai-repo-name pr-number",
	Short: "Format a commit message for a PR in a rapidsai GitHub repo",
	Args:  cobra.ExactArgs(2),
	Long: `This command formats a commit message from a string of commits in
a pull request. It is meant to help RAPIDS Ops members force-merge PRs that
are failing their status checks.

It formats commit messages similar to how ops-bot does in its auto-merge.`,
	Run: runPRCommitMsgCmd,
}

func init() {
	rootCmd.AddCommand(PRCommitMsgCmd)
}

func runPRCommitMsgCmd(cmd *cobra.Command, args []string) {
	repo := args[0]

	prnum, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatalf("pr '%s' couldn't be converted to int", args[1])
	}

	log.Debugf("looking up pr %d in repo %s...", prnum, repo)

	ctx := context.Background()
	ghcli := github.GetGitHubClient(ctx)

	commitMessage := github.FormatCommitMessageForPR(ghcli, ctx, RAPIDSAI_ORG, repo, prnum)

	fmt.Println(commitMessage)
}
