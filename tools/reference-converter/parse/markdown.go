package parse

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/exp/slog"
)

// unmarshalMarkdownXML reads the XML in-order and converts it to markdown.
//
// Use it from xml.Unmashaler implementations for elements that need to convert
// their inner XML to markdown.
func unmarshalMarkdownXML(d *xml.Decoder, parent xml.StartElement) (string, error) {
	var content strings.Builder
LOOP:
	for {
		token, err := d.Token()
		if errors.Is(err, io.EOF) || token == nil {
			break
		}
		if err != nil {
			return "", err
		}

		switch t := token.(type) {
		case xml.CharData: // consume inline text
			content.WriteString(strings.Trim(string(t), "\t"))
		case xml.StartElement:
			md := chooseMarkdowner(t.Name)

			// consume child element
			if err := d.DecodeElement(md, &t); err != nil {
				return "", fmt.Errorf("failed to decode <%s>: %w", t.Name.Local, err)
			}
			content.WriteString(md.ToMarkdown())

		case xml.EndElement:
			if t.Name.Local != parent.Name.Local {
				return "", fmt.Errorf("unexpected </%s>, wanted </%s>", t.Name.Local, parent.Name.Local)
			}
			break LOOP
		case xml.Comment, xml.ProcInst, xml.Directive:
			// no processing needed
		}
	}

	return strings.Trim(content.String(), "\n "), nil
}

type markdowner interface {
	ToMarkdown() string
}

func chooseMarkdowner(name xml.Name) markdowner {
	switch name.Local {
	case "literal", "var", "command", "path", "c-def":
		return &code{}
	case "c-func":
		return &code{suffix: "()"}
	case "value":
		return &code{isEmphasized: true}
	case "example":
		return &fence{}
	case "link":
		return &link{}
	case "list":
		return &list{}
	case "para":
		return &Paragraph{}
	// TODO(AMPEX-72): handle other prose-y tags
	case "note", "http-status", "header", "commercial_version", "emphasis":
		return &unsupportedTag{}
	default:
		slog.Warn("unsupported tag", slog.String("name", name.Local))
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
	isEmphasized bool
	suffix       string // additional content to add, inside the code block
	Content      string `xml:",chardata"`
}

func (t *code) ToMarkdown() string {
	s := t.Content
	if t.suffix != "" {
		s += t.suffix
	}
	s = fmt.Sprintf("`%s`", s)

	if t.isEmphasized {
		return fmt.Sprintf("*%s*", s)
	}
	return s
}

type fence struct {
	Content string `xml:",chardata"`
}

func (t *fence) ToMarkdown() string {
	return fmt.Sprintf("```\n%s\n```", t.Content)
}
