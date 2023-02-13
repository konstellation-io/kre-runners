package main

import (
	"krt-files-downloader/v2/internal/adapter/config"
	"krt-files-downloader/v2/internal/adapter/repository/mongodb"
	"krt-files-downloader/v2/internal/domain/usecase"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}

	krtRepository := mongodb.NewMongoKRTRepository(cfg)
	krtInteractor := usecase.NewKRTInteractor(krtRepository)

	err = krtInteractor.DownloadKRTFiles(cfg.RuntimeID, cfg.VersionID, cfg.KRTBasePath)
	if err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}

	log.Print("All done")
}
