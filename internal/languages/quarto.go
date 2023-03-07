package languages

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/samcofer/tam-offline-download-email/internal/config"
	"net/http"
	"time"
)

// RStudio contains product information
type Assets []struct {
	BrowserDownloadURL string `json:"browser_download_url"`
}
type Quarto struct {
	Assets Assets `json:"assets"`
}

func RetrieveQuartoInstallerInfo() (Quarto, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodGet, "https://api.github.com/repos/quarto-dev/quarto-cli/releases/latest", nil)
	if err != nil {
		return Quarto{}, errors.New("error creating request")
	}
	res, err := client.Do(req)
	if err != nil {
		return Quarto{}, errors.New("error retrieving JSON data")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return Quarto{}, errors.New("error retrieving JSON data")
	}
	var quarto Quarto
	err = json.NewDecoder(res.Body).Decode(&quarto)
	if err != nil {
		return Quarto{}, errors.New("error unmarshalling JSON data")
	}
	return quarto, nil
}

func DownloadAndInstallQuarto(osType config.OperatingSystem) (string, error) {
	// Retrieve JSON data

	quarto, err := RetrieveQuartoInstallerInfo()
	if err != nil {
		return "error", fmt.Errorf("RetrieveWorkbenchInstallerInfo: %w", err)
	}
	// Retrieve installer info
	QuartoDownload, err := quarto.GetInstallerInfo(osType)
	if err != nil {
		return "error", fmt.Errorf("GetInstallerInfo: %w", err)
	}

	return QuartoDownload, nil
}

func (q *Quarto) GetInstallerInfo(osType config.OperatingSystem) (string, error) {
	switch osType {
	case config.Ubuntu18, config.Ubuntu20, config.Ubuntu22, config.Redhat8:
		for _, val := range q.Assets {
			fmt.Println(val.BrowserDownloadURL)
		}
		return q.Assets[0].BrowserDownloadURL, nil
	case config.Redhat7:
		for _, val := range q.Assets {
			fmt.Println(val.BrowserDownloadURL)
		}
		return q.Assets[0].BrowserDownloadURL, nil
	default:
		return "", errors.New("operating system not supported")
	}
}
