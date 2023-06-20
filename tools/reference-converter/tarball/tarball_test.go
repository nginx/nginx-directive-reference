package tarball_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nginxinc/ampex-apps/tools/reference-converter/tarball"
	"github.com/stretchr/testify/require"
)

func TestOpen_File(t *testing.T) {
	t.Parallel()
	files, err := tarball.Open(context.Background(), "./testdata/test.tar.gz")

	require.NoError(t, err)
	require.ElementsMatch(t, files, []tarball.File{
		{Name: "foo.txt", Contents: []byte("foo\n")},
		{Name: "bar.txt", Contents: []byte("bar\n")},
	})
}

func TestOpen_Url(t *testing.T) {
	t.Parallel()
	// serve the tarball via HTTP
	fs := http.FileServer(http.Dir("./testdata"))
	srv := httptest.NewServer(fs)
	defer srv.Close()
	url := srv.URL + "/test.tar.gz"

	files, err := tarball.Open(context.Background(), url, tarball.WithHttpClient(*srv.Client()))

	require.NoError(t, err)
	require.ElementsMatch(t, files, []tarball.File{
		{Name: "foo.txt", Contents: []byte("foo\n")},
		{Name: "bar.txt", Contents: []byte("bar\n")},
	})
}
