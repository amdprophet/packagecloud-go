package promote

import (
	"fmt"

	commanderrors "github.com/amdprophet/packagecloud-go/command/errors"
	"github.com/amdprophet/packagecloud-go/packagecloud"
	"github.com/spf13/cobra"
)

func ByFilenameCommand(getClientFn packagecloud.GetClientFn) *cobra.Command {
	var srcRepo packagecloud.Repo
	var dstRepo packagecloud.Repo
	var distro packagecloud.Distro
	var filename string

	name := "by-filename"
	usage := fmt.Sprintf("%s <%s> <%s> <%s> <%s>",
		name,
		"source repository",
		"destination repository",
		"distro",
		"filename",
	)
	example := fmt.Sprintf("%s %s %s %s %s",
		name,
		"ecorp/staging",
		"ecorp/production",
		"ubuntu/jammy",
		"package.deb",
	)

	cmd := &cobra.Command{
		Use:     usage,
		Short:   "Promote a single package from one repository to another by filename",
		Example: example,
		Args: func(cmd *cobra.Command, args []string) error {
			newErrWithUsage := commanderrors.NewErrorWithUsageFactory(cmd.Help)

			if len(args) != 4 {
				return newErrWithUsage("requires exactly 4 arguments")
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

			if arg, err := packagecloud.NewDistroFromString(args[2]); err != nil {
				msg := fmt.Sprintf("invalid distro: %s", err)
				return newErrWithUsage(msg)
			} else {
				distro = arg
			}

			filename = args[3]

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClientFn()
			if err != nil {
				return err
			}

			if err := client.PromoteByFilename(srcRepo, dstRepo, distro, filename); err != nil {
				return fmt.Errorf("failed to promote package: %s", err)
			}

			return nil
		},
	}

	return cmd
}
