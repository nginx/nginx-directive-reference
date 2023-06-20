package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/nginxinc/ampex-apps/tools/reference-converter/tarball"
	"golang.org/x/exp/slog"
)

var (
	destFlag    = flag.String("dst", "reference.json", "where to write JSON output")
	sourceFlag  = flag.String("src", "http://hg.nginx.org/nginx.org/archive/tip.tar.gz", "where to get the XML sources")
	feedURLFlag = flag.String("feed-url", "http://hg.nginx.org/nginx.org/atom-log", "where to get the atom feed for XML changes")
	baseURLFlag = flag.String("base-url", "https://nginx.org/en/docs/", "base URL for rendering links inside the docs")
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

	// TODO: get the latest version from the atom feed
	// TODO: get the latest version from the destination
	// TODO: exit if the versions match
	_, err := tarball.Open(ctx, *sourceFlag)
	if err != nil {
		slog.ErrorCtx(ctx, "failed to read", slog.Any("error", err), slog.String("src", *sourceFlag))
		return
	}
	// TODO: find module XML files
	// TODO: parse into structs
	// TODO: marshall to json
}
