package github

import (
	"context"

	"github.com/apex/log"
	"github.com/google/go-github/v42/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func getGHOauthToken() string {
	ghOauthToken, ok := viper.Get("gh_token").(string)

	if !ok || ghOauthToken == "" {
		log.Fatal("roc config file didn't contain valid GitHub oauth token")
	}

	ghUser := viper.Get("gh_username")

	log.Debugf("using GitHub token '%s' for user '%s'", ghOauthToken, ghUser)
	return ghOauthToken
}

func GetGitHubClient(ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: getGHOauthToken()},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}
