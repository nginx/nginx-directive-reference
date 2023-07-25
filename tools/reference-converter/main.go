package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/nginxinc/ampex-apps/tools/reference-converter/atom"
	"github.com/nginxinc/ampex-apps/tools/reference-converter/output"
	"github.com/nginxinc/ampex-apps/tools/reference-converter/parse"
	"github.com/nginxinc/ampex-apps/tools/reference-converter/tarball"
	"golang.org/x/exp/slog"
)

var (
	destFlag      = flag.String("dst", "reference.json", "where to write JSON output")
	sourceFlag    = flag.String("src", "http://hg.nginx.org/nginx.org/archive/tip.tar.gz", "where to get the XML sources")
	feedURLFlag   = flag.String("feed-url", "http://hg.nginx.org/nginx.org/atom-log", "where to get the atom feed for XML changes")
	baseURLFlag   = flag.String("base-url", "https://nginx.org/en/docs/", "base URL for rendering links inside the docs")
	upsellURLFlag = flag.String("upsell-url", "https://nginx.com/products/", "URL for linking people to NGINX+")
)

func main() {
	opts := slog.HandlerOptions{Level: slog.LevelDebug}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &opts)))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGHUP, syscall.SIGABRT, syscall.SIGINT)
	defer stop()

	flag.Parse()

	slog.InfoCtx(ctx, "started", slog.Group("opts",
		slog.String("dst", *destFlag),
		slog.String("src", *sourceFlag),
		slog.String("feed-url", *feedURLFlag),
		slog.String("base-url", *baseURLFlag)))
	defer slog.InfoCtx(ctx, "finished")

	// TODO: get the latest version from the atom feed (atom.go)
	v1, err := atom.GetVersion(ctx, *feedURLFlag)
	if err != nil {
		slog.ErrorCtx(ctx, "failed to get the version", slog.Any("error", err), slog.String("src", *feedURLFlag))
	}
	slog.InfoCtx(ctx, "Comparing Versions", slog.String("atom", v1))
	// TODO: get the latest version from the destination

	//v2 :=
	// TODO: get the latest version from the destination
	// TODO: exit if the versions match
	// unpack the tarball
	files, err := tarball.Open(ctx, *sourceFlag)
	if err != nil {
		slog.ErrorCtx(ctx, "failed to read", slog.Any("error", err), slog.String("src", *sourceFlag))
		return
	}

	// reading files, converts XML to markdown
	r, err := parse.Parse(files, *baseURLFlag, *upsellURLFlag)
	if err != nil {
		slog.ErrorCtx(ctx, "failed to parse", slog.Any("error", err))
		return
	}
	slog.InfoCtx(ctx, "parsed into modules", slog.Int("n", len(r.Modules)))

	// convert XML types to JSON types
	ref := output.New(v1, r.Modules)

	dst, err := os.Create(*destFlag)
	if err != nil {
		slog.ErrorCtx(ctx, "failed to open dst", slog.Any("error", err))
		return
	}
	defer dst.Close()
	if err := ref.Write(ctx, dst); err != nil {
		slog.ErrorCtx(ctx, "failed to save", slog.Any("error", err))
		return
	}
}
