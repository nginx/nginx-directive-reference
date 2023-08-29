package parse_test

import (
	"testing"

	"github.com/nginxinc/nginx-directive-reference/reference-converter/parse"
	"github.com/nginxinc/nginx-directive-reference/reference-converter/tarball"
	"github.com/stretchr/testify/require"
)

const baseURL = "http://example.org"
const upsellURL = "http://example.com"

// TODO: move into TestMarkdown
func TestLink_ToMarkdown(t *testing.T) {
	t.Parallel()
	// simulate other pages in the tarball
	pages := []tarball.File{
		testArticleFile("/xml/en/debugging.xml", "I'm another file"),
		testArticleFile("/xml/above.xml", "I'm in the parent dir"),
		testModuleFile(t, withPath("/xml/en/child/below.xml")),
	}
	testcases := map[string]struct {
		XML  string
		want string
	}{
		"same page link": {
			XML:  `<link id="accept_mutex"/>`,
			want: "[`accept_mutex`](http://example.org/en/test.html#accept_mutex)",
		},
		"sibling page uses it's data": {
			XML:  `<link doc="debugging.xml"/>`,
			want: `[I'm another file](http://example.org/en/debugging.html)`,
		},
		"child page with inner text": {
			XML:  `<link doc="child/below.xml">ngx_http_perl_module</link>`,
			want: `[ngx_http_perl_module](http://example.org/en/child/below.html)`,
		},
		"child page with anchor": {
			XML:  `<link doc="child/below.xml" id="reuseport"/>`,
			want: "[`reuseport`](http://example.org/en/child/below.html#reuseport)",
		},
		"parent page with anchor": {
			XML:  `<link doc="../above.xml" id="epoll"/>`,
			want: "[`epoll`](http://example.org/above.html#epoll)",
		},
		"sibling page with inner xml": {
			XML:  `<link doc="debugging.xml"><literal>--foo</literal></link>`,
			want: "[`--foo`](http://example.org/en/debugging.html)",
		},
		"url": {
			XML:  `<link url="http://nginx.org">hello</link>`,
			want: "[hello](http://nginx.org)",
		},
		"title with new line": {
			XML:  `<link url="http://nginx.org">hello&#xA;world</link>`,
			want: "[hello world](http://nginx.org)",
		},
	}
	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			files := append(pages, testModuleFile(t, withContent(tc.XML)))

			ref, err := parse.Parse(files, baseURL, upsellURL)
			require.NoError(t, err)
			require.NotNil(t, ref)

			require.Equal(t, 2, len(ref.Modules))
			got := ref.Modules[1].Sections[0].Directives[0].Prose.ToMarkdown()

			require.Equal(t, tc.want, got, "failed on `%s`", tc.XML)
		})
	}
}
