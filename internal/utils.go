package internal

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"os"
	"path/filepath"
)

func extractTarGz(filename string, destination string) error {
	gzipStream, err := os.Open(filename)
	if err != nil {
		log.Fatal("Can't open file")
		return err
	}

	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		log.Fatal("NewReader failed")
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Next failed: %s", err.Error())
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(header.Name, 0755); err != nil {
				log.Fatalf("Mkdir failed: %s", err.Error())
				return err
			}
		case tar.TypeReg:
			destinationFilename := filepath.Join(destination, "trivy.db")

			_, err := os.Stat(destinationFilename)
			if !os.IsNotExist(err) {
				_ = os.Remove(destinationFilename)
			}

			outputFile, err := os.Create(destinationFilename)
			if err != nil {
				log.Fatalf("Create failed: %s", err.Error())
				return err
			}
			if _, err := io.Copy(outputFile, tarReader); err != nil {
				log.Fatalf("Copy failed: %s", err.Error())
				return err
			}

			_ = outputFile.Close()

		default:
			log.Fatalf(
				"Uknown type: %s in %s",
				header.Typeflag,
				header.Name)
		}

	}

	_ = gzipStream.Close()

	_ = os.Remove(filename)

	return nil
}
