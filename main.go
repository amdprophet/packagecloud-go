package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/amdprophet/packagecloud-go/command"
	"github.com/amdprophet/packagecloud-go/packagecloud"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	url     string
	token   string
	verbose bool

	rootCmd = &cobra.Command{
		Use:   "packagecloud",
		Short: "A Go alternative to the official packagecloud command-line client",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			return nil
		},
	}
)

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.packagecloud)")
	rootCmd.PersistentFlags().StringVar(&url, "url", "", "website url to use (default is https://packagecloud.io)")
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "token to use for authentication")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "enable verbose mode")

	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

func initConfig() {
	viper.SetConfigType("json")

	if cfgFile != "" {
		fmt.Println("set")
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}

		viper.SetConfigType("yaml")
		viper.SetConfigFile(filepath.Join(home, ".packagecloud"))
	}

	viper.SetEnvPrefix("packagecloud")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		er(fmt.Errorf("error reading config: %s", err))
	}
}

func main() {
	rootCmd, err := newRootCmd()
	if err != nil {
		er(err)
	}
	cobra.OnInitialize(initCobra(rootCmd))

	if err := rootCmd.Execute(); err != nil {
		er(err)
	}
	getClientFn := func() (*packagecloud.Client, error) {
		var config packagecloud.Config
		if err := viper.Unmarshal(&config); err != nil {
			return nil, fmt.Errorf("unable to unmarshal config: %s", err)
		}
		if err := config.Validate(); err != nil {
			return nil, fmt.Errorf("config validation failed: %s", err)
		}

		return packagecloud.NewClient(config), nil
	}
	command.AddCommands(rootCmd, getClientFn)

	if err := rootCmd.Execute(); err != nil {
		er(err)
	}
}
