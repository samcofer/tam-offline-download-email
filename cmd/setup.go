package cmd

import (
	"fmt"
	"github.com/samcofer/tam-offline-download-email/internal/prodrivers"
	"github.com/samcofer/tam-offline-download-email/internal/workbench"
	"github.com/spf13/viper"

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

func newSetup(OS string) error {

	//fmt.Println("Welcome to the Workbench Installer!\n")

	//// Check if running as root
	//err := os.CheckIfRunningAsRoot()
	//if err != nil {
	//	return err
	//}
	// Determine OS and install pre-requisites

	osType, err := os.DetectOS(OS)
	if err != nil {
		return fmt.Errorf("issue selecting languages: %w", err)
	}

	// R
	err = languages.ScanAndHandleRVersions(osType)
	if err != nil {
		return fmt.Errorf("issue finding R locations: %w", err)
	}

	//Python
	err = languages.ScanAndHandlePythonVersions(osType)
	if err != nil {
		return fmt.Errorf("issue finding Python locations: %w", err)
	}

	//workbench
	workbench.DownloadAndInstallWorkbench(osType)
	if err != nil {
		return fmt.Errorf("issue installing Workbench: %w", err)
	}

	//drivers
	prodrivers.DownloadAndInstallProDrivers(osType)

	return nil
}

func setSetupOpts() {

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
			setSetupOpts()
			if err := root.opts.Validate(); err != nil {
				return err
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			//TODO: Add your logic to gather config to pass code here
			log.WithField("opts", fmt.Sprintf("%+v", root.opts)).Trace("setup-opts")
			if err := newSetup(viper.GetString("os")); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.PersistentFlags().String("os", "rhel8", "destinationOS")
	viper.BindPFlag("os", cmd.PersistentFlags().Lookup("os"))
	root.cmd = cmd
	return root
}
