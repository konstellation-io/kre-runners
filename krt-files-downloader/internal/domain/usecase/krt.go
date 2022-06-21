package usecase

import (
	"fmt"
	"krt-files-downloader/internal/adapter/util"
	"krt-files-downloader/internal/domain/repository"
)

type KRTInteractor struct {
	krtRepo repository.KRTRepository
}

func NewKRTInteractor(krtRepo repository.KRTRepository) *KRTInteractor {
	return &KRTInteractor{krtRepo: krtRepo}
}

func (i KRTInteractor) DownloadKRTFiles(runtimeID, versionID, destination string) error {
	krt, err := i.krtRepo.DownloadKRT(runtimeID, versionID)
	if err != nil {
		return fmt.Errorf("downloading the KRT file: %w", err)
	}

	err = util.ExtractTarGz(krt, destination)
	if err != nil {
		return fmt.Errorf("extracting the KRT file: %w", err)
	}

	return nil
}
