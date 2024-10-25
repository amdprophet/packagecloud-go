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
	"github.com/amdprophet/packagecloud-go/types"
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

			if len(args) < 1 {
				return newErrWithUsage("requires at least 2 arguments")
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

			if query == "" && filter == "" && dist == "" {
				return newErrWithUsage("one or more of the query, filter and/or dist flags must be specified")
			}

			client, err := getClientFn()
			if err != nil {
				return err
			}

			format, err := cmd.Flags().GetString(flagFormat)
			if err != nil {
				return fmt.Errorf("failed to parse format: %s", err)
			}

			options := packagecloud.SearchOptions{
				RepoUser: repo[0],
				RepoName: repo[1],
				Query:    query,
				Filter:   filter,
				Dist:     dist,
			}

			waitRetries := 0
			for {
				bytes, err := client.Search(options)
				if err != nil {
					return fmt.Errorf("failed to retrieve search results: %s", err)
				}

				if format == "json" {
					fmt.Println(string(bytes))
					return nil
				}

				var packages types.Packages
				if err := json.Unmarshal(bytes, &packages); err != nil {
					return fmt.Errorf("failed to unmarshal search results json: %s", err)
				}

				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"Name", "Distro", "Version", "Release", "Epoch", "Indexed"})
				table.SetAutoMergeCells(false)

				indexed := true
				for _, pkg := range packages {
					if !pkg.Indexed {
						indexed = false
					}
					row := []string{
						pkg.Name,
						pkg.DistroVersion,
						pkg.Version,
						pkg.Release,
						strconv.Itoa(pkg.Epoch),
						strconv.FormatBool(pkg.Indexed),
					}
					table.Append(row)
				}
				table.Render()

				if !waitForIndexing || indexed || waitRetries >= waitMaxRetries {
					break
				}

				fmt.Println("\nOne or more packages have not yet been indexed.")
				fmt.Printf("Waiting %d seconds before trying again.\n", waitSeconds)
				for i := 0; i < waitSeconds; i++ {
					fmt.Printf(".")
					time.Sleep(1 * time.Second)
				}
				fmt.Println("")

				waitRetries++
			}

			return nil
		},
	}

	cmd.Flags().StringP(flagFormat, shortFlagFormat, defaultFormat, "output format to use - table or json")
	cmd.Flags().StringP(flagQuery, shortFlagQuery, "", "search string to search for package filename(s)")
	cmd.Flags().StringP(flagFilter, shortFlagFilter, "", "name of package type to search for packages (ignored when --dist is set)")
	cmd.Flags().StringP(flagDist, shortFlagDist, "", "name of the distribution to filter packages by (overrides --filter)")
	cmd.Flags().BoolP(flagWaitForIndexing, shortFlagWaitForIndexing, false, "wait for packages matching the search string to be indexed")
	cmd.Flags().IntP(flagWaitSeconds, shortFlagWaitSeconds, 5, "seconds to wait for retrying to check if packages have been indexed")
	cmd.Flags().IntP(flagWaitMaxRetries, shortFlagWaitMaxRetries, 12, "maximum amount of retry attempts to check if packages have been indexed")

	return cmd
}
