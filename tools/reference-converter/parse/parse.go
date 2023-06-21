package parse

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/nginxinc/ampex-apps/tools/reference-converter/tarball"
	"golang.org/x/exp/slog"
)

// IsModule checks if this looks like it contains module definitions.
func IsModule(f tarball.File) bool {
	return strings.HasSuffix(f.Name, ".xml") && strings.Contains(string(f.Contents), "dtd/module.dtd")
}

// NewModule parses the module XML, converting to markdown while reading.
func NewModule(f tarball.File) (*Module, error) {

	contents := f.Contents
	// some files are missing a closing tag
	if !strings.Contains(string(contents), "</module>") {
		slog.Warn("fixed missing </module>", slog.String("file", f.Name))
		contents = append(f.Contents, []byte("</module>")...)
	}

	var res Module
	// needs a custom decoder to handle HTML entities
	decoder := xml.NewDecoder(bytes.NewReader(contents))
	decoder.Entity = map[string]string{
		"nbsp":  " ",
		"mdash": "—",
	}

	if err := decoder.Decode(&res); err != nil {
		return nil, fmt.Errorf("unable to parse %s: %w", f.Name, err)
	}
	return &res, nil
}
