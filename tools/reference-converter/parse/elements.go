package parse

import "encoding/xml"

// TODO: this might need in-order parsing, it's almost BNF
type Syntax struct {
	Values   []string `xml:"value"`
	Literals []string `xml:"literal"`
}

type Directive struct {
	Name     string   `xml:"name,attr"`
	Default  string   `xml:"default"`
	Contexts []string `xml:"context"`
	Syntax   Syntax   `xml:"syntax"`
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
