package middware

import (
	"compress/gzip"
	"github.com/VicShved/shorturl/internal/logger"
	"go.uber.org/zap"
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
		logger.Log.Error("new gzip reader", zap.Error(err))
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
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cntEnc := strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
		accEnc := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
		cntJSON := strings.Contains(r.Header.Get("Content-Type"), "application/json")
		contTxt := strings.Contains(r.Header.Get("Content-Type"), "text/html")
		logger.Log.Info("GzipMiddleware", zap.Bool("cntEnc", cntEnc),
			zap.Bool("accEnc", accEnc),
			zap.Bool("cntJSON", cntJSON),
			zap.Bool("contTxt", contTxt),
		)

		if cntEnc {
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
			defer gz.Close()
			next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
