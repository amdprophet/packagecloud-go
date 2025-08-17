package versions

import (
	"fmt"

	commanderrors "github.com/amdprophet/packagecloud-go/command/errors"
	"github.com/amdprophet/packagecloud-go/packagecloud"
	"github.com/spf13/cobra"
)

func LatestVersionCommand(getClientFn packagecloud.GetClientFn) *cobra.Command {
	var repo packagecloud.Repo
	var name string

	cmd := &cobra.Command{
		Use:   "latest <user/repo> <package name>",
		Short: "Show the latest version of a package with a given name in a given repository",
		Args: func(cmd *cobra.Command, args []string) error {
			newErrWithUsage := commanderrors.NewErrorWithUsageFactory(cmd.Help)

			if len(args) != 2 {
				return newErrWithUsage("requires exactly 2 arguments")
			}

			if arg, err := packagecloud.NewRepoFromString(args[0]); err != nil {
				return newErrWithUsage(fmt.Sprintf("invalid repo: %s", err))
			} else {
				repo = arg
			}

			name = args[1]

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			newErrWithUsage := commanderrors.NewErrorWithUsageFactory(cmd.Help)

			filter, err := cmd.Flags().GetString(flagFilter)
			if err != nil {
				return err
			}

			dist, err := cmd.Flags().GetString(flagDist)
			if err != nil {
				return err
			}

			arch, err := cmd.Flags().GetString(flagArch)
			if err != nil {
				return err
			}

			perPage, err := cmd.Flags().GetString(flagPerPage)
			if err != nil {
				return err
			}

			options := packagecloud.ListVersionsOptions{
				Repo:        repo,
				PackageName: name,
				Filter:      filter,
				Dist:        dist,
				Arch:        arch,
				PerPage:     perPage,
			}

			if err := options.Validate(); err != nil {
				return newErrWithUsage(err.Error())
			}

			client, err := getClientFn()
			if err != nil {
				return err
			}

			latest, err := client.LatestVersion(options)
			if err != nil {
				return fmt.Errorf("failed to list versions: %w", err)
			}

			fmt.Println(latest)

			return nil
		},
	}

	cmd.Flags().StringP(flagFilter, shortFlagFilter, "", "name of package type to search for packages (ignored when --dist is set)")
	cmd.Flags().StringP(flagDist, shortFlagDist, "", "name of the distribution to filter packages by (overrides --filter)")
	cmd.Flags().StringP(flagArch, shortFlagArch, "", "architecture to filter packages by (alpine/rpm/debian only)")
	cmd.Flags().StringP(flagPerPage, shortFlagPerPage, "256", "number of packages to return from the result set with each request")

	return cmd
}
