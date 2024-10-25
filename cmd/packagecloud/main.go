package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/amdprophet/packagecloud-go/command"
	commanderrors "github.com/amdprophet/packagecloud-go/command/errors"
	"github.com/amdprophet/packagecloud-go/packagecloud"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	flagConfig  = "config"
	flagURL     = "url"
	flagToken   = "token"
	flagVerbose = "verbose"

	defaultURL     = "https://packagecloud.io"
	defaultToken   = ""
	defaultVerbose = false
)

var (
	cfgFile string
)

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func defaultConfigPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".packagecloud"), nil
}

func initCobra(cmd *cobra.Command) func() {
	return func() {
		viper.SetEnvPrefix("packagecloud")
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
		viper.AutomaticEnv()

		if err := postInitCommands(cmd.Commands()); err != nil {
			er(err)
		}

		viper.SetConfigType("json")
		viper.SetConfigFile(cfgFile)
		viper.ReadInConfig()
	}
}

func postInitCommands(commands []*cobra.Command) error {
	for _, cmd := range commands {
		if err := presetRequiredFlags(cmd); err != nil {
			return err
		}

		if cmd.HasSubCommands() {
			if err := postInitCommands(cmd.Commands()); err != nil {
				return err
			}
		}
	}
	return nil
}

func presetRequiredFlags(cmd *cobra.Command) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	var fErr error
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			if err := cmd.Flags().Set(f.Name, viper.GetString(f.Name)); err != nil {
				fErr = err
				return
			}
		}
	})
	return fErr
}

func newRootCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:           "packagecloud",
		Short:         "A Go alternative to the official packagecloud command-line client",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	defaultConfig, err := defaultConfigPath()
	if err != nil {
		return nil, err
	}

	cmd.PersistentFlags().StringVar(&cfgFile, flagConfig, defaultConfig, "config file")
	cmd.PersistentFlags().String(flagURL, defaultURL, "website url to use")
	cmd.PersistentFlags().String(flagToken, defaultToken, "token to use for authentication")
	cmd.PersistentFlags().Bool(flagVerbose, defaultVerbose, "enable verbose mode")

	cmd.MarkFlagRequired(flagToken)

	viper.BindPFlag(flagURL, cmd.PersistentFlags().Lookup(flagURL))
	viper.BindPFlag(flagToken, cmd.PersistentFlags().Lookup(flagToken))
	viper.BindPFlag(flagVerbose, cmd.PersistentFlags().Lookup(flagVerbose))

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

	command.AddCommands(cmd, getClientFn)

	return cmd, nil
}

func main() {
	rootCmd, err := newRootCmd()
	if err != nil {
		er(err)
	}
	cobra.OnInitialize(initCobra(rootCmd))

	if err := rootCmd.Execute(); err != nil {
		if argsErr, ok := err.(*commanderrors.ErrInvalidArgs); ok {
			fmt.Printf("Error: %s\n\n", argsErr)
			rootCmd.Help()
			os.Exit(1)
		}
		if argsErr, ok := err.(*commanderrors.ErrWithUsage); ok {
			fmt.Printf("Error: %s\n\n", argsErr)
			argsErr.Usage()
			os.Exit(1)
		}
		er(err)
	}
}
