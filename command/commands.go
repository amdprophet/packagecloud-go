package command

import (
	"github.com/amdprophet/packagecloud-go/command/distro"
	"github.com/amdprophet/packagecloud-go/command/promote"
	"github.com/amdprophet/packagecloud-go/command/push"
	"github.com/amdprophet/packagecloud-go/command/search"
	"github.com/amdprophet/packagecloud-go/packagecloud"
	"github.com/spf13/cobra"
)

func AddCommands(rootCmd *cobra.Command, getClientFn packagecloud.GetClientFn) {
	rootCmd.AddCommand(
		distro.HelpCommand(getClientFn),
		push.PushCommand(getClientFn),
		promote.HelpCommand(getClientFn),
		search.SearchCommand(getClientFn),
	)
}
