package parse

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log/slog"

	"github.com/nginxinc/nginx-directive-reference/reference-converter/tarball"
)

// attrMap is a key/value version of []xml.Attr, ignoring XML namespaces. Create
// with newAttrs.
type attrMap map[string]string

func newAttrs(attrs []xml.Attr) attrMap {
	ret := make(attrMap, len(attrs))
	for _, a := range attrs {
		ret[a.Name.Local] = a.Value
	}
	return ret
}

// unmarshalXML works like xml.Unmarshal, but configured to handle the HTML
// entities we see in NGINX docs and other quirks in the XML.
func unmarshalXML(v any, f tarball.File) error {
	// some files are missing a closing tag
	if f.Contains("<module") && !f.Contains("</module>") {
		slog.Warn("fixed missing </module>", slog.String("file", f.Name))
		f.Contents = append(f.Contents, []byte("</module>")...)
	}

	decoder := xml.NewDecoder(bytes.NewReader(f.Contents))
	decoder.Entity = map[string]string{
		"nbsp":  " ",
		"mdash": "—",
		"ldquo": "“",
		"rdquo": "”",
		"lsquo": "‘",
		"rsquo": "’",
		"times": "×",
	}

	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("unable to parse %s: %w", f.Name, err)
	}
	return nil
}
