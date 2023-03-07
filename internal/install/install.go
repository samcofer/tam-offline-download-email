package install

import (
	"errors"
	"fmt"
	"github.com/samcofer/tam-offline-download-email/internal/config"
)

// Installs R/Python in a certain way based on the operating system
func InstallLanguage(osType config.OperatingSystem) error {

	_, err := RetrieveInstallCommand(osType)
	//print("Here is the install command" + installCommand)
	if err != nil {
		return fmt.Errorf("RetrieveInstallCommand: %w", err)
	}

	//err = system.RunCommand(installCommand)
	//if err != nil {
	//	return fmt.Errorf("RunCommand: %w", err)
	//}

	//successMessage := "\n" + languageTitleCase + " version " + version + " successfully installed!\n"
	//fmt.Println(successMessage)
	return nil
}

// Creates the proper command to install R/Python based on the operating system
func RetrieveInstallCommand(osType config.OperatingSystem) (string, error) {

	filepath := "/path/to/file"

	switch osType {
	case config.Ubuntu22, config.Ubuntu20, config.Ubuntu18:
		return "gdebi -n " + filepath, nil
	case config.Redhat7, config.Redhat8:
		return "yum install -y " + filepath, nil
	default:
		return "", errors.New("operating system not supported")
	}
}
