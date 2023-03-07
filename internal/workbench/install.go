package workbench

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/samcofer/tam-offline-download-email/internal/config"
)

// InstallerInfo contains the information needed to download and install Workbench
type InstallerInfo struct {
	BaseName string `json:"basename"`
	URL      string `json:"url"`
	Version  string `json:"version"`
	Label    string `json:"label"`
}

// OperatingSystems contains the installer information for each supported operating system
type OperatingSystems struct {
	Bionic   InstallerInfo `json:"bionic"`
	Ubuntu64 InstallerInfo `json:"ubuntu64"`
	Focal    InstallerInfo `json:"focal"`
	Jammy    InstallerInfo `json:"jammy"`
	Redhat7  InstallerInfo `json:"redhat7_64"`
	RHEL8    InstallerInfo `json:"rhel8"`
	Redhat8  InstallerInfo `json:"redhat8"`
	Fedora28 InstallerInfo `json:"fedora28"`
}

// Installer contains the installer information for a product
type Installer struct {
	Installer OperatingSystems `json:"installer"`
}

// ProductType contains the installer for each product type
type ProductType struct {
	Server Installer `json:"server"`
}

// Category contains information for stable and preview product types
type Category struct {
	Stable ProductType `json:"stable"`
}

// Product contains information for each RStudio product
type Product struct {
	Pro Category `json:"pro"`
}

// RStudio contains product information
type RStudio struct {
	Rstudio Product   `json:"rstudio"`
	Connect Installer `json:"connect"`
	Rspm    Installer `json:"rspm"`
}

// Retrieves JSON data from Posit, downloads the Workbench installer, and installs Workbench
func DownloadAndInstallWorkbench(osType config.OperatingSystem) (string, string, string, error) {
	// Retrieve JSON data

	rstudio, err := RetrieveWorkbenchInstallerInfo()
	if err != nil {
		return "error", "error", "error", fmt.Errorf("RetrieveWorkbenchInstallerInfo: %w", err)
	}
	// Retrieve installer info
	installerInfoWorkbench, installerInfoConnect, installerInfoPackagemanager, err := rstudio.GetInstallerInfo(osType)
	if err != nil {
		return "error", "error", "error", fmt.Errorf("GetInstallerInfo: %w", err)
	}
	// Download installer
	WorkbenchDownLoad := "Workbench download URL: " + installerInfoWorkbench.URL
	ConnectDownLoad := "Connect download URL: " + installerInfoConnect.URL
	PackageManagerDownLoad := "Package Manager download URL: " + installerInfoPackagemanager.URL
	return WorkbenchDownLoad, ConnectDownLoad, PackageManagerDownLoad, nil
}

// Installs Workbench in a certain way based on the operating system

// Creates the proper command to install Workbench based on the operating system

// Pulls out the installer information from the JSON data based on the operating system
func (r *RStudio) GetInstallerInfo(osType config.OperatingSystem) (InstallerInfo, InstallerInfo, InstallerInfo, error) {
	switch osType {
	case config.Ubuntu18:
		return r.Rstudio.Pro.Stable.Server.Installer.Bionic, r.Connect.Installer.Bionic, r.Rspm.Installer.Ubuntu64, nil
	case config.Ubuntu20:
		return r.Rstudio.Pro.Stable.Server.Installer.Bionic, r.Connect.Installer.Focal, r.Rspm.Installer.Focal, nil
	case config.Ubuntu22:
		return r.Rstudio.Pro.Stable.Server.Installer.Jammy, r.Connect.Installer.Jammy, r.Rspm.Installer.Jammy, nil
	case config.Redhat7:
		return r.Rstudio.Pro.Stable.Server.Installer.Redhat7, r.Connect.Installer.Redhat7, r.Rspm.Installer.Redhat7, nil
	case config.Redhat8:
		return r.Rstudio.Pro.Stable.Server.Installer.RHEL8, r.Connect.Installer.Redhat8, r.Rspm.Installer.Fedora28, nil
	default:
		return InstallerInfo{}, InstallerInfo{}, InstallerInfo{}, errors.New("operating system not supported")
	}
}

// Retrieves JSON data from Posit
func RetrieveWorkbenchInstallerInfo() (RStudio, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodGet, "https://www.rstudio.com/wp-content/downloads.json", nil)
	if err != nil {
		return RStudio{}, errors.New("error creating request")
	}
	res, err := client.Do(req)
	if err != nil {
		return RStudio{}, errors.New("error retrieving JSON data")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return RStudio{}, errors.New("error retrieving JSON data")
	}
	var rstudio RStudio
	err = json.NewDecoder(res.Body).Decode(&rstudio)
	if err != nil {
		return RStudio{}, errors.New("error unmarshalling JSON data")
	}
	return rstudio, nil
}
