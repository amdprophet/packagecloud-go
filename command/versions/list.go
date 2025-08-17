package versions

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/Masterminds/semver/v3"
	commanderrors "github.com/amdprophet/packagecloud-go/command/errors"
	"github.com/amdprophet/packagecloud-go/packagecloud"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

const (
	flagFormat      = "format"
	shortFlagFormat = "f"

	defaultFormat = "table"

	flagFilter      = "filter"
	shortFlagFilter = "i"

	flagDist      = "dist"
	shortFlagDist = "d"

	flagArch      = "arch"
	shortFlagArch = "a"

	flagPerPage      = "per-page"
	shortFlagPerPage = "p"
)

func ListCommand(getClientFn packagecloud.GetClientFn) *cobra.Command {
	var repo packagecloud.Repo
	var name string

	cmd := &cobra.Command{
		Use:   "list <user/repo> <package name>",
		Short: "List all versions of packages with a given name in a given repository",
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

			format, err := cmd.Flags().GetString(flagFormat)
			if err != nil {
				return fmt.Errorf("failed to parse format: %s", err)
			}

			versions, err := client.ListVersions(options)
			if err != nil {
				return fmt.Errorf("failed to list versions: %w", err)
			}

			if format == "json" {
				bytes, err := json.Marshal(versions)
				if err != nil {
					return fmt.Errorf("failed to marshal versions: %w", err)
				}
				fmt.Println(string(bytes))
				return nil
			}

			keys := make([]*semver.Version, 0, len(versions))
			for key := range versions {
				version, err := semver.NewVersion(key)
				if err != nil {
					return fmt.Errorf("failed to parse version %s: %w", key, err)
				}

				keys = append(keys, version)
			}
			sort.Sort(sort.Reverse(semver.Collection(keys)))

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Version", "Number of packages"})
			table.SetAutoMergeCells(false)

			for _, key := range keys {
				row := []string{
					options.PackageName,
					key.String(),
					strconv.Itoa(versions[key.String()]),
				}
				table.Append(row)
			}
			table.Render()

			return nil
		},
	}

	cmd.Flags().StringP(flagFormat, shortFlagFormat, defaultFormat, "output format to use - table or json")
	cmd.Flags().StringP(flagFilter, shortFlagFilter, "", "name of package type to search for packages (ignored when --dist is set)")
	cmd.Flags().StringP(flagDist, shortFlagDist, "", "name of the distribution to filter packages by (overrides --filter)")
	cmd.Flags().StringP(flagArch, shortFlagArch, "", "architecture to filter packages by (alpine/rpm/debian only)")
	cmd.Flags().StringP(flagPerPage, shortFlagPerPage, "256", "number of packages to return from the result set with each request")

	return cmd
}
