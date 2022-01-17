package github

import (
	"context"

	"github.com/apex/log"
	"github.com/google/go-github/v42/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func getGHOauthToken() string {
	// Using ~/.config/gh/hosts.yml
	ghcomVars := viper.Get("github.com")

	ghcomVarsMap := ghcomVars.(map[string]interface{})

	ghUser := ""
	ghOauthToken := ""
	var ok bool

	for key, value := range ghcomVarsMap {
		switch key {
		case "user":
			ghUser, ok = value.(string)
			if !ok {
				log.Fatal("expected gh cli config 'user' key to be a string")
			}
			continue
		case "oauth_token":
			ghOauthToken, ok = value.(string)
			if !ok {
				log.Fatal("expected gh cli config 'oauth_token' key to be a string")
			}
			continue
		}
	}

	if ghUser == "" || ghOauthToken == "" {
		log.Fatal("gh cli config file didn't contain GitHub user or oauth token")
	}

	log.Debugf("using GitHub username and token %s %s", ghUser, ghOauthToken)
	return ghOauthToken
}

func GetGHClient(ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: getGHOauthToken()},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}
