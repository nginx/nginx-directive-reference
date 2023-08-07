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

	return strings.TrimSuffix(content.String(), "\n "), nil
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
		return &example{}
	case "link":
		return &link{}
	case "list":
		return &list{}
	case "para":
		return &Paragraph{}
	case "header":
		return &header{}
	case "emphasis":
		return &emphasis{}
	case "http-status":
		return &httpStatus{}
	case "commercial_version":
		return &commercialVersion{}
	case "note":
		return &note{}
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

// <example> elements shows snippets of config, C code, etc.
type example struct {
	content string
}

func (e *example) ToMarkdown() string { return e.content }

// UnmarshalXML processes the elements in-order to generate correct content.
// Some <example>s contain <emphasis>, so needs to be parsed in-order.
func (e *example) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	content, err := unmarshalMarkdownXML(d, start)
	if err != nil {
		return err
	}
	content = strings.Trim(content, "\n")
	*e = example{content: fmt.Sprintf("```\n%s\n```", content)}
	return nil
}

// <header> elements are for HTTP headers
type header struct {
	Content string `xml:",chardata"`
}

func (h *header) ToMarkdown() string { return fmt.Sprintf(`"%s"`, h.Content) }

// <emphasis> elements are used to bold some of the code inside an <example>.
// There is no easy translation to markdown, so this bolding is dropped.
type emphasis struct {
	Content string `xml:",chardata"`
}

func (e *emphasis) ToMarkdown() string { return e.Content }

// <http-status> elements describe a HTTP status code.
type httpStatus struct {
	Code int    `xml:"code,attr"`
	Text string `xml:"text,attr"`
}

func (h *httpStatus) ToMarkdown() string {
	return fmt.Sprintf("%d (%s)", h.Code, h.Text)
}

// <commercial_version> elements are upsell links.
type commercialVersion struct {
	Content string `xml:",chardata"`
}

func (e *commercialVersion) ToMarkdown() string {
	return fmt.Sprintf("[%s](%s)", e.Content, current.upsellURL)
}

// <note> elements highlight some quirks or changes over time, rendered as
// blockquotes.
type note struct {
	content string
}

func (n *note) ToMarkdown() string { return n.content }

// UnmarshalXML processes the elements in-order to generate correct content.
// Some <note>s contain <literal>s, so needs to be parsed in-order.
func (n *note) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	content, err := unmarshalMarkdownXML(d, start)
	if err != nil {
		return err
	}
	// add a '>' prefix to each line
	var sb strings.Builder
	for _, line := range strings.Split(content, "\n") {
		if len(line) == 0 {
			continue
		}
		sb.WriteString(fmt.Sprintf("> %s\n", line))
	}
	*n = note{content: sb.String()}
	return nil
}
