package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nginxinc/nginx-directive-reference/reference-converter/internal/atom"
	"github.com/nginxinc/nginx-directive-reference/reference-converter/internal/output"
	"github.com/nginxinc/nginx-directive-reference/reference-converter/internal/parse"
	"github.com/nginxinc/nginx-directive-reference/reference-converter/internal/tarball"
)

var (
	destFlag      = flag.String("dst", "reference.json", "where to write JSON output")
	sourceFlag    = flag.String("src", "https://github.com/nginx/nginx.org/archive/refs/heads/main.tar.gz", "where to get the XML sources")
	feedURLFlag   = flag.String("feed-url", "https://github.com/nginx/nginx.org/commits/main.atom", "where to get the atom feed for XML changes")
	baseURLFlag   = flag.String("base-url", "https://nginx.org", "base URL for rendering links inside the docs")
	upsellURLFlag = flag.String("upsell-url", "https://nginx.com/products/", "URL for linking people to NGINX+")
)

func main() {
	err := runConverter()
	if err != nil {
		os.Exit(1)
	}
}

func runConverter() error {
	opts := slog.HandlerOptions{Level: slog.LevelDebug}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &opts)))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGHUP, syscall.SIGABRT, syscall.SIGINT)
	defer stop()

	flag.Parse()

	slog.InfoContext(ctx, "started", slog.Group("opts",
		slog.String("dst", *destFlag),
		slog.String("src", *sourceFlag),
		slog.String("feed-url", *feedURLFlag),
		slog.String("base-url", *baseURLFlag)))
	defer slog.InfoContext(ctx, "finished")

	v1, err := atom.GetVersion(ctx, *feedURLFlag)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get the version", slog.Any("error", err), slog.String("src", *feedURLFlag))
	}
	slog.InfoContext(ctx, "Comparing Versions", slog.String("atom", v1))

	// unpack the tarball
	files, err := tarball.Open(ctx, *sourceFlag)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read", slog.Any("error", err), slog.String("src", *sourceFlag))
		return err
	}

	// reading files, converts XML to markdown
	r, err := parse.Parse(files, *baseURLFlag, *upsellURLFlag)
	if err != nil {
		slog.ErrorContext(ctx, "failed to parse", slog.Any("error", err))
		return err
	}
	slog.InfoContext(ctx, "parsed into modules", slog.Int("n", len(r.Modules)))

	// convert XML types to JSON types
	ref := output.New(v1, r.Modules)

	dst, err := os.Create(*destFlag)
	if err != nil {
		slog.ErrorContext(ctx, "failed to open dst", slog.Any("error", err))
		return err
	}
	defer dst.Close()
	if err := ref.Write(ctx, dst); err != nil {
		slog.ErrorContext(ctx, "failed to save", slog.Any("error", err))
		return err
	}
	return nil
}
