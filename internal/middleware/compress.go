package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}

type gzipReader struct {
	r  io.ReadCloser
	gz *gzip.Reader
}

func newGzipReader(r io.ReadCloser) (*gzipReader, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &gzipReader{r: r, gz: gz}, nil
}

func (gz *gzipReader) Read(data []byte) (int, error) {
	return gz.gz.Read(data)
}

func (gz *gzipReader) Close() error {
	if err := gz.r.Close(); err != nil {
		return err
	}
	return gz.gz.Close()
}
func CompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cntenc := strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
		accEnc := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
		cntjson := strings.Contains(r.Header.Get("Content-Type"), "application/json")
		conttxt := strings.Contains(r.Header.Get("Content-Type"), "text/html")

		if !(cntjson || conttxt) {
			next.ServeHTTP(w, r)
			return
		}

		if cntenc {
			gr, err := newGzipReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = gr
			defer gr.Close()
		}

		if accEnc {
			w.Header().Set("Content-Encoding", "gzip")
			gz, err := gzip.NewWriterLevel(w, gzip.BestCompression)
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}
			next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
			return
		}
		next.ServeHTTP(w, r)

	})
}
