package usecase_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"krt-files-downloader/internal/domain/usecase"
	"krt-files-downloader/internal/mocks"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type krtSuite struct {
	ctrl          *gomock.Controller
	krtInteractor *usecase.KRTInteractor
	krtRepository *mocks.MockKRTRepository
}

func newKRTSuite(t *testing.T) *krtSuite {
	ctrl := gomock.NewController(t)

	krtRepository := mocks.NewMockKRTRepository(ctrl)
	krtInteractor := usecase.NewKRTInteractor(krtRepository)

	return &krtSuite{
		ctrl:          ctrl,
		krtInteractor: krtInteractor,
		krtRepository: krtRepository,
	}
}

func TestKRTInteractor_DownloadKRTFiles(t *testing.T) {
	s := newKRTSuite(t)
	defer s.ctrl.Finish()

	const (
		runtimeID   = "runtime1234"
		versionID   = "version1234"
		runtimeID   = "runtime1234"
		destination = "_test/output"
		krtFile     = "_test/greeter-v1.krt"
	)

	expectedFiles := [...]string{
		"krt.yml",
		"public_input.proto",
		"src/py-greeter/main.py",
		"src/py-greeter/public_input_pb2.py",
	}

	defer os.RemoveAll(destination)

	krtBuf, err := ioutil.ReadFile(krtFile)
	require.Nil(t, err)

	krtReader := bytes.NewReader(krtBuf)
	s.krtRepository.EXPECT().DownloadKRT(runtimeID, versionID).Return(krtReader, nil)

	err = s.krtInteractor.DownloadKRTFiles(runtimeID, versionID, destination)
	require.Nil(t, err)

	for _, f := range expectedFiles {
		if _, err := os.Stat(fmt.Sprintf("%s/%s", destination, f)); os.IsNotExist(err) {
			t.Fatalf("file \"%s\" not found in the destination folder", f)
		}
	}
}
