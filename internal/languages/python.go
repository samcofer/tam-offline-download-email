package languages

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/samcofer/tam-offline-download-email/internal/config"
	"github.com/samcofer/tam-offline-download-email/internal/install"
)

var availablePythonVersions = []string{
	"3.11.2",
	"3.11.1",
	"3.11.0",
	"3.10.10",
	"3.10.9",
	"3.10.8",
	"3.10.7",
	"3.10.6",
	"3.10.5",
	"3.10.4",
	"3.10.3",
	"3.10.2",
	"3.10.1",
	"3.10.0",
	"3.9.16",
	"3.9.15",
	"3.9.14",
	"3.9.13",
	"3.9.12",
	"3.9.11",
	"3.9.10",
	"3.9.9",
	"3.9.8",
	"3.9.7",
	"3.9.6",
	"3.9.5",
	"3.9.4",
	"3.9.3",
	"3.9.2",
	"3.9.1",
	"3.9.0",
	"3.8.16",
	"3.8.15",
	"3.8.14",
	"3.8.13",
	"3.8.12",
	"3.8.11",
	"3.8.10",
	"3.8.9",
	"3.8.8",
	"3.8.7",
	"3.8.6",
	"3.8.5",
	"3.8.4",
	"3.8.3",
	"3.8.2",
	"3.8.1",
	"3.8.0",
	"3.7.16",
	"3.7.15",
	"3.7.14",
	"3.7.13",
	"3.7.12",
	"3.7.11",
	"3.7.10",
	"3.7.9",
	"3.7.8",
	"3.7.7",
	"3.7.6",
	"3.7.5",
	"3.7.4",
	"3.7.3",
	"3.7.2",
	"3.7.1",
	"3.7.0",
}

// GetPythonRootDirs returns the root directories for Python

// GetPythonPaths returns the paths workbench will look for Python
// underneath the root directories with the format
// /root/{pythonVersion}/bin/python

// PromptAndSetPythonPATH prompts user to set Python PATH

// PythonLocationPATHPrompt asks users which Python binary they want to add to PATH

// Prompts user if they want to install Python and does the installation
func PromptAndInstallPython(osType config.OperatingSystem) ([]string, error) {

	validPythonVersions, err := RetrieveValidPythonVersions()
	if err != nil {
		return []string{}, fmt.Errorf("issue retrieving Python versions: %w", err)
	}
	installPythonVersions, err := PythonSelectVersionsPrompt(validPythonVersions)
	if err != nil {
		return []string{}, fmt.Errorf("issue selecting Python versions: %w", err)
	}
	for _, pythonVersion := range installPythonVersions {
		err = DownloadAndInstallPython(pythonVersion, osType)
		if err != nil {
			return []string{}, fmt.Errorf("issue installing Python version: %w", err)
		}
	}
	return installPythonVersions, nil
}

// ScanAndHandlePythonVersions scans for Python versions, handles result/errors and creates PythonConfig
func ScanAndHandlePythonVersions(osType config.OperatingSystem) error {
	//pythonVersionsOrig, err := ScanForPythonVersions()
	//if err != nil {
	//	return []string{}, fmt.Errorf("issue occured in scanning for Python versions: %w", err)
	//}

	//fmt.Println("\nFound Python versions: ", strings.Join(pythonVersionsOrig, ", "), "\n")

	_, _ = PromptAndInstallPython(osType)

	return nil
}

// ScanForPythonVersions scans for Python versions in locations workbench will also look

// Prompt users if they would like to install Python versions

// PythonPATHPrompt asks users if they would like to set Python PATH

func RetrieveValidPythonVersions() ([]string, error) {
	// TODO make this dynamic based on https://cdn.posit.co/python/versions.json
	return availablePythonVersions, nil
}

// Prompt asking users which Python version(s) they would like to install
func PythonSelectVersionsPrompt(availablePythonVersions []string) ([]string, error) {
	var qs = []*survey.Question{
		{
			Name: "pythonVersions",
			Prompt: &survey.MultiSelect{
				Message: "Which version(s) of Python would you like to install?",
				Options: availablePythonVersions,
				Default: availablePythonVersions[0],
			},
		},
	}
	pythonVersionsAnswers := struct {
		PythonVersions []string `survey:"pythonVersions"`
	}{}
	err := survey.Ask(qs, &pythonVersionsAnswers)
	if err != nil {
		return []string{}, errors.New("there was an issue with the Python versions selection prompt")
	}
	if len(pythonVersionsAnswers.PythonVersions) == 0 {
		return []string{}, errors.New("at least one Python version must be selected")
	}
	return pythonVersionsAnswers.PythonVersions, nil
}

// Downloads the Python installer, and installs Python
func DownloadAndInstallPython(pythonVersion string, osType config.OperatingSystem) error {
	// Create InstallerInfoPython with the proper information
	installerInfo, err := PopulateInstallerInfo("python", pythonVersion, osType)
	if err != nil {
		return fmt.Errorf("PopulateInstallerInfoPython: %w", err)
	}
	// Download installer
	downloadurl := "Python Version Download URL: " + installerInfo.URL
	fmt.Println(downloadurl)
	if err != nil {
		return fmt.Errorf("DownloadPython: %w", err)
	}
	// Install Python
	err = install.InstallLanguage(osType)
	if err != nil {
		return fmt.Errorf("InstallLanguage: %w", err)
	}
	// Upgrade pip, setuptools, and wheel
	//err = UpgradePythonTools(pythonVersion)
	//if err != nil {
	//	return fmt.Errorf("UpgradePythonTools: %w", err)
	//}

	return nil
}

// RemovePythonFromPath removes python or python3 from the end of a path so the directory can be used

// RemovePythonFromPathSlice removes python or python3 from the end of a set of path strings in a slice so the directories can be used
