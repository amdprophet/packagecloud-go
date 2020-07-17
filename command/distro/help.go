package distro

import (
	"github.com/amdprophet/packagecloud-go/packagecloud"
	"github.com/spf13/cobra"
)

func HelpCommand(getClientFn packagecloud.GetClientFn) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "distro",
		Short: "Manage distros",
	}

	cmd.AddCommand(ListCommand(getClientFn))

	return cmd
}
