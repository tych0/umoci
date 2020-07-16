package mutate

import (
	"io"
	"runtime"

	gzip "github.com/klauspost/pgzip"
	ispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

type Compressor interface {
	Compress(io.Reader) (io.ReadCloser, error)
	MediaType() string
}

type withClose struct {
	r io.Reader
}

func (wc withClose) Read(p []byte) (n int, err error) {
	return wc.r.Read(p)
}

func (wc withClose) Close() error {
	return nil
}

type noopCompressor struct {
	mediaType string
}

func NewNoopCompressor(mediaType string) Compressor {
	return noopCompressor{mediaType}
}

func (nc noopCompressor) Compress(r io.Reader) (io.ReadCloser, error) {
	return withClose{r}, nil
}

func (nc noopCompressor) MediaType() string {
	return nc.mediaType
}

var GzipCompressor Compressor = gzipCompressor{}

type gzipCompressor struct{}

func (gz gzipCompressor) Compress(reader io.Reader) (io.ReadCloser, error) {
	pipeReader, pipeWriter := io.Pipe()

	gzw := gzip.NewWriter(pipeWriter)
	if err := gzw.SetConcurrency(256<<10, 2*runtime.NumCPU()); err != nil {
		return nil, errors.Wrapf(err, "set concurrency level to %v blocks", 2*runtime.NumCPU())
	}
	go func() {
		if _, err := io.Copy(gzw, reader); err != nil {
			// #nosec G104
			_ = pipeWriter.CloseWithError(errors.Wrap(err, "compressing layer"))
		}
		if err := gzw.Close(); err != nil {
			// #nosec G104
			_ = pipeWriter.CloseWithError(errors.Wrap(err, "close gzip writer"))
		}
		if err := pipeWriter.Close(); err != nil {
			// #nosec G104
			_ = pipeWriter.CloseWithError(errors.Wrap(err, "close pipe writer"))
		}
	}()

	return pipeReader, nil
}

func (gz gzipCompressor) MediaType() string {
	return ispec.MediaTypeImageLayerGzip
}
