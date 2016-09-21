// zip
package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func ZipFiles(files []string, outFilePath string) {
	outFile, _ := os.Create(outFilePath)
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)

	for _, file := range files {
		f, _ := os.Open(file)
		archivePath := filepath.Base(file)
		zipFileWriter, _ := zipWriter.Create(archivePath)
		_, _ = io.Copy(zipFileWriter, f)
		f.Close()
		os.Remove(file)
	}
	zipWriter.Close()
}
