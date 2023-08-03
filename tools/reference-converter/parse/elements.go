package parse

import (
	"encoding/xml"
	"regexp"
	"strings"
)

// Syntax contain the markdown formatted syntax for the directive, very close to
// https://en.wikipedia.org/wiki/Backus%E2%80%93Naur_form
type Syntax struct {
	Content string
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
	*s = Syntax{Content: content}
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

type Directive struct {
	Name     string   `xml:"name,attr"`
	Default  string   `xml:"default"`
	Contexts []string `xml:"context"`
	Syntax   Syntaxes `xml:"syntax"`
	Prose    Prose    `xml:"para"`
}

type Section struct {
	ID         string      `xml:"id,attr"`
	Directives []Directive `xml:"directive"`
	Prose      Prose       `xml:"para"`
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
