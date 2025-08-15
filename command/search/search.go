package search

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	commanderrors "github.com/amdprophet/packagecloud-go/command/errors"
	"github.com/amdprophet/packagecloud-go/packagecloud"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

const (
	flagFormat      = "format"
	shortFlagFormat = "f"

	defaultFormat = "table"

	flagQuery      = "query"
	shortFlagQuery = "q"

	flagFilter      = "filter"
	shortFlagFilter = "i"

	flagDist      = "dist"
	shortFlagDist = "d"

	flagArch      = "arch"
	shortFlagArch = "a"

	flagPerPage      = "per-page"
	shortFlagPerPage = "p"

	flagWaitForIndexing      = "wait-for-indexing"
	shortFlagWaitForIndexing = "w"

	flagWaitSeconds      = "wait-seconds"
	shortFlagWaitSeconds = "s"

	flagWaitMaxRetries      = "wait-max-retries"
	shortFlagWaitMaxRetries = "r"
)

func SearchCommand(getClientFn packagecloud.GetClientFn) *cobra.Command {
	var repo []string

	cmd := &cobra.Command{
		Use:   "search user/repo",
		Short: "Search for packages matching given search parameters",
		Args: func(cmd *cobra.Command, args []string) error {
			newErrWithUsage := commanderrors.NewErrorWithUsageFactory(cmd.Help)

			if len(args) != 1 {
				return newErrWithUsage("requires exactly 1 argument")
			}

			repo = strings.Split(args[0], "/")
			if len(repo) != 2 {
				return newErrWithUsage("invalid repo, use format user/repo")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			newErrWithUsage := commanderrors.NewErrorWithUsageFactory(cmd.Help)

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

			perPage, err := cmd.Flags().GetString(flagPerPage)
			if err != nil {
				return err
			}

			waitForIndexing, err := cmd.Flags().GetBool(flagWaitForIndexing)
			if err != nil {
				return err
			}

			waitSeconds, err := cmd.Flags().GetInt(flagWaitSeconds)
			if err != nil {
				return err
			}

			waitMaxRetries, err := cmd.Flags().GetInt(flagWaitMaxRetries)
			if err != nil {
				return err
			}

			options := packagecloud.SearchOptions{
				RepoUser: repo[0],
				RepoName: repo[1],
				Query:    query,
				Filter:   filter,
				Dist:     dist,
				Arch:     arch,
				PerPage:  perPage,
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

			waitRetries := 0
			for {
				packages, err := client.Search(options)
				if err != nil {
					return fmt.Errorf("failed to retrieve search results: %s", err)
				}

				indexed := packages.Indexed()

				if waitForIndexing && !indexed {
					if waitRetries >= waitMaxRetries {
						// packages have not been indexed after max retries has been reached
						return fmt.Errorf("Packages have not finished indexing after %d seconds", waitSeconds*waitMaxRetries)
					}

					if format != "json" {
						fmt.Println("\nOne or more packages have not yet been indexed")
						fmt.Printf("Waiting %d seconds before trying again\n", waitSeconds)
						fmt.Println("")
					}

					for i := 0; i < waitSeconds; i++ {
						if format != "json" {
							fmt.Printf(".")
						}
						time.Sleep(1 * time.Second)
					}

					waitRetries++
					continue
				}

				if format == "json" {
					bytes, err := json.Marshal(packages)
					if err != nil {
						return fmt.Errorf("failed to marshal packages: %w", err)
					}
					fmt.Println(string(bytes))
					return nil
				}

				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Name", "Distro", "Version", "Release", "Epoch", "Indexed", "Filename", "Type"})
				table.SetAutoMergeCells(false)

				for _, pkg := range packages {
					row := []string{
						pkg.Name,
						pkg.DistroVersion,
						pkg.Version,
						pkg.Release,
						strconv.Itoa(pkg.Epoch),
						strconv.FormatBool(pkg.Indexed),
						pkg.Filename,
						pkg.Type,
					}
					table.Append(row)
				}
				table.Render()

				return nil
			}
		},
	}

	cmd.Flags().StringP(flagFormat, shortFlagFormat, defaultFormat, "output format to use - table or json")
	cmd.Flags().StringP(flagQuery, shortFlagQuery, "", "search string to search for package filename(s)")
	cmd.Flags().StringP(flagFilter, shortFlagFilter, "", "name of package type to search for packages (ignored when --dist is set)")
	cmd.Flags().StringP(flagDist, shortFlagDist, "", "name of the distribution to filter packages by (overrides --filter)")
	cmd.Flags().StringP(flagArch, shortFlagArch, "", "architecture to filter packages by (alpine/rpm/debian only)")
	cmd.Flags().StringP(flagPerPage, shortFlagPerPage, "256", "number of packages to return from the result set with each request")
	cmd.Flags().BoolP(flagWaitForIndexing, shortFlagWaitForIndexing, false, "wait for packages matching the search string to be indexed")
	cmd.Flags().IntP(flagWaitSeconds, shortFlagWaitSeconds, 10, "seconds to wait for retrying to check if packages have been indexed")
	cmd.Flags().IntP(flagWaitMaxRetries, shortFlagWaitMaxRetries, 12, "maximum amount of retry attempts to check if packages have been indexed")

	return cmd
}
