package parse

import (
	"fmt"
	"sync"

	"github.com/nginxinc/nginx-directive-reference/reference-converter/internal/tarball"
)

// current refers to the Reference currently being parsed.
//
// xml.Decoder.Decode gives no way to pass contextual information. Deep in the
// XML tree we need to know things like the current module being parsed or an
// attribute from another XML file via relative path.
var current *Reference
var currentMu sync.Mutex // protects current

// Parse reads and parses all the XML files, converting prose to markdown on the
// way to respect the ordering of XML elements.
func Parse(xmlFiles []tarball.File, baseURL, upsellURL string) (*Reference, error) {
	ref := &Reference{baseURL: baseURL, upsellURL: upsellURL}

	// read all the files so we can build links
	if err := ref.parsePages(xmlFiles); err != nil {
		return nil, fmt.Errorf("failed to parse articles: %w", err)
	}

	// read all modules
	if err := ref.parseModules(xmlFiles); err != nil {
		return nil, fmt.Errorf("failed to parse modules: %w", err)
	}

	return ref, nil
}
