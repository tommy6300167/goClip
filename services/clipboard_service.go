package services

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os/exec"
	"strings"
	"unicode/utf8"
)

type ClipboardService struct{}

func NewClipboardService() *ClipboardService {
	return &ClipboardService{}
}

func (cs *ClipboardService) HasImageInClipboard() bool {
	cmd := exec.Command("/usr/bin/osascript", "-e", "return (clipboard info) as string")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return false
	}
	s := out.String()
	return strings.Contains(s, "«class PNGf»") || strings.Contains(s, "«class TIFF»")
}

func (cs *ClipboardService) ReadClipboardImage() ([]byte, error) {
	try := func(typ string) ([]byte, error) {
		cmd := exec.Command("/usr/bin/osascript",
			"-e", fmt.Sprintf(`set d to the clipboard as «class %s»`, typ),
			"-e", `return d`)
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			return nil, err
		}
		hexStr := strings.TrimSpace(out.String())
		hexStr = strings.TrimPrefix(hexStr, "«data "+typ)
		hexStr = strings.TrimSuffix(hexStr, "»")
		hexStr = strings.ReplaceAll(hexStr, " ", "")
		if len(hexStr)%2 != 0 {
			return nil, fmt.Errorf("bad hex length")
		}
		data := make([]byte, len(hexStr)/2)
		for i := 0; i < len(hexStr); i += 2 {
			var b byte
			fmt.Sscanf(hexStr[i:i+2], "%02x", &b)
			data[i/2] = b
		}
		return data, nil
	}
	if b, err := try("PNGf"); err == nil && len(b) > 0 {
		return b, nil
	}
	return try("TIFF")
}

func (cs *ClipboardService) ReadClipboardText() (string, error) {
	cmd := cs.utf8Env(exec.Command("/usr/bin/pbpaste", "-Prefer", "txt"))
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	result := strings.TrimRight(out.String(), "\r\n")
	if !utf8.ValidString(result) {
		result = strings.ToValidUTF8(result, "�")
	}
	return result, nil
}

func (cs *ClipboardService) CopyTextToClipboard(text string) error {
	cmd := exec.Command("/usr/bin/pbcopy")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

func (cs *ClipboardService) CopyImageToClipboard(imagePath string) error {
	script := fmt.Sprintf(`
		set imageData to read (POSIX file "%s") as «class PNGf»
		set the clipboard to imageData
	`, imagePath)
	cmd := exec.Command("/usr/bin/osascript", "-e", script)
	return cmd.Run()
}

func (cs *ClipboardService) ClearSystemClipboard() error {
	return exec.Command("/usr/bin/pbcopy").Run()
}

func (cs *ClipboardService) GetImageHash(data []byte) string {
	h := sha1.Sum(data)
	return hex.EncodeToString(h[:])
}

func (cs *ClipboardService) utf8Env(cmd *exec.Cmd) *exec.Cmd {
	cmd.Env = append(cmd.Env,
		"LANG=zh_TW.UTF-8",
		"LC_ALL=zh_TW.UTF-8",
		"LC_CTYPE=zh_TW.UTF-8",
	)
	return cmd
}