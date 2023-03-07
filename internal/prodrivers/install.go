package prodrivers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/samcofer/tam-offline-download-email/internal/config"
	"net/http"
	"time"
)

// InstallerInfo contains the information needed to download and install Posit Pro Drivers
type InstallerInfo struct {
	BaseName string `json:"basename"`
	URL      string `json:"url"`
	Version  string `json:"version"`
	Label    string `json:"label"`
}

// OperatingSystems contains the installer information for each supported operating system
type OperatingSystems struct {
	// Posit Pro Drivers are the same for all Ubuntu versions, so we only need one
	Focal   InstallerInfo `json:"focal"`
	Redhat7 InstallerInfo `json:"redhat7_64"`
	Redhat8 InstallerInfo `json:"rhel8"`
}

// Installer contains the installer information for a product
type Installer struct {
	Installer OperatingSystems `json:"installer"`
}

// ProDrivers contains product information
type ProDrivers struct {
	ProDrivers Installer `json:"pro_drivers"`
}

// DownloadAndInstallProDrivers Retrieves JSON data from Posit, downloads the Pro Drivers installer, and installs Pro Drivers
func DownloadAndInstallProDrivers(osType config.OperatingSystem) error {
	// Retrieve JSON data
	rstudio, err := RetrieveProDriversInstallerInfo()
	if err != nil {
		return fmt.Errorf("RetrieveProDriversInstallerInfo: %w", err)
	}
	// Retrieve installer info
	installerInfo, err := rstudio.GetInstallerInfo(osType)
	if err != nil {
		return fmt.Errorf("GetInstallerInfo: %w", err)
	}
	// Install prerequisites
	//err = InstallUnixODBC(osType)
	//if err != nil {
	//	return fmt.Errorf("InstallUnixODBC: %w", err)
	//}
	// Download installer
	//install.DownloadFile("Pro Drivers", installerInfo.URL, installerInfo.BaseName)

	fmt.Println("Pro Driver Download URL: " + installerInfo.URL)
	if err != nil {
		return fmt.Errorf("DownloadFile: %w", err)
	}
	// Install Pro Drivers
	//err = InstallProDrivers(filepath, osType)
	//if err != nil {
	//	return fmt.Errorf("InstallProDrivers: %w", err)
	//}
	// Configure ODBC driver name and locations
	//err = BackupAndAppendODBCConfiguration()
	//if err != nil {
	//	return fmt.Errorf("BackupAndAppendODBCConfiguration: %w", err)
	//}
	//fmt.Println("\nPosit Pro Drivers next steps:\nNow that the Pro Drivers are installed and /etc/odbcinst.ini is set up, the next step is to test database connectivity and/or create DSNs in your /etc/odbc.ini file.\n\n More information about each of these steps can be found here: https://docs.posit.co/pro-drivers/workbench-connect/#step-4-testing-database-connectivity\n")
	return nil
}

// RetrieveProDriversInstallerInfo Retrieves JSON data from Posit
func RetrieveProDriversInstallerInfo() (ProDrivers, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodGet, "https://www.rstudio.com/wp-content/downloads.json", nil)
	if err != nil {
		return ProDrivers{}, errors.New("error creating request")
	}
	res, err := client.Do(req)
	if err != nil {
		return ProDrivers{}, errors.New("error retrieving JSON data")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return ProDrivers{}, errors.New("error retrieving JSON data")
	}
	var proDrivers ProDrivers
	err = json.NewDecoder(res.Body).Decode(&proDrivers)
	if err != nil {
		return ProDrivers{}, errors.New("error unmarshalling JSON data")
	}
	return proDrivers, nil
}

// GetInstallerInfo Pulls out the installer information from the JSON data based on the operating system
func (pd *ProDrivers) GetInstallerInfo(osType config.OperatingSystem) (InstallerInfo, error) {
	switch osType {
	// Posit Pro Drivers are the same for all Ubuntu versions
	case config.Ubuntu18, config.Ubuntu20, config.Ubuntu22:
		return pd.ProDrivers.Installer.Focal, nil
	case config.Redhat7:
		return pd.ProDrivers.Installer.Redhat7, nil
	case config.Redhat8:
		return pd.ProDrivers.Installer.Redhat8, nil
	default:
		return InstallerInfo{}, errors.New("operating system not supported")
	}
}
