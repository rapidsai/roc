package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/text"
	"github.com/cli/oauth"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rapidsai/roc/internal/build"
	"github.com/rapidsai/roc/pkg/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	PROGNAME        = "roc"
	RAPIDSAI_ORG    = "rapidsai"
	ROC_CONFIG_NAME = ".roc"
)

var (
	rootCmd = &cobra.Command{
		Use:              PROGNAME,
		Short:            "A `gh` cli helper tool for RAPIDS Ops admins",
		Version:          fmt.Sprintf("%s-%s", build.Version, build.Date),
		PersistentPreRun: setupLogging,
		Long: `This tool adds some convenient commands for RAPIDS Ops admins to
do some routine tasks on GitHub.

Read more at https://github.com/rapidsai/roc`,
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	flags := rootCmd.PersistentFlags()

	// Set up cobra command line flags
	flags.BoolP("verbose", "v", false, "enable verbose logging")

	// Bind flags to viper
	if err := viper.BindPFlag("verbose", flags.Lookup("verbose")); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up flags: '%s'", err.Error())
		os.Exit(1)
	}
}

func initConfig() {
	home, err := homedir.Dir()
	cobra.CheckErr(err)

	viper.AddConfigPath(home)
	viper.SetConfigName(ROC_CONFIG_NAME)

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		file, err := os.Create(filepath.Join(home, fmt.Sprintf("%s.yml", ROC_CONFIG_NAME)))
		cobra.CheckErr(err)

		err = file.Chmod(0644)
		cobra.CheckErr(err)

		fmt.Println("Please authenticate yourself with GitHub")
		flow := &oauth.Flow{
			Hostname: "github.com",
			ClientID: "a1414cacc50a4d34227c",
			Scopes:   []string{"repo", "user"},
		}

		githubToken, err := flow.DetectFlow()
		if err != nil {
			panic(err)
		}
		fmt.Println("Authentication success!")
		viper.Set("gh_token", githubToken.Token)
		err = viper.WriteConfig()
		cobra.CheckErr(err)

		ctx := context.TODO()
		client := github.GetGitHubClient(ctx)

		user, _, err := client.Users.Get(ctx, "")
		if err != nil {
			panic(err)
		}

		viper.Set("gh_username", user.GetLogin())
		err = viper.WriteConfig()
		cobra.CheckErr(err)
	}

	viper.AutomaticEnv()
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("Couldn't run %s, please double check the usage", PROGNAME)
	}

	os.Exit(0)
}

func setupLogging(cmd *cobra.Command, args []string) {
	// Set up logging library - we have to do this in main, not init, cause of flags
	log.SetHandler(text.Default)

	if viper.GetBool("verbose") {
		log.SetLevel(log.DebugLevel)
	}

	log.SetHandler(cli.Default)
}
