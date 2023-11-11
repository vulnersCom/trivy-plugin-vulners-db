package internal

import (
	"github.com/aquasecurity/trivy-db/pkg/metadata"
	"github.com/cavaliergopher/grab/v3"
	"log"
	"path/filepath"
	"time"
)

func Download(cacheDir string, apiKey string) {
	dbPath := filepath.Join(cacheDir, "db")
	client := grab.NewClient()
	req, _ := grab.NewRequest(dbPath, "https://vulners.com/api/v3/trivy/free?apiKey="+apiKey)

	log.Printf("Downloading: %v", req.URL())
	resp := client.Do(req)
	log.Printf("Response statys: %v", resp.HTTPResponse.Status)

	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			log.Printf("Transferred: %v bytes", resp.BytesComplete())
		case <-resp.Done:
			break Loop
		}
	}

	if err := resp.Err(); err != nil {
		log.Fatalf("Download failed: %v", err)
	}

	log.Printf("Download saved to: %v ", resp.Filename)

	err := extractTarGz(resp.Filename, dbPath)
	if err != nil {
		log.Fatalf("Can't extract: %v", err)
	}

	dbMetadata := createMetadata()

	dbMetadataClient := metadata.NewClient(cacheDir)
	err = dbMetadataClient.Update(dbMetadata)
	if err != nil {
		log.Fatalf("Update dbMetadata failed: %v", err)
	}

}

func createMetadata() metadata.Metadata {
	dbMetadata := metadata.Metadata{}
	dbMetadata.Version = 2
	dbMetadata.DownloadedAt = time.Now().UTC()
	dbMetadata.UpdatedAt = time.Now().UTC()
	dbMetadata.NextUpdate = time.Now().UTC().AddDate(1, 0, 0)
	return dbMetadata
}
