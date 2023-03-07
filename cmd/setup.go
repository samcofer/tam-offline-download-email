package cmd

import (
	"fmt"
	"github.com/samcofer/tam-offline-download-email/internal/prodrivers"
	"github.com/samcofer/tam-offline-download-email/internal/workbench"

	"github.com/samber/lo"
	"github.com/samcofer/tam-offline-download-email/internal/config"
	"github.com/samcofer/tam-offline-download-email/internal/languages"
	"github.com/samcofer/tam-offline-download-email/internal/os"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type setupCmd struct {
	cmd  *cobra.Command
	opts setupOpts
}

type setupOpts struct {
}

func newSetup(setupOpts setupOpts) error {

	var WBConfig config.WBConfig

	//fmt.Println("Welcome to the Workbench Installer!\n")

	//// Check if running as root
	//err := os.CheckIfRunningAsRoot()
	//if err != nil {
	//	return err
	//}

	// Determine OS and install pre-requisites
	osType, err := os.DetectOS()
	if err != nil {
		return err
	}

	// Languages
	selectedLanguages, err := languages.PromptAndRespond()
	if err != nil {
		return fmt.Errorf("issue selecting languages: %w", err)
	}

	// R
	WBConfig.RConfig.Paths, err = languages.ScanAndHandleRVersions(osType)
	if err != nil {
		return fmt.Errorf("issue finding R locations: %w", err)
	}
	// remove any path that starts with /usr and only offer symlinks for those that don't (i.e. /opt directories)
	rPathsFiltered := languages.RemoveSystemRPaths(WBConfig.RConfig.Paths)
	// check if R and Rscript has already been symlinked
	rSymlinked := languages.CheckIfRSymlinkExists()
	rScriptSymlinked := languages.CheckIfRscriptSymlinkExists()
	if (len(rPathsFiltered) > 0) && !rSymlinked && !rScriptSymlinked {
		err = languages.PromptAndSetRSymlinks(rPathsFiltered)
		if err != nil {
			return fmt.Errorf("issue setting R symlinks: %w", err)
		}
	}
	if lo.Contains(selectedLanguages, "python") {
		WBConfig.PythonConfig.Paths, err = languages.ScanAndHandlePythonVersions(osType)
		if err != nil {
			return fmt.Errorf("issue finding Python locations: %w", err)
		}
	}

	// If Workbench is not detected then prompt to install

	workbench.DownloadAndInstallWorkbench(osType)
	if err != nil {
		return fmt.Errorf("issue installing Workbench: %w", err)
	}

	// Pro Drivers
	//proDriversExistingStatus, err := prodrivers.CheckExistingProDrivers()
	if err != nil {
		return fmt.Errorf("issue in checking for prior pro driver installation: %w", err)
	}

	prodrivers.DownloadAndInstallProDrivers(osType)

	return nil
}

func setSetupOpts(setupOpts *setupOpts) {

}

func (opts *setupOpts) Validate() error {
	return nil
}

func newSetupCmd() *setupCmd {
	root := &setupCmd{opts: setupOpts{}}

	cmd := &cobra.Command{
		Use:   "setup",
		Short: "setup",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			setSetupOpts(&root.opts)
			if err := root.opts.Validate(); err != nil {
				return err
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			//TODO: Add your logic to gather config to pass code here
			log.WithField("opts", fmt.Sprintf("%+v", root.opts)).Trace("setup-opts")
			if err := newSetup(root.opts); err != nil {
				return err
			}
			return nil
		},
	}
	root.cmd = cmd
	return root
}
