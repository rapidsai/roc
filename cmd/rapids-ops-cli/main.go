package main

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/text"
	"github.com/rapidsai/rapids-ops-cli/internal/build"
	"github.com/rapidsai/rapids-ops-cli/internal/ghcli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	PROGNAME     = "roc"
	RAPIDSAI_ORG = "rapidsai"
)

var (
	rootCmd = &cobra.Command{
		Use:              PROGNAME,
		Short:            "A `gh` cli helper tool for RAPIDS Ops admins",
		Version:          fmt.Sprintf("%s-%s", build.Version, build.Date),
		PersistentPreRun: setupLogging,
		Long: `This tool adds some convenient commands for RAPIDS Ops admins to
do some routine tasks on GitHub.

Read more at https://github.com/rapidsai/rapids-ops-cli`,
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
	// Get the ghcli config file
	configPath := ghcli.GetGHCliConfigPath()

	viper.AddConfigPath(configPath)

	// gh cli stores the oauth key in hosts.yml
	viper.SetConfigName("hosts")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Couldn't find `gh` config file: please set up gh cli first")
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
