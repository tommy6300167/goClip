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

// DeviceFingerprint 獲取設備指紋
func DeviceFingerprint() (string, error) {
	var components []string

	// 1. 系統架構
	components = append(components, runtime.GOOS)
	components = append(components, runtime.GOARCH)

	// 2. 硬體序號 (macOS)
	if runtime.GOOS == "darwin" {
		if serial, err := getMacSerial(); err == nil {
			components = append(components, serial)
		}

		if uuid, err := getMacHardwareUUID(); err == nil {
			components = append(components, uuid)
		}
	}

	// 3. 使用者名稱
	if username := os.Getenv("USER"); username != "" {
		components = append(components, username)
	}

	// 4. 主機名稱
	if hostname, err := os.Hostname(); err == nil {
		components = append(components, hostname)
	}

	// 組合所有組件並生成雜湊
	combined := strings.Join(components, "|")
	hash := sha256.Sum256([]byte(combined))

	return hex.EncodeToString(hash[:]), nil
}

// getMacSerial 獲取Mac序號
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

// getMacHardwareUUID 獲取Mac硬體UUID
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

// IsValidDeviceFingerprint 檢查設備指紋是否匹配
func IsValidDeviceFingerprint(storedFingerprint string) bool {
	currentFingerprint, err := DeviceFingerprint()
	if err != nil {
		return false
	}

	return currentFingerprint == storedFingerprint
}
