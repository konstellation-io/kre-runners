package repository

import (
	"io"
)

//go:generate mockgen -source=${GOFILE} -destination=../../mocks/repo_${GOFILE} -package=mocks

type KRTRepository interface {
	DownloadKRT(runtimeID, versionID string) (io.Reader, error)
}
