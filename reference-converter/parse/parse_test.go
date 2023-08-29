package parse_test

import (
	"encoding/xml"
	"testing"

	"github.com/nginxinc/nginx-directive-reference/reference-converter/parse"
	"github.com/nginxinc/nginx-directive-reference/reference-converter/tarball"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Parallel()
	testcases := map[string]struct {
		filename  string
		wantError bool
		want      *parse.Module
	}{
		"module file": {
			filename: "module.xml",
			want: &parse.Module{
				XMLName: xml.Name{Local: "module"},
				Name:    "Module ngx_FAKE_TEST_module",
				Link:    "/en/docs/FAKE/ngx_FAKE_TEST_module.html",
				Lang:    "en",
				Sections: []parse.Section{
					{
						ID: "directives",
						Directives: []parse.Directive{
							{
								Name:     "testing",
								Default:  "on",
								Contexts: []string{"http", "server", "location"},
								Syntax: parse.Syntaxes{{
									Content: "`on` | `off`",
								}},
								Prose: parse.Prose{
									{Content: "\nFree form test.\n"},
									{Content: "\nCan have more than one, with some html—ish entities and `verbatim` text.\n"},
								},
							},
						},
					},
				},
			},
		},
		"incomplete module": {
			filename: "incomplete.xml",
			want: &parse.Module{
				XMLName: xml.Name{Local: "module"},
				Name:    "Module ngx_FAKE_TEST_module",
				Link:    "/en/docs/FAKE/ngx_FAKE_TEST_module.html",
				Lang:    "en",
				Sections: []parse.Section{
					{
						ID: "directives",
						Directives: []parse.Directive{
							{Name: "who_needs_closing_tags"},
						},
					},
				},
			},
		},
		"not a module": {
			filename: "not-a-module.xml",
		},
	}
	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			f := readTestFile(t, tc.filename)
			got, err := parse.Parse([]tarball.File{f}, baseURL, upsellURL)

			if tc.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if tc.want == nil {
				require.Equal(t, 0, len(got.Modules))
			} else {
				require.Equal(t, 1, len(got.Modules))
				require.Equal(t, *tc.want, *got.Modules[0])
			}
		})
	}
}
