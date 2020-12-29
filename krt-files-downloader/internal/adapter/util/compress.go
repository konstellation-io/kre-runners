package util

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func ExtractTarGz(input io.Reader, destinationFolder string) error {
	start := time.Now()

	log.Printf("Extracting KRT file to \"%s\"...", destinationFolder)

	err := os.MkdirAll(destinationFolder, os.ModePerm)
	if err != nil {
		return fmt.Errorf("creating destination folders: %w", err)
	}

	uncompressedStream, err := gzip.NewReader(input)
	if err != nil {
		return fmt.Errorf("uncompressing tar.gz file: %w", err)
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("executing tarReader next: %w", err)
		}

		target := filepath.Join(destinationFolder, header.Name) // nolint: gosec

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.ModePerm); err != nil {
				return fmt.Errorf("creating dir: %w", err)
			}
		case tar.TypeReg:
			outFile, err := os.Create(target)
			if err != nil {
				return fmt.Errorf("creating file: %w", err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil { // nolint: gosec
				return fmt.Errorf("copying file content: %w", err)
			}

			err = outFile.Close()
			if err != nil {
				return fmt.Errorf("closing out file: %w", err)
			}
		}
	}

	elapsed := time.Since(start)
	log.Printf("Extract took %s", elapsed)

	return nil
}
