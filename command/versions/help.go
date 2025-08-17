package versions

import (
	"github.com/amdprophet/packagecloud-go/packagecloud"
	"github.com/spf13/cobra"
)

func HelpCommand(getClientFn packagecloud.GetClientFn) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "versions",
		Short: "Perform various actions related to versions",
	}

	cmd.AddCommand(ListCommand(getClientFn))
	cmd.AddCommand(LatestVersionCommand(getClientFn))
	cmd.AddCommand(PreviousVersionCommand(getClientFn))

	return cmd
}
