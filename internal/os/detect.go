package os

import (
	"errors"
	"github.com/samcofer/tam-offline-download-email/internal/config"
)

// DetectOS Detect which operating system WBI is running on
func DetectOS(OS string) (config.OperatingSystem, error) {

	var err error
	switch OS {
	case "ubuntu22":
		return config.Ubuntu22, err
	case "ubuntu20":
		return config.Ubuntu20, err
	case "ubuntu18":
		return config.Ubuntu18, err
	case "rhel7":
		return config.Redhat7, err
	case "rhel8":
		return config.Redhat8, err
	default:
		return config.Unknown, errors.New("unsupported OS")
	}
}
