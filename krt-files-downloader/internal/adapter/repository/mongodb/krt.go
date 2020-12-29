package mongodb

import (
	"fmt"
	"io"
	"krt-files-downloader/internal/adapter/config"
	"log"
)

type MongoKRTRepository struct {
	cfg     *config.Config
	manager *MongoManager
}

func NewMongoKRTRepository(cfg *config.Config) *MongoKRTRepository {
	manager := NewMongoManger()
	return &MongoKRTRepository{manager: manager, cfg: cfg}
}

func (m MongoKRTRepository) DownloadKRT(versionID string) (io.Reader, error) {
	err := m.manager.Connect(m.cfg.MongoDB.URI)
	if err != nil {
		return nil, fmt.Errorf("connecting to mongodb: %w", err)
	}

	defer func() {
		if err := m.manager.Disconnect(); err != nil {
			log.Printf("Error disconecting from MongoDB: %s", err)
		}
	}()

	return m.manager.DownloadFile(m.cfg.MongoDB.DBName, m.cfg.MongoDB.Bucket, versionID)
}
