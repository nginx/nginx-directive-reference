package output

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/nginxinc/ampex-apps/tools/reference-converter/parse"
)

type Directive struct {
	Name        string   `json:"name"`
	Default     string   `json:"default"`
	Contexts    []string `json:"contexts"`
	Syntax      []string `json:"syntax"`
	Description string   `json:"description"`
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
				Name:        directive.Name,
				Default:     directive.Default,
				Contexts:    directive.Contexts,
				Syntax:      directive.Syntax.ToMarkdown(),
				Description: directive.Prose.ToMarkdown(),
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
			res.Modules = append(res.Modules, toModule(m))
		}
	}

	return &res
}

func (r *Reference) Write(ctx context.Context, dst io.Writer) error {
	jsonData, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal reference struct: %w", err)
	}
	_, err = dst.Write(jsonData)
	if err != nil {
		return fmt.Errorf("unable to write data to io.Writer: %w", err)
	}
	return nil
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
