package parse_test

import (
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/nginxinc/ampex-apps/tools/reference-converter/parse"
	"github.com/stretchr/testify/require"
)

func TestPara_ToMarkdown(t *testing.T) {
	t.Parallel()
	testcases := map[string]struct {
		XML  string
		want string
	}{
		"multiple <para>s are combined": {
			XML:  `<para>A</para><para>B</para>`,
			want: "A\n\nB",
		},
		"<literal> tags are code": {
			XML:  `<para>A <literal>B</literal></para>`,
			want: "A `B`",
		},
		"comments are ignored": {
			XML:  `<para>A <!-- B --></para>`,
			want: "A",
		},
		"examples are fences": {
			XML:  `<para><example> A</example></para>`,
			want: "```\n A\n```",
		},
		"unknown tags show a TODO": {
			XML:  `<para><what>??</what></para>`,
			want: "`TODO: handle <what>`",
		},
		// TODO(AMPEX-72): handle more prose-y elements
	}
	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			type doc struct {
				Prose parse.Prose `xml:"para"`
			}

			var d doc
			err := xml.Unmarshal([]byte(fmt.Sprintf("<doc>%s</doc>", tc.XML)), &d)
			require.NoError(t, err)

			got := d.Prose.ToMarkdown()

			require.Equal(t, tc.want, got, "failed on %#v", d.Prose)
		})
	}
}
