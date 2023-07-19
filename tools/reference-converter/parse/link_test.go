package parse_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nginxinc/ampex-apps/tools/reference-converter/parse"
	"github.com/nginxinc/ampex-apps/tools/reference-converter/tarball"
	"github.com/stretchr/testify/require"
)

const baseURL = "http://example.com"

func TestLink_ToMarkdown(t *testing.T) {
	t.Parallel()
	testModulePath := "/xml/en/test.xml"
	// simulate other pages in the tarball
	pages := []tarball.File{
		testArticleFile("/xml/en/debugging.xml", "I'm another file"),
		testArticleFile("/xml/above.xml", "I'm in the parent dir"),
		testArticleFile("/xml/en/child/below.xml", "I'm in the child dir"),
	}
	testcases := map[string]struct {
		XML  string
		want string
	}{
		"same page link": {
			XML:  `<link id="accept_mutex"/>`,
			want: "[`accept_mutex`](http://example.com/en/test.html#accept_mutex)",
		},
		"sibling page uses it's data": {
			XML:  `<link doc="debugging.xml"/>`,
			want: `[I'm another file](http://example.com/en/debugging.html)`,
		},
		"child page with inner text": {
			XML:  `<link doc="child/below.xml">ngx_http_perl_module</link>`,
			want: `[ngx_http_perl_module](http://example.com/en/child/below.html)`,
		},
		"child page with anchor": {
			XML:  `<link doc="child/below.xml" id="reuseport"/>`,
			want: "[`reuseport`](http://example.com/en/child/below.html#reuseport)",
		},
		"parent page with anchor": {
			XML:  `<link doc="../above.xml" id="epoll"/>`,
			want: "[`epoll`](http://example.com/above.html#epoll)",
		},
		"sibling page with inner xml": {
			XML:  `<link doc="debugging.xml"><literal>--foo</literal></link>`,
			want: "[`--foo`](http://example.com/en/debugging.html)",
		},
		"url": {
			XML:  `<link url="http://nginx.org">hello</link>`,
			want: "[hello](http://nginx.org)",
		},
	}
	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			files := append(pages, testModuleFile(testModulePath, tc.XML))

			ref, err := parse.Parse(files, baseURL)
			require.NoError(t, err)
			require.NotNil(t, ref)

			require.Equal(t, 1, len(ref.Modules))
			got := ref.Modules[0].Sections[0].Prose.ToMarkdown()

			require.Equal(t, tc.want, got, "failed on `%s`", tc.XML)
		})
	}
}

func testModuleFile(path, paraXML string) tarball.File {
	tmpl := `<!DOCTYPE article SYSTEM "../dtd/module.dtd">
	<module link="%s">
		<section>
			<para>%s</para>
		</section>
	</module>`

	link := strings.TrimPrefix(strings.Replace(path, ".xml", ".html", 1), "/xml")

	return tarball.File{
		Name:     path,
		Contents: []byte(fmt.Sprintf(tmpl, link, paraXML)),
	}
}

func testArticleFile(path, name string) tarball.File {
	tmpl := `<!DOCTYPE article SYSTEM "../dtd/article.dtd">
	<article link="%s" name="%s" />`

	link := strings.TrimPrefix(strings.Replace(path, ".xml", ".html", 1), "/xml")

	return tarball.File{
		Name:     path,
		Contents: []byte(fmt.Sprintf(tmpl, link, name)),
	}
}
