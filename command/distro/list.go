package distro

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/amdprophet/packagecloud-go/packagecloud"
	"github.com/amdprophet/packagecloud-go/types"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

const (
	flagFormat      = "format"
	shortFlagFormat = "f"

	defaultFormat = "table"
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

			format, err := cmd.Flags().GetString(flagFormat)
			if err != nil {
				return fmt.Errorf("failed to parse format: %s", err)
			}

			bytes, err := client.GetDistributions()
			if err != nil {
				return fmt.Errorf("failed to retrieve distributions: %s", err)
			}

			if format == "json" {
				fmt.Println(string(bytes))
				return nil
			}

			var packageTypes types.PackageTypes
			if err := json.Unmarshal(bytes, &packageTypes); err != nil {
				return fmt.Errorf("failed to unmarshal distributions json: %s", err)
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Package Type", "Distro/Version"})
			table.SetAutoMergeCells(true)

			for packageType, distros := range packageTypes {
				for _, distro := range distros {
					for _, distroVersion := range distro.Versions {
						row := []string{
							packageType,
							fmt.Sprintf("%s/%s", distro.IndexName, distroVersion.IndexName),
						}
						table.Append(row)
					}
				}
			}
			table.Render()

			return nil
		},
	}

	cmd.Flags().StringP(flagFormat, shortFlagFormat, defaultFormat, "output format to use - table or json")

	return cmd
}
