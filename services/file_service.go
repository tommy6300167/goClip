package services

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"clipmini/models"
)

type FileService struct {
	config *models.AppConfig
}

func NewFileService(config *models.AppConfig) *FileService {
	return &FileService{
		config: config,
	}
}

func (fs *FileService) ReadHistoryLines() ([]string, error) {
	f, err := os.Open(fs.config.LogFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		if !utf8.ValidString(line) {
			line = strings.ToValidUTF8(line, "?")
		}
		lines = append(lines, line)
	}
	return lines, nil
}

func (fs *FileService) WriteHistoryLines(lines []string) error {
	if err := os.MkdirAll(fs.config.LogDirPath, 0o755); err != nil {
		return err
	}
	
	f, err := os.Create(fs.config.LogFilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	
	for _, line := range lines {
		if _, err := f.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return nil
}

func (fs *FileService) SaveImage(data []byte, timestamp string) (string, error) {
	if err := os.MkdirAll(fs.config.ImageDirPath, 0o755); err != nil {
		return "", err
	}
	
	filename := fmt.Sprintf("image_%s.png",
		strings.ReplaceAll(strings.ReplaceAll(timestamp, ":", "-"), " ", "_"))
	path := filepath.Join(fs.config.ImageDirPath, filename)
	
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}
	return path, nil
}

func (fs *FileService) CleanupImageFiles(imagePaths []string) {
	for _, path := range imagePaths {
		_ = os.Remove(path)
	}
}

func (fs *FileService) DeleteHistoryFile() error {
	return os.Remove(fs.config.LogFilePath)
}

func (fs *FileService) DeleteImageDirectory() error {
	return os.RemoveAll(fs.config.ImageDirPath)
}

func (fs *FileService) ImageExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}