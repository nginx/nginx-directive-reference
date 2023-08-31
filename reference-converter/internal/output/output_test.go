package output_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/nginxinc/nginx-directive-reference/reference-converter/internal/output"
	"github.com/nginxinc/nginx-directive-reference/reference-converter/internal/parse"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	modules := []*parse.Module{
		{Name: "Module 1", Lang: "en"},
		{Name: "Module 2", Lang: "en", Sections: []parse.Section{
			{Directives: []parse.Directive{
				{
					Name:     "directive 2",
					Default:  "default 2",
					Contexts: []string{"context 1", "context 2"},
					Syntax: []parse.Syntax{{
						Content: "syntax 1",
					}, {
						Content: "syntax 2",
					}},
					Prose: parse.Prose{
						{Content: "Test"},
					},
				},
			}},
		}},
	}
	got := output.New("1.0", modules)
	want := &output.Reference{
		Modules: []output.Module{
			{
				Name: "2",
				Directives: []output.Directive{
					{
						Name:            "directive 2",
						Default:         "default 2",
						Contexts:        []string{"context 1", "context 2"},
						SyntaxMd:        []string{"syntax 1", "syntax 2"},
						SyntaxHtml:      []string{"<p>syntax 1</p>\n", "<p>syntax 2</p>\n"},
						DescriptionMd:   "Test",
						DescriptionHtml: "<p>Test</p>\n",
					},
				},
			},
		},
		Version: "1.0",
	}
	require.Equal(t, want, got)

}
func TestWrite(t *testing.T) {

	want := &output.Reference{
		Modules: []output.Module{
			{
				Name: "Module 1",
				Directives: []output.Directive{
					{
						Name:            "Directive 1",
						Default:         "Default",
						Contexts:        []string{"Context 1", "Context 2"},
						SyntaxMd:        []string{"I am Syntax"},
						SyntaxHtml:      []string{"<p>I am Syntax</p>"},
						DescriptionMd:   "It is a directive",
						DescriptionHtml: "<p>It is a directive</p>",
					},
				},
			},
		},
		Version: "1.0",
	}

	buf := new(bytes.Buffer)

	err := want.Write(context.Background(), buf)
	require.NoError(t, err)

	jsonData := buf.Bytes()

	require.Contains(t, string(jsonData), "<p>", "doesn't do any wacky encoding")

	var got output.Reference
	err = json.Unmarshal(jsonData, &got)
	require.NoError(t, err)

	require.Equal(t, want, &got)

}

func TestGetVersion(t *testing.T) {
	jsonData := `{"modules": [{ "name": "Module 1", "directives": [] }], "version": "1.0"}`

	got, err := output.GetVersion(context.Background(), strings.NewReader(jsonData))
	require.NoError(t, err)

	want := "1.0"
	require.Equal(t, want, got)
}
