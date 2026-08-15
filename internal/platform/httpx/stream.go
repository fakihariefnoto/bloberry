package httpx

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var ErrBodyTooLarge = errors.New("request body exceeds size ceiling")

// MaxBytesReader wraps the request body with a hard ceiling (TRD R5).
func MaxBytesReader(w http.ResponseWriter, r *http.Request, n int64) {
	r.Body = http.MaxBytesReader(w, r.Body, n)
}

// ServeStream writes bytes to the response with HTTP range-request support.
// Never buffers the whole object in memory (R5).
func ServeStream(w http.ResponseWriter, r *http.Request, content io.ReadSeeker, contentType string, size int64, modtime time.Time) {
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", contentType)
	http.ServeContent(w, r, "", modtime, content)
}

// Copy is the streaming copy used on the write path; the error must be
// handled by callers (gosec G110 — no unchecked io.Copy).
func Copy(dst io.Writer, src io.Reader) error {
	_, err := io.Copy(dst, src)
	return err
}

// ParseRangeBytes returns (start, end) for a Range header. end is inclusive.
func ParseRangeBytes(header string, size int64) (start, end int64, ok bool) {
	const prefix = "bytes="
	if !strings.HasPrefix(header, prefix) {
		return 0, 0, false
	}
	spec := strings.TrimPrefix(header, prefix)
	if strings.Contains(spec, ",") {
		return 0, 0, false // multiple ranges unsupported
	}
	parts := strings.SplitN(spec, "-", 2)
	if len(parts) != 2 {
		return 0, 0, false
	}
	if parts[0] == "" {
		// suffix range: last N bytes
		n, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil || n <= 0 {
			return 0, 0, false
		}
		if n > size {
			n = size
		}
		return size - n, size - 1, true
	}
	start, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || start < 0 || start >= size {
		return 0, 0, false
	}
	if parts[1] == "" {
		return start, size - 1, true
	}
	e, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || e < start {
		return 0, 0, false
	}
	if e >= size {
		e = size - 1
	}
	return start, e, true
}
