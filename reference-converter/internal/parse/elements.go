package parse

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// Syntax contain the markdown formatted syntax for the directive, very close to
// https://en.wikipedia.org/wiki/Backus%E2%80%93Naur_form
type Syntax struct {
	Content string
	IsBlock bool
}

func (s *Syntax) ToMarkdown() string { return s.Content }

type Syntaxes []Syntax

func (ss Syntaxes) ToMarkdown() []string {
	if len(ss) == 0 {
		return nil
	}

	ret := make([]string, 0, len(ss))
	for _, s := range ss {
		ret = append(ret, s.ToMarkdown())
	}
	return ret
}

func (ss Syntaxes) ToHTML() []string {
	if len(ss) == 0 {
		return nil
	}

	ret := make([]string, 0, len(ss))
	for _, s := range ss {
		md := []byte(s.ToMarkdown())
		ret = append(ret, string(mdToHTML(md)))
	}
	return ret
}

func (ss Syntaxes) IsBlock() bool {
	isBlock := false
	for _, s := range ss {
		if s.IsBlock {
			isBlock = true
			break
		}
	}
	return isBlock
}

var whitespace = regexp.MustCompile(`\s+`)

// UnmarshalXML processes the elements in-order to generate correct content,
// dropping incidental whitespace present in the source XML.
func (s *Syntax) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	content, err := unmarshalMarkdownXML(d, start)
	if err != nil {
		return err
	}
	content = whitespace.ReplaceAllString(content, " ")
	content = strings.Trim(content, " \n")
	attrs := newAttrs(start.Attr)
	isBlock := attrs["block"] == "yes"
	if attrs["block"] == "yes" {
		content = fmt.Sprintf("%s `%s`", content, "{...}")
	}
	*s = Syntax{
		Content: content,
		IsBlock: isBlock,
	}
	return nil
}

// Paragraphs contain the markdown converted content
type Paragraph struct {
	Content string
}

func (p *Paragraph) ToMarkdown() string { return p.Content }

// ToTrimmedMarkdown trims leading/trailing whitespace, useful for ignoring
// newlines and from the XML.
func (p *Paragraph) ToTrimmedMarkdown() string { return strings.Trim(p.ToMarkdown(), "\n\t ") }

func (p *Paragraph) ToIndentedMarkdown(isTagList bool) string {
	lines := strings.Split(strings.Trim(p.Content, "\n\t "), "\n")
	indentedLines := make([]string, 0, len(lines))
	for i, line := range lines {
		prefix := "    "
		if i == 0 && !isTagList {
			// The first line in an ordered and unordered list starts with one space,
			// remaining lines will have 4 spaces
			prefix = " "
		}
		indentedLines = append(indentedLines, prefix+line)
	}

	return strings.Join(indentedLines, "\n")
}

// UnmarshalXML processes the elements in-order to generate correct content
func (p *Paragraph) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	content, err := unmarshalMarkdownXML(d, start)
	if err != nil {
		return err
	}
	*p = Paragraph{Content: content}
	return nil
}

// Prose is a collection of paragraphs
type Prose []Paragraph

func (t Prose) ToMarkdown() string {
	paras := make([]string, 0, len(t))
	for _, p := range t {
		paras = append(paras, p.ToTrimmedMarkdown())
	}
	return strings.Join(paras, "\n\n")
}

func mdToHTML(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func (t Prose) ToHTML() string {
	md := []byte(t.ToMarkdown())
	return string(mdToHTML(md))
}

type Directive struct {
	Name     string   `xml:"name,attr"`
	Default  string   `xml:"default"`
	Contexts []string `xml:"context"`
	Syntax   Syntaxes `xml:"syntax"`
	Prose    Prose    `xml:"para"`
}

// Variable represents an NGINX variable defined by a module, e.g $binary_remote_addr.
type Variable struct {
	Name  string
	Prose Prose
}

// unmarshalVariablesCML extracts NGINX variables from the common pattern:
//
//	<section id="variables">
//	<list type="tag">
//	<tag-name><var>$VARIABLE_NAME</var></tag-name>
//	<tag-desc>$DOCUMENTATION</tag-desc>
//	<tag-name><var>$VARIABLE_NAME</var><value>$DYNAMIC_SUFFIX</value></tag-name>
//	<tag-desc>$DOCUMENTATION</tag-desc>
//	</list>
//	</section>
func unmarshalVariablesCML(d *xml.Decoder, start xml.StartElement) ([]Variable, error) {
	var v struct {
		ID         string `xml:"id,attr"`
		Paragraphs []struct {
			List struct {
				TagNames []struct {
					Name   string `xml:"var"`
					Suffix string `xml:"value"`
				} `xml:"tag-name"`
				TagDesc []Prose `xml:"tag-desc"`
			} `xml:"list"`
		} `xml:"para"`
	}
	if err := d.DecodeElement(&v, &start); err != nil {
		return nil, fmt.Errorf("failed to parse variables: %w", err)
	}
	var vs []Variable
	for _, para := range v.Paragraphs {
		if len(para.List.TagDesc) != len(para.List.TagNames) {
			return nil, fmt.Errorf(
				"invalid variables section, need to have the same number of names (%d) and descriptions (%d)",
				len(para.List.TagNames), len(para.List.TagDesc),
			)
		}
		for idx, tn := range para.List.TagNames {
			name := tn.Name
			if tn.Suffix != "" {
				name += strings.ToUpper(tn.Suffix)
			}
			vs = append(vs, Variable{
				Name:  name,
				Prose: para.List.TagDesc[idx],
			})
		}
	}

	return vs, nil
}

type Section struct {
	ID         string
	Directives []Directive
	Prose      Prose
	Variables  []Variable
}

// UnmarshalXML handles parsing sections with directives vs sections with variables.
func (s *Section) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	attrs := newAttrs(start.Attr)
	if attrs["id"] == "variables" {
		// parse as a list of variables
		vs, err := unmarshalVariablesCML(d, start)
		if err != nil {
			return fmt.Errorf("failed to unmarshall variables: %w", err)
		}
		*s = Section{
			ID:        "variables",
			Variables: vs,
		}
		return nil
	}

	// parse as a normal section
	var sec struct {
		ID         string      `xml:"id,attr"`
		Directives []Directive `xml:"directive"`
		Prose      Prose       `xml:"para"`
	}
	if err := d.DecodeElement(&sec, &start); err != nil {
		return err
	}

	*s = Section{
		ID:         sec.ID,
		Directives: sec.Directives,
		Prose:      sec.Prose,
	}
	return nil
}

type Module struct {
	XMLName  xml.Name  `xml:"module"`
	Name     string    `xml:"name,attr"`
	Link     string    `xml:"link,attr"`
	Lang     string    `xml:"lang,attr"`
	Sections []Section `xml:"section"`
}

// page represents <article>s or <module>s that are used with <link>
type page struct {
	Name string `xml:"name,attr"`
	Link string `xml:"link,attr"`
	path string // Path to the xml file
}
