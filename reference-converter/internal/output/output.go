package output

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/nginxinc/nginx-directive-reference/reference-converter/internal/parse"
)

type Directive struct {
	Name            string   `json:"name"`
	Default         string   `json:"default"`
	Contexts        []string `json:"contexts"`
	SyntaxMd        []string `json:"syntax_md"`
	SyntaxHtml      []string `json:"syntax_html"`
	IsBlock         bool     `json:"isBlock"`
	DescriptionMd   string   `json:"description_md"`
	DescriptionHtml string   `json:"description_html"`
}

type Module struct {
	Id         string      `json:"id"`
	Name       string      `json:"name"`
	Directives []Directive `json:"directives"`
}

func toModule(m *parse.Module) Module {
	module := Module{
		Name: strings.TrimLeft(m.Name, "Module "),
		Id:   m.Link,
	}
	for _, section := range m.Sections {
		for _, directive := range section.Directives {
			module.Directives = append(module.Directives, Directive{
				Name:            directive.Name,
				Default:         directive.Default,
				Contexts:        directive.Contexts,
				SyntaxMd:        directive.Syntax.ToMarkdown(),
				SyntaxHtml:      directive.Syntax.ToHTML(),
				IsBlock:         directive.Syntax.IsBlock(),
				DescriptionMd:   directive.Prose.ToMarkdown(),
				DescriptionHtml: directive.Prose.ToHTML(),
			})
		}
	}
	return module
}

type Reference struct {
	Modules []Module `json:"modules"`
	Version string   `json:"version"`
}

func New(version string, modules []*parse.Module) *Reference {
	res := Reference{
		Modules: make([]Module, 0, len(modules)),
		Version: version,
	}

	for _, m := range modules {
		if m.Lang == "en" {
			mod := toModule(m)
			// filter modules with zero directives
			if len(mod.Directives) > 0 {
				res.Modules = append(res.Modules, mod)
			}
		}
	}

	return &res
}

func (r *Reference) Write(ctx context.Context, dst io.Writer) error {
	enc := json.NewEncoder(dst)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func GetVersion(ctx context.Context, r io.Reader) (string, error) {
	var reference Reference
	refData, err := io.ReadAll(r)
	if err != nil {
		return "", fmt.Errorf("unable to read the reference data: %w", err)
	}
	err = json.Unmarshal(refData, &reference)
	if err != nil {
		return "", fmt.Errorf("unable to unmarshal json data: %w", err)
	}
	return reference.Version, nil
}
