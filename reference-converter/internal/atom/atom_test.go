package atom_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/nginxinc/nginx-directive-reference/reference-converter/internal/atom"
	"github.com/stretchr/testify/require"
)

func readTestDataFile(t *testing.T) string {
	t.Helper()
	body, err := os.ReadFile("./testdata/test.xml")
	require.NoError(t, err)
	return string(body)
}
func TestGetVersion(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		statusCode   int
		responseBody string
		wantError    bool
		want         string
	}{
		"BadStatus": {
			statusCode:   http.StatusBadRequest,
			responseBody: "",
			wantError:    true,
		},
		"OKStatusWithRealisticXML": {
			statusCode:   http.StatusOK,
			responseBody: readTestDataFile(t),
			wantError:    false,
			want:         "http://hg.nginx.org/nginx.org/rev/c80a7cb452e8",
		},
		"OKStatusWithNoEntries": {
			statusCode:   http.StatusOK,
			responseBody: "<feed></feed>",
			wantError:    true,
		},
		"OKStatusWithInvalidXML": {
			statusCode:   http.StatusOK,
			responseBody: "I'm not XML",
			wantError:    true,
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {

			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(testCase.statusCode)
				_, err := w.Write([]byte(testCase.responseBody))
				if err != nil {
					t.Errorf("Response body failed to write: %s", err)
				}
			}))
			defer srv.Close()

			got, err := atom.GetVersion(context.Background(), srv.URL, atom.WithHttpClient(*srv.Client()))

			if testCase.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, testCase.want, got)
			}

		})
	}
}
