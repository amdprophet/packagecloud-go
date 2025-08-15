package promote

import (
	"fmt"

	commanderrors "github.com/amdprophet/packagecloud-go/command/errors"
	"github.com/amdprophet/packagecloud-go/packagecloud"
	"github.com/spf13/cobra"
)

const (
	flagQuery      = "query"
	shortFlagQuery = "q"

	flagFilter      = "filter"
	shortFlagFilter = "i"

	flagDist      = "dist"
	shortFlagDist = "d"

	flagArch      = "arch"
	shortFlagArch = "a"
)

func BySearchCommand(getClientFn packagecloud.GetClientFn) *cobra.Command {
	var srcRepo packagecloud.Repo
	var dstRepo packagecloud.Repo

	name := "by-search"
	usage := fmt.Sprintf("%s <%s> <%s> (%s)",
		name,
		"source repository",
		"destination repository",
		"-q | -i | -d | -a",
	)
	example := fmt.Sprintf("%s %s %s %s",
		name,
		"ecorp/staging",
		"ecorp/production",
		"-q '1.4.3-3258'",
	)

	cmd := &cobra.Command{
		Use:     usage,
		Short:   "Search for packages matching search options and promote all matches",
		Example: example,
		Args: func(cmd *cobra.Command, args []string) error {
			newErrWithUsage := commanderrors.NewErrorWithUsageFactory(cmd.Help)

			if len(args) != 2 {
				return newErrWithUsage("requires exactly 2 arguments")
			}

			if arg, err := packagecloud.NewRepoFromString(args[0]); err != nil {
				msg := fmt.Sprintf("invalid source repo: %s", err)
				return newErrWithUsage(msg)
			} else {
				srcRepo = arg
			}

			if arg, err := packagecloud.NewRepoFromString(args[1]); err != nil {
				msg := fmt.Sprintf("invalid destination repo: %s", err)
				return newErrWithUsage(msg)
			} else {
				dstRepo = arg
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			newErrWithUsage := commanderrors.NewErrorWithUsageFactory(cmd.Help)

			client, err := getClientFn()
			if err != nil {
				return err
			}

			query, err := cmd.Flags().GetString(flagQuery)
			if err != nil {
				return err
			}

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

			options := packagecloud.SearchOptions{
				RepoUser: srcRepo.User,
				RepoName: srcRepo.Name,
				Query:    query,
				Filter:   filter,
				Dist:     dist,
				Arch:     arch,
			}

			if err := options.Validate(); err != nil {
				return newErrWithUsage(err.Error())
			}

			if err := client.PromoteBySearch(dstRepo, options); err != nil {
				return fmt.Errorf("failed to promote packages: %s", err)
			}

			return nil
		},
	}

	cmd.Flags().StringP(flagQuery, shortFlagQuery, "", "search string to search for package filename(s)")
	cmd.Flags().StringP(flagFilter, shortFlagFilter, "", "name of package type to search for packages (ignored when --dist is set)")
	cmd.Flags().StringP(flagDist, shortFlagDist, "", "name of the distribution to filter packages by (overrides --filter)")
	cmd.Flags().StringP(flagArch, shortFlagArch, "", "architecture to filter packages by (alpine/rpm/debian only)")

	return cmd
}
