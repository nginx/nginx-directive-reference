package parse

import (
	"path"
	"strings"

	"github.com/nginxinc/ampex-apps/tools/reference-converter/tarball"
)

// Reference is the collection of parsed docs for NGINX
type Reference struct {
	Modules     []*Module       // parsed and processed NGINX modules
	baseURL     string          // where the official docs live
	upsellURL   string          // where we link people when pushing the NGINX+
	pages       map[string]page // used to build links from directives
	currentPage page            // file currently being parsed, used to build links
	listDepth   int             // tracks nested <list> usage
}

func (r *Reference) parsePages(files []tarball.File) error {
	r.pages = make(map[string]page)

	for _, f := range files {
		if f.Contains("dtd/article.dtd") || f.Contains("dtd/module.dtd") {
			p := page{path: f.Name}
			if err := unmarshalXML(&p, f); err != nil {
				return err
			}
			r.pages[p.path] = p
		}
	}
	return nil
}

func (r *Reference) parseModules(files []tarball.File) error {
	currentMu.Lock()
	defer currentMu.Unlock()
	// set the context for UnmarshalXML implementations to read during parsing
	current = r
	defer func() { current = nil }()

	for _, f := range files {
		if !f.Contains("dtd/module.dtd") || strings.HasSuffix(f.Name, "_head.xml") {
			continue
		}
		res, err := r.parseModule(f)
		if err != nil {
			return err
		}
		r.Modules = append(r.Modules, res)
	}
	return nil
}

func (r *Reference) parseModule(f tarball.File) (*Module, error) {
	// set the context for UnmarshalXML implementations to read during parsing
	r.currentPage = r.pages[f.Name]
	defer func() { r.currentPage = page{} }()

	var res Module
	if err := unmarshalXML(&res, f); err != nil {
		return nil, err
	}
	return &res, nil
}

// getPage looks up another page, given relative path from the currentModule.
func (r *Reference) getPage(relpath string) (page, bool) {
	if relpath == "" {
		return page{}, false
	}
	p := path.Join(path.Dir(r.currentPage.path), relpath)
	page, ok := r.pages[p]
	return page, ok
}
