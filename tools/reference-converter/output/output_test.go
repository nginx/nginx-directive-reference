package output_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/nginxinc/ampex-apps/tools/reference-converter/output"
	"github.com/nginxinc/ampex-apps/tools/reference-converter/parse"
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
				Name: "1",
			},
			{
				Name: "2",
				Directives: []output.Directive{
					{
						Name:        "directive 2",
						Default:     "default 2",
						Contexts:    []string{"context 1", "context 2"},
						Syntax:      []string{"syntax 1", "syntax 2"},
						Description: "Test",
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
						Name:        "Directive 1",
						Default:     "Default",
						Contexts:    []string{"Context 1", "Context 2"},
						Syntax:      []string{"I am Syntax"},
						Description: "It is a directive",
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
