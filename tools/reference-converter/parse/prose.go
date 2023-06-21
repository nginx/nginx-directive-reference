package parse

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Paragraphs contain the markdown converted content
type Paragraph struct {
	Content string
}

func (p *Paragraph) ToMarkdown() string { return p.Content }

// UnmarshalXML processes the elements in-order to generate correct content
func (p *Paragraph) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var content strings.Builder
LOOP:
	for {
		token, err := d.Token()
		if errors.Is(err, io.EOF) || token == nil {
			break
		}
		if err != nil {
			return err
		}

		switch t := token.(type) {
		case xml.CharData: // consume inline text
			content.WriteString(string(t))
		case xml.StartElement:
			md := chooseMarkdowner(t.Name)

			// consume child element
			if err := d.DecodeElement(md, &t); err != nil {
				return fmt.Errorf("failed to decode <%s>: %w", t.Name.Local, err)
			}
			content.WriteString(md.ToMarkdown())

		case xml.EndElement:
			if t.Name.Local != start.Name.Local {
				return fmt.Errorf("unexpected </%s>, wanted </%s>", t.Name.Local, start.Name.Local)
			}
			break LOOP
		case xml.Comment, xml.ProcInst, xml.Directive:
			// no processing needed
		}
	}
	*p = Paragraph{Content: strings.Trim(content.String(), "\n ")}
	return nil
}

// Prose is a collection of paragraphs
type Prose []Paragraph

func (t Prose) ToMarkdown() string {
	paras := make([]string, 0, len(t))
	for _, p := range t {
		paras = append(paras, p.ToMarkdown())
	}
	return strings.Join(paras, "\n\n")
}

type markdowner interface {
	ToMarkdown() string
}

func chooseMarkdowner(name xml.Name) markdowner {
	switch name.Local {
	// TODO(AMPEX-72): handle other prose-y tags
	case "literal":
		return &code{}
	case "example":
		return &fence{}
	default:
		return &unsupportedTag{}
	}
}

type unsupportedTag struct {
	XMLName  xml.Name
	Contents string `xml:",innerxml"`
}

func (t *unsupportedTag) ToMarkdown() string {
	return fmt.Sprintf("`TODO: handle <%s>`", t.XMLName.Local)
}

type code struct {
	Content string `xml:",chardata"`
}

func (t *code) ToMarkdown() string {
	return fmt.Sprintf("`%s`", t.Content)
}

type fence struct {
	Content string `xml:",chardata"`
}

func (t *fence) ToMarkdown() string {
	return fmt.Sprintf("```\n%s\n```", t.Content)
}
