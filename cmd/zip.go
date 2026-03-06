package cmd

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// downloadAndExtractZip downloads a ZIP from a URL and extracts it to dest.
// It uses downloadFile from add.go
func downloadAndExtractZip(url, dest string) error {
	data, err := downloadFile(url)
	if err != nil {
		return err
	}

	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return fmt.Errorf("failed to parse zip archive: %w", err)
	}

	// Create dest folder
	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}
	destInfo, err := os.Stat(dest)
	if err != nil {
		return err
	}
	if !destInfo.IsDir() {
		return fmt.Errorf("destination is not a directory")
	}

	destClean := filepath.Clean(dest)

	for _, file := range reader.File {
		// Prevent Zip Slip vulnerability
		path := filepath.Join(dest, file.Name)
		if !strings.HasPrefix(path, destClean+string(os.PathSeparator)) && destClean != path {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		if err = os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		content, err := io.ReadAll(rc)
		if err != nil {
			outFile.Close()
			rc.Close()
			return err
		}

		// Apply Go template engine on the extracted file
		content = processTemplate(content)

		_, err = outFile.Write(content)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
