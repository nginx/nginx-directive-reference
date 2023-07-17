package parse

import (
	"encoding/xml"
	"fmt"
)

// link converts <link> elements into markdown links
type link struct {
	content string
}

func (l *link) ToMarkdown() string { return l.content }

// UnmarshalXML processes the elements in-order to generate correct content
func (l *link) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// handle inner content like <link>title</link> or
	// <link><literal>title</literal></link>
	title, err := unmarshalMarkdownXML(d, start)
	if err != nil {
		return err
	}

	// manually work with attrs, unmarshalMarkdownXML consumes the whole element
	attrs := newAttrs(start.Attr)
	p, hasPage := current.getPage(attrs["doc"])

	// linking to a directive, e.g. <link id="anchor" />
	if title == "" && attrs["id"] != "" {
		title = fmt.Sprintf("`%s`", attrs["id"])
	}
	// find the name of the other page, e.g. <link doc="page.xml" />
	if title == "" && hasPage {
		title = p.Name
	}

	href := attrs["url"]
	if href == "" {
		// default to self-link, e.g. <link id="anchor">
		href = current.currentPage.Link
		if hasPage {
			// linking to another page, e.g. <link doc="page.xml">
			href = p.Link
		}

		if attrs["id"] != "" {
			href += "#" + attrs["id"]
		}
		href = current.baseURL + href
	}

	*l = link{
		content: fmt.Sprintf("[%s](%s)", title, href),
	}
	return nil
}
