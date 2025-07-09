package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
)

// DecompressMiddleware распаковывает тело запроса, если оно сжато в gzip
func DecompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to decompress gzip body", http.StatusBadRequest)
				return
			}
			defer gz.Close()

			// Заменяем тело запроса на распакованное
			r.Body = io.NopCloser(gz)
		}
		next.ServeHTTP(w, r)
	})
}
