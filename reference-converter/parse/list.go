package parse

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// unorderedList handles <list type="bullet">.
type unorderedList struct {
	Items []Paragraph `xml:"listitem"`
}

func (t *unorderedList) ToMarkdown() string {
	var sb strings.Builder
	for _, item := range t.Items {
		sb.WriteString(fmt.Sprintf("-%s\n", item.ToIndentedMarkdown(false)))
	}
	return sb.String()
}

// orderedList handles <list type="enum">.
type orderedList struct {
	Items []Paragraph `xml:"listitem"`
}

func (t *orderedList) ToMarkdown() string {
	var sb strings.Builder
	for i, item := range t.Items {
		sb.WriteString(fmt.Sprintf("%d.%s\n", i+1, item.ToIndentedMarkdown(false)))
	}
	return sb.String()
}

// taglist handles <list type="tag">. These are rendered as <dl>s in the
// official docs, which don't have a direct mapping in pure markdown. Simulates
// it using unordered lists and indentation.
type taglist struct {
	TagNames []Paragraph `xml:"tag-name"`
	TagDesc  []Paragraph `xml:"tag-desc"`
}

func (t *taglist) ToMarkdown() string {
	if len(t.TagNames) != len(t.TagDesc) {
		panic(fmt.Sprintf("tag lists must have same number of names (%d) as descs (%d)", len(t.TagNames), len(t.TagDesc)))
	}
	var sb strings.Builder

	for i := range t.TagNames {
		name := t.TagNames[i]
		desc := t.TagDesc[i]
		sb.WriteString(fmt.Sprintf(
			"- %s\n\n%s\n", name.ToTrimmedMarkdown(), desc.ToIndentedMarkdown(true)))

	}
	return sb.String()
}

// list parses a variety of `<list>` types to markdown.
type list struct {
	content string
}

func (t *list) ToMarkdown() string { return t.content }

// UnmarshalXML processes the elements in-order to generate correct content
func (l *list) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	attrs := newAttrs(start.Attr)
	listType := attrs["type"]
	var sub markdowner
	switch listType {
	case "bullet":
		sub = &unorderedList{}
	case "tag":
		sub = &taglist{}
	case "enum":
		sub = &orderedList{}
	default:
		return fmt.Errorf("unknown list type '%s'", listType)
	}

	if err := d.DecodeElement(sub, &start); err != nil {
		return fmt.Errorf("failed to parse %s list: %w", listType, err)
	}

	*l = list{
		content: sub.ToMarkdown(),
	}
	return nil
}
