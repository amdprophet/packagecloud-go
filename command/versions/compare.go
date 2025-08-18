package versions

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	commanderrors "github.com/amdprophet/packagecloud-go/command/errors"
	"github.com/amdprophet/packagecloud-go/packagecloud"
	"github.com/spf13/cobra"
)

func CompareCommand(getClientFn packagecloud.GetClientFn) *cobra.Command {
	var a string
	var b string

	cmd := &cobra.Command{
		Use:   "compare <semver a> <semver b>",
		Short: "Compares semantic version 'a' to semantic version 'b'",
		Args: func(cmd *cobra.Command, args []string) error {
			newErrWithUsage := commanderrors.NewErrorWithUsageFactory(cmd.Help)

			if len(args) != 2 {
				return newErrWithUsage("requires exactly 2 arguments")
			}

			a = args[0]
			b = args[1]

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			newErrWithUsage := commanderrors.NewErrorWithUsageFactory(cmd.Help)

			semverA, err := semver.NewVersion(a)
			if err != nil {
				return newErrWithUsage(fmt.Sprintf("invalid semver a: %s", err))
			}

			semverB, err := semver.NewVersion(b)
			if err != nil {
				return newErrWithUsage(fmt.Sprintf("invalid semver b: %s", err))
			}

			if semverA.Equal(semverB) {
				fmt.Println("equal")
			} else if semverA.GreaterThan(semverB) {
				fmt.Println("greater")
			} else {
				fmt.Println("lesser")
			}

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
