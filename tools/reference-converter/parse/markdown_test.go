package parse_test

import (
	"fmt"
	"testing"

	"github.com/nginxinc/ampex-apps/tools/reference-converter/parse"
	"github.com/nginxinc/ampex-apps/tools/reference-converter/tarball"
	"github.com/stretchr/testify/require"
)

func TestMarkdown(t *testing.T) {
	t.Parallel()
	testcases := map[string]struct {
		XML, want string
		opts      []xmlOption
	}{
		"multiple <para>s are combined": {
			XML: `<para>A</para><para>B</para>`,
			want: lines(
				"A",
				"",
				"B"),
			opts: []xmlOption{withPara(false)},
		},
		"<literal> are code": {
			XML:  `A <literal>B</literal>`,
			want: "A `B`",
		},
		"comments are ignored": {
			XML:  `A <!-- B -->`,
			want: "A",
		},
		"<example> are fences": {
			XML: `<example> A</example>`,
			want: lines(
				"```",
				" A",
				"```",
			),
		},
		"unknown tags show a TODO": {
			XML:  `<what>??</what>`,
			want: "`TODO: handle <what>`",
		},
		`<list type="tag">`: {
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
		`<list type="bullet">`: {
			XML: `<list type="bullet">
			<listitem><para>content</para></listitem>
			<listitem>more <literal>content</literal></listitem>
			</list>`,
			want: lines(
				"- content",
				"- more `content`",
			),
		},
		`<list type="enum">`: {
			XML: `<list type="enum">
			<listitem>content</listitem>
			<listitem>more <literal>content</literal></listitem>
			</list>`,
			want: lines(
				"1. content",
				"2. more `content`",
			),
		},
		"nested <list>": {
			XML: `<list type="tag">
			<tag-name>tag</tag-name>
			<tag-desc>
				stuff
				<list type="bullet">
					<listitem>another list!</listitem>
					<listitem>
						but wait
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
		"<header> are quoted": {
			XML:  `<header>User-Agent</header>`,
			want: `"User-Agent"`,
		},
		"<emphasis>": {
			XML: `<example>upstream <emphasis>name</emphasis></example>`,
			want: lines(
				"```",
				"upstream name",
				"```",
			),
		},
		"<http-status>": {
			XML:  `<http-status code="418" text="I'm a teapot"/>`,
			want: "418 (I'm a teapot)",
		},
		"<commercial_version>": {
			XML:  `<commercial_version>title</commercial_version>`,
			want: fmt.Sprintf("[title](%s)", upsellURL),
		},
	}
	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ref, err := parse.Parse([]tarball.File{
				testModuleFile(t, tc.XML, tc.opts...),
			}, baseURL, upsellURL)
			require.NoError(t, err)

			require.Equal(t, 1, len(ref.Modules))
			got := ref.Modules[0].Sections[0].Prose.ToMarkdown()

			require.Equal(t, tc.want, got, "failed on `%s`", tc.XML)
		})
	}
}
