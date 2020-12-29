package mongodb

import (
	"bytes"
	"io"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

func (m *MongoManager) DownloadFile(dbName, bucketName, fileID string) (io.Reader, error) {
	start := time.Now()

	bucket, err := gridfs.NewBucket(
		m.client.Database(dbName),
		options.GridFSBucket().SetName(bucketName),
	)
	if err != nil {
		return nil, err
	}

	log.Printf("Downloading \"%s\" file from \"%s\" bucket...", fileID, bucketName)

	var buf bytes.Buffer

	dSize, err := bucket.DownloadToStream(fileID, &buf)
	if err != nil {
		return nil, err
	}

	log.Printf("Downloaded %d bytes", dSize)

	elapsed := time.Since(start)
	log.Printf("Download took %s", elapsed)

	return &buf, nil
}
