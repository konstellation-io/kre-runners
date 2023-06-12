package common

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

const (
	CompressLevel = 9
	gzipID1       = 0x1f
	gzipID2       = 0x8b
)

// IsCompressed check if the input string is compressed.
func IsCompressed(data []byte) bool {
	return data[0] == gzipID1 && data[1] == gzipID2
}

// CompressData creates compressed []byte.
func CompressData(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz, err := gzip.NewWriterLevel(&b, CompressLevel)
	if err != nil {
		return nil, err
	}
	_, err = gz.Write(data)
	if err != nil {
		return nil, err
	}
	err = gz.Close()
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// UncompressData open gzip and return uncompressed []byte.
func UncompressData(data []byte) ([]byte, error) {
	rd := bytes.NewReader(data)
	gr, err := gzip.NewReader(rd)
	if err != nil {
		return nil, err
	}

	defer gr.Close()
	return ioutil.ReadAll(gr)
}
