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
	output "os"
	"text/template"
)

type setupCmd struct {
	cmd  *cobra.Command
	opts setupOpts
}

type setupOpts struct {
}

type CustomerData struct {
	Customer       string
	R              []string
	Python         []string
	Quarto         string
	Workbench      string
	Connect        string
	PackageManager string
	ProDriver      string
}

func newSetup(OS string, customer string) error {

	var URLs CustomerData

	// Determine OS and install pre-requisites
	osType, err := os.DetectOS(OS)
	if err != nil {
		return fmt.Errorf("issue selecting languages: %w", err)
	}

	// R
	URLs.R, err = languages.ScanAndHandleRVersions(osType)
	if err != nil {
		return fmt.Errorf("issue finding R locations: %w", err)
	}

	//Python
	URLs.Python, err = languages.ScanAndHandlePythonVersions(osType)
	if err != nil {
		return fmt.Errorf("issue finding Python locations: %w", err)
	}

	//workbench
	URLs.Workbench, URLs.Connect, URLs.PackageManager, err = workbench.DownloadAndInstallWorkbench(osType)
	if err != nil {
		return fmt.Errorf("issue installing Workbench: %w", err)
	}

	//drivers
	URLs.ProDriver, err = prodrivers.DownloadAndInstallProDrivers(osType)

	//Quarto
	URLs.Quarto, err = languages.DownloadAndInstallQuarto(osType)
	//fmt.Println(URLs.R)
	//fmt.Println(URLs.Python)
	//fmt.Println(URLs.Workbench)
	//fmt.Println(URLs.Connect)
	//fmt.Println(URLs.PackageManager)
	//fmt.Println(URLs.ProDriver)

	URLs.Customer = customer
	err = EmailTemplate(URLs)
	if err != nil {
		return fmt.Errorf("template issue: %w", err)
	}
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
			if err := newSetup(viper.GetString("os"), viper.GetString("customer")); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.PersistentFlags().String("os", "rhel8", "destinationOS")
	viper.BindPFlag("os", cmd.PersistentFlags().Lookup("os"))
	cmd.PersistentFlags().String("customer", "REPLACE_ME", "Customer Name")
	viper.BindPFlag("customer", cmd.PersistentFlags().Lookup("customer"))
	root.cmd = cmd
	return root
}

// Declare type pointer to a template
var temp *template.Template

func EmailTemplate(custData CustomerData) error {

	temp = template.Must(template.ParseFiles("email.md"))
	err := temp.Execute(output.Stdout, custData)
	if err != nil {
		return err
	}
	return nil
}
