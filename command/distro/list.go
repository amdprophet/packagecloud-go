package distro

import (
	"github.com/amdprophet/packagecloud-go/packagecloud"
	"github.com/spf13/cobra"
)

func ListCommand(getClientFn packagecloud.GetClientFn) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available distros and versions for package_type",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClientFn()
			if err != nil {
				return err
			}

			client.ListDistributions()

			return nil
		},
	}

	return cmd
}
