package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func DeviceFingerprint() (string, error) {
	var components []string

	components = append(components, runtime.GOOS)
	components = append(components, runtime.GOARCH)

	if runtime.GOOS == "darwin" {
		if serial, err := getMacSerial(); err == nil {
			components = append(components, serial)
		}

		if uuid, err := getMacHardwareUUID(); err == nil {
			components = append(components, uuid)
		}
	}

	if username := os.Getenv("USER"); username != "" {
		components = append(components, username)
	}

	if hostname, err := os.Hostname(); err == nil {
		components = append(components, hostname)
	}

	combined := strings.Join(components, "|")
	hash := sha256.Sum256([]byte(combined))

	return hex.EncodeToString(hash[:]), nil
}

func getMacSerial() (string, error) {
	cmd := exec.Command("system_profiler", "SPHardwareDataType")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Serial Number") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}

	return "", fmt.Errorf("serial number not found")
}

func getMacHardwareUUID() (string, error) {
	cmd := exec.Command("system_profiler", "SPHardwareDataType")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Hardware UUID") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}

	return "", fmt.Errorf("hardware UUID not found")
}

func IsValidDeviceFingerprint(storedFingerprint string) bool {
	currentFingerprint, err := DeviceFingerprint()
	if err != nil {
		return false
	}

	return currentFingerprint == storedFingerprint
}
