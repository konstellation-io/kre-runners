package repository

import (
	"io"
)

//go:generate mockgen -source=${GOFILE} -destination=$PWD/internal/mocks/repo_${GOFILE} -package=mocks

type KRTRepository interface {
	DownloadKRT(versionID string) (io.Reader, error)
}
