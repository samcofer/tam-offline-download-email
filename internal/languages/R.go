package languages

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/samcofer/tam-offline-download-email/internal/config"
)

var availableRVersions = []string{
	"4.2.2", "4.2.1", "4.2.0", "4.1.3", "4.1.2", "4.1.1", "4.1.0", "4.0.5", "4.0.4", "4.0.3", "4.0.2", "4.0.1", "4.0.0", "3.6.3", "3.6.2", "3.6.1", "3.6.0", "3.5.3", "3.5.2", "3.5.1", "3.5.0", "3.4.4", "3.4.3", "3.4.2", "3.4.1", "3.4.0", "3.3.3", "3.3.2", "3.3.1", "3.3.0",
}

// GetRRootDirs returns the root directories for R

// GetRPaths returns the paths workbench will look for R
// underneath the root directories with the format
// /root/{rversion}/bin/R

// Detects if the path is an R directory

// PromptAndInstallR Prompts user if they want to install R and does the installation
func PromptAndInstallR(osType config.OperatingSystem) ([]string, error) {
	var urlR []string
	var urlDown string
	validRVersions, err := RetrieveValidRVersions()
	if err != nil {
		return []string{}, fmt.Errorf("issue retrieving R versions: %w", err)
	}
	installRVersions, err := RSelectVersionsPrompt(validRVersions)
	if err != nil {
		return []string{}, fmt.Errorf("issue selecting R versions: %w", err)
	}
	for _, rVersion := range installRVersions {
		urlDown, err = DownloadAndInstallR(rVersion, osType)
		urlR = append(urlR, urlDown)
		if err != nil {
			return []string{}, fmt.Errorf("issue installing R version: %w", err)
		}
	}
	return urlR, nil
}

// ScanAndHandleRVersions scans for R versions, handles result/errors and creates RConfig
func ScanAndHandleRVersions(osType config.OperatingSystem) ([]string, error) {

	urlR, err := PromptAndInstallR(osType)

	return urlR, err
}

// Append to a string slice only if the string is not yet in the slice

// ScanForRVersions scans for R versions in locations workbench will also look

// Prompt users if they would like to install R versions

func RetrieveValidRVersions() ([]string, error) {
	// TODO make this dynamic based on https://cran.r-project.org/src/base/R-4/ and https://cran.r-project.org/src/base/R-3/
	return availableRVersions, nil
}

// RSelectVersionsPrompt Prompt asking users which R version(s) they would like to install
func RSelectVersionsPrompt(availableRVersions []string) ([]string, error) {
	var qs = []*survey.Question{
		{
			Name: "rversions",
			Prompt: &survey.MultiSelect{
				Message: "Which version(s) of R would you like to install?",
				Options: availableRVersions,
				Default: availableRVersions[0],
			},
		},
	}
	rVersionsAnswers := struct {
		RVersions []string `survey:"rversions"`
	}{}
	err := survey.Ask(qs, &rVersionsAnswers)
	if err != nil {
		return []string{}, errors.New("there was an issue with the R versions selection prompt")
	}
	if len(rVersionsAnswers.RVersions) == 0 {
		return []string{}, errors.New("at least one R version must be selected")
	}
	return rVersionsAnswers.RVersions, nil
}

// DownloadAndInstallR Downloads the R installer, and installs R
func DownloadAndInstallR(rVersion string, osType config.OperatingSystem) (string, error) {
	// Create InstallerInfo with the proper information
	installerInfo, err := PopulateInstallerInfo("r", rVersion, osType)
	if err != nil {
		return "error: ", fmt.Errorf("PopulateInstallerInfo: %w", err)
	}
	// Download installer
	filepath := "\n  - R " + rVersion + ": " + installerInfo.URL
	if err != nil {
		return "error: ", fmt.Errorf("DownloadR: %w", err)
	}
	return filepath, nil
}
