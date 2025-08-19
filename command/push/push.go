package push

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	commanderrors "github.com/amdprophet/packagecloud-go/command/errors"
	"github.com/amdprophet/packagecloud-go/packagecloud"
	"github.com/amdprophet/packagecloud-go/types"
	"github.com/spf13/cobra"
)

const (
	flagSkipExists = "skip-exists"

	defaultSkipExists = false
)

func PushCommand(getClientFn packagecloud.GetClientFn) *cobra.Command {
	var repo []string

	cmd := &cobra.Command{
		Use:   "push user/repo/distro/version /path/to/packages",
		Short: "Push package(s) to repository",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return &commanderrors.ErrInvalidArgs{Msg: "requires at least 2 arguments"}
			}

			repo = strings.Split(args[0], "/")
			if len(repo) != 4 {
				return &commanderrors.ErrInvalidArgs{Msg: "invalid repo, use format user/repo/distro/version"}
			}

			if err := packagecloud.ValidateFileExtensions(args[1:]); err != nil {
				return &commanderrors.ErrInvalidArgs{Msg: err.Error()}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClientFn()
			if err != nil {
				return err
			}

			skipExists, err := cmd.Flags().GetBool(flagSkipExists)
			if err != nil {
				return fmt.Errorf("failed to parse %s: %s", flagSkipExists, err)
			}

			filePaths := args[1:]

			packageTypes, err := client.GetDistributions()
			if err != nil {
				return fmt.Errorf("failed to fetch distributions: %s", err)
			}

			packageType := filepath.Ext(filePaths[0])[1:]
			if _, ok := packageTypes[packageType]; !ok {
				return fmt.Errorf("failed to find package type in distributions: %s", packageType)
			}

			distroID, err := types.GetDistroID(packageTypes[packageType], packageType, repo[2], repo[3])
			if err != nil {
				return err
			}

			for _, filePath := range filePaths {
				options := packagecloud.PushPackageOptions{
					RepoUser: repo[0],
					RepoName: repo[1],
					DistroID: strconv.Itoa(distroID),
					FilePath: filePaths[0],
				}

				fmt.Println("uploading package:", filePath)
				_, err := client.PushPackage(options)
				if err != nil {
					if skipExists && errors.Is(err, packagecloud.ErrPackageAlreadyExists) {
						fmt.Println("package already exists, skipping...")
						continue
					}
					return fmt.Errorf("failed to upload package: %s", err)
				}
			}

			return nil
		},
	}

	cmd.Flags().Bool(flagSkipExists, defaultSkipExists, "skip over packages that already exist")

	return cmd
}
