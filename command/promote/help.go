package promote

import (
	"github.com/amdprophet/packagecloud-go/packagecloud"
	"github.com/spf13/cobra"
)

func HelpCommand(getClientFn packagecloud.GetClientFn) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "promote",
		Short: "Promote a package or packages",
	}

	cmd.AddCommand(ByFilenameCommand(getClientFn))
	cmd.AddCommand(BySearchCommand(getClientFn))

	return cmd
}
