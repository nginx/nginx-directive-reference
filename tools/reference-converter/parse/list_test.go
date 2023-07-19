package parse_test

import (
	"strings"
	"testing"

	"github.com/nginxinc/ampex-apps/tools/reference-converter/parse"
	"github.com/nginxinc/ampex-apps/tools/reference-converter/tarball"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Parallel()
	testcases := map[string]struct {
		XML, want string
	}{
		"tag list": {
			XML: `<list type="tag">
			<tag-name>tag <literal>one</literal></tag-name>
			<tag-desc>contents</tag-desc>
			<tag-name>tag two</tag-name>
			<tag-desc>more <var>contents</var></tag-desc>
			</list>`,
			want: lines(
				"- tag `one`",
				"",
				"  contents",
				"- tag two",
				"",
				"  more `contents`",
			),
		},
		"bullet list": {
			XML: `<list type="bullet">
			<listitem><para>content</para></listitem>
			<listitem>more <literal>content</literal></listitem>
			</list>`,
			want: lines(
				"- content",
				"- more `content`",
			),
		},
		"enum list": {
			XML: `<list type="enum">
			<listitem>content</listitem>
			<listitem>more <literal>content</literal></listitem>
			</list>`,
			want: lines(
				"1. content",
				"2. more `content`",
			),
		},
		"nested list": {
			XML: `<list type="tag">
			<tag-name>tag</tag-name>
			<tag-desc> stuff
				<list type="bullet">
					<listitem>another list!</listitem>
					<listitem> but wait
						<list type="enum">
							<listitem>there's more!</listitem>
						</list>
					</listitem>
				</list>
			</tag-desc>
			</list>`,
			want: lines(
				"- tag",
				"",
				"  stuff",
				"  - another list!",
				"  - but wait",
				"    1. there's more!",
			),
		},
	}
	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ref, err := parse.Parse([]tarball.File{
				testModuleFile("test.xml", tc.XML),
			}, baseURL)
			require.NoError(t, err)

			require.Equal(t, 1, len(ref.Modules))
			got := ref.Modules[0].Sections[0].Prose.ToMarkdown()

			require.Equal(t, tc.want, got, "failed on `%s`", tc.XML)
		})
	}
}

func lines(l ...string) string { return strings.Join(l, "\n") }
