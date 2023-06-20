package tarball

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/exp/slog"
)

type File struct {
	Name     string
	Contents []byte
}

type config struct {
	Client http.Client
}
type Option = func(*config)

// WithHttpClient uses the provided client instead of the default for
// downloading tarballs.
func WithHttpClient(c http.Client) Option {
	return func(o *config) { o.Client = c }
}

// Open reads a tarball from the given path or url, and returns a slice of all
// the files inside.
func Open(ctx context.Context, pathOrURL string, opts ...Option) ([]File, error) {
	if !strings.HasSuffix(pathOrURL, ".tar.gz") {
		return nil, errors.New("invalid source, must be a tar.gz")
	}

	if _, err := os.Stat(pathOrURL); err == nil {
		return openFile(ctx, pathOrURL)
	}

	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}
	return openURL(ctx, pathOrURL, cfg.Client)
}

func openURL(ctx context.Context, url string, client http.Client) ([]File, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unabled to download %s: %w", url, err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unabled to download %s: %s", url, res.Status)
	}
	return open(ctx, res.Body, slog.With(slog.String("url", url)))
}

func openFile(ctx context.Context, path string) ([]File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return open(ctx, f, slog.With(slog.String("path", path)))
}

func open(ctx context.Context, raw io.ReadCloser, log *slog.Logger) ([]File, error) {
	log.DebugCtx(ctx, "opening tarball")
	defer raw.Close()
	gz, err := gzip.NewReader(raw)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)

	log.DebugCtx(ctx, "reading tarball")
	defer log.DebugCtx(ctx, "done reading")
	var res []File
	for {
		// stop if the context is canceled
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		header, err := tr.Next()

		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read next tarball entry: %w", err)
		}

		// we only care about regular files
		if header.Typeflag != tar.TypeReg {
			continue
		}

		buf, err := io.ReadAll(tr)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s contents: %w", header.Name, err)
		}
		res = append(res, File{Name: header.Name, Contents: buf})

	}
	log.DebugCtx(ctx, "read tarball", slog.Int("numFiles", len(res)))
	return res, nil
}
