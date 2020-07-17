package command

import (
	"github.com/amdprophet/packagecloud-go/command/distro"
	"github.com/amdprophet/packagecloud-go/packagecloud"
	"github.com/spf13/cobra"
)

// distro SUBCMD ...ARGS
// gpg_key SUBCMD ...ARGS
// help [COMMAND]
// master_token SUBCMD ...ARGS
// promote user/repo[/distro/version] [@scope/]package_name user/destination_repo
// push user/repo[/distro/version] /path/to/packages
// read_token SUBCMD ...ARGS
// repository SUBCMD ...ARGS
// version
// yank user/repo[/distro/version] [@scope/]package_name

func AddCommands(rootCmd *cobra.Command, getClientFn packagecloud.GetClientFn) {
	rootCmd.AddCommand(
		distro.HelpCommand(getClientFn),
	)
}
