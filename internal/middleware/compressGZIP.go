package middleware

import (
	"bytes"
	"compress/gzip"
)

// CompressGZIP сжимает []byte в формат gzip
func CompressGZIP(data []byte) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write(data)
	if err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return &buf, nil
}
