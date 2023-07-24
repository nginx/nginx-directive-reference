package parse_test

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/nginxinc/ampex-apps/tools/reference-converter/tarball"
	"github.com/stretchr/testify/require"
)

func lines(l ...string) string { return strings.Join(l, "\n") }

type xmlConfig struct {
	AddPara         bool
	Path, Link, XML string
}

type xmlOption func(*xmlConfig)

// withPara controls whether the test XML content is wrapped by a <para>. Defaults to true.
func withPara(wrap bool) xmlOption {
	return func(xc *xmlConfig) {
		xc.AddPara = wrap
	}
}

// withPath sets the file path for the module XML and the <module link> attribute.
func withPath(path string) xmlOption {
	return func(xc *xmlConfig) {
		xc.Path = path
		xc.Link = strings.TrimPrefix(strings.Replace(path, ".xml", ".html", 1), "/xml")
	}
}

// moduleTemplate creates realistic nginx doc XML
var moduleTemplate = template.Must(template.New("mod").Parse(`
<!DOCTYPE module SYSTEM "../dtd/module.dtd">
<module link="{{.Link}}">
<section>
{{ if .AddPara -}} <para> {{- end }}
{{.XML}}
{{ if .AddPara -}} </para> {{- end }}
</section>
</module>
`))

func testModuleFile(t *testing.T, XML string, opts ...xmlOption) tarball.File {
	t.Helper()
	cfg := xmlConfig{
		AddPara: true,
		Path:    "/xml/en/test.xml",
		Link:    "/en/test.html",
		XML:     XML,
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	var buf bytes.Buffer
	require.NoError(t, moduleTemplate.Execute(&buf, cfg))

	return tarball.File{
		Name:     cfg.Path,
		Contents: buf.Bytes(),
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

func readTestFile(t *testing.T, filename string) tarball.File {
	t.Helper()
	content, err := os.ReadFile("./testdata/" + filename)
	require.NoError(t, err)
	return tarball.File{
		Name:     filename,
		Contents: content,
	}
}
