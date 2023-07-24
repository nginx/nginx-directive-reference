package parse_test

import (
	"fmt"
	"testing"

	"github.com/nginxinc/ampex-apps/tools/reference-converter/parse"
	"github.com/nginxinc/ampex-apps/tools/reference-converter/tarball"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarkdown(t *testing.T) {
	t.Parallel()
	testcases := map[string]struct {
		content, wantContent string
		syntax, wantSyntax   string
		opts                 []xmlOption
	}{
		"multiple <para>s are combined": {
			content: `<para>A</para><para>B</para>`,
			wantContent: lines(
				"A",
				"",
				"B"),
			opts: []xmlOption{withPara(false)},
		},
		"<literal> are code": {
			content:     `A <literal>B</literal>`,
			wantContent: "A `B`",
		},
		"comments are ignored": {
			content:     `A <!-- B -->`,
			wantContent: "A",
		},
		"<example> are fences": {
			content: `<example> A</example>`,
			wantContent: lines(
				"```",
				" A",
				"```",
			),
		},
		"unknown tags show a TODO": {
			content:     `<what>??</what>`,
			wantContent: "`TODO: handle <what>`",
		},
		`<list type="tag">`: {
			content: `<list type="tag">
			<tag-name>tag <literal>one</literal></tag-name>
			<tag-desc>contents</tag-desc>
			<tag-name>tag two</tag-name>
			<tag-desc>more <var>contents</var></tag-desc>
			</list>`,
			wantContent: lines(
				"- tag `one`",
				"",
				"  contents",
				"- tag two",
				"",
				"  more `contents`",
			),
		},
		`<list type="bullet">`: {
			content: `<list type="bullet">
			<listitem><para>content</para></listitem>
			<listitem>more <literal>content</literal></listitem>
			</list>`,
			wantContent: lines(
				"- content",
				"- more `content`",
			),
		},
		`<list type="enum">`: {
			content: `<list type="enum">
			<listitem>content</listitem>
			<listitem>more <literal>content</literal></listitem>
			</list>`,
			wantContent: lines(
				"1. content",
				"2. more `content`",
			),
		},
		"nested <list>": {
			content: `<list type="tag">
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
			wantContent: lines(
				"- tag",
				"",
				"  stuff",
				"  - another list!",
				"  - but wait",
				"    1. there's more!",
			),
		},
		"<header> are quoted": {
			content:     `<header>User-Agent</header>`,
			wantContent: `"User-Agent"`,
		},
		"<emphasis>": {
			content: `<example>upstream <emphasis>name</emphasis></example>`,
			wantContent: lines(
				"```",
				"upstream name",
				"```",
			),
		},
		"<http-status>": {
			content:     `<http-status code="418" text="I'm a teapot"/>`,
			wantContent: "418 (I'm a teapot)",
		},
		"<commercial_version>": {
			content:     `<commercial_version>title</commercial_version>`,
			wantContent: fmt.Sprintf("[title](%s)", upsellURL),
		},
		"<syntax> enum": {
			syntax:     "<literal>enumA</literal> | <literal>enumB</literal>",
			wantSyntax: "`enumA` | `enumB`",
		},
		"<syntax> arg": {
			syntax:     "<value>arg</value>",
			wantSyntax: "*`arg`*",
		},
		"<syntax> args": {
			syntax:     "<value>argA</value> <value>argB</value>",
			wantSyntax: "*`argA`* *`argB`*",
		},
		"<syntax> optional flags": {
			syntax: lines(
				"[<literal>flagA</literal>]",
				"[<literal>flagB</literal>]",
				"[<literal>flagC</literal>]"),
			wantSyntax: "[`flagA`] [`flagB`] [`flagC`]",
		},
		"<syntax> arg with optional flag": {
			syntax:     "<value>arg</value> [<literal>flag</literal>]",
			wantSyntax: "*`arg`* [`flag`]",
		},
		"<syntax> arg or flag": {
			syntax:     "<value>argA</value> | <value>argB</value> | <literal>flag</literal>",
			wantSyntax: "*`argA`* | *`argB`* | `flag`",
		},
		"<syntax> arg with optional flag or flag": {
			syntax:     "<value>arg</value> [<literal>flagA</literal>] | <literal>flagB</literal>",
			wantSyntax: "*`arg`* [`flagA`] | `flagB`",
		},
		"<syntax> enum with optional flag": {
			syntax:     "<literal>enumA</literal> | <literal>enumB</literal> [<literal>flag</literal>]",
			wantSyntax: "`enumA` | `enumB` [`flag`]",
		},
		"<syntax> arg with named options": {
			syntax:     "<value>arg</value> [<literal>opt</literal>=<value>val</value>]",
			wantSyntax: "*`arg`* [`opt`=*`val`*]",
		},
		"<syntax> multi line indented": {
			syntax: lines(
				"",
				"    [<literal>SSLv2</literal>]",
				"    [<literal>SSLv4</literal>]",
			),
			wantSyntax: "[`SSLv2`] [`SSLv4`]",
		},
		"<note>": {
			content:     "<note>Hey, I'm <value>important</value></note>",
			wantContent: "> Hey, I'm *`important`*",
		},
		"<note> multi-line": {
			content: lines(
				"<note>",
				"The <literal>TLSv1.1</literal> and <literal>TLSv1.2</literal> parameters",
				"(1.1.13, 1.0.12) work only when OpenSSL 1.0.1 or higher is used.",
				"</note>"),
			wantContent: lines(
				"> The `TLSv1.1` and `TLSv1.2` parameters",
				"> (1.1.13, 1.0.12) work only when OpenSSL 1.0.1 or higher is used.",
			),
		},
	}
	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			opts := append(tc.opts, withContent(tc.content), withSyntax(tc.syntax))
			f := testModuleFile(t, opts...)
			ref, err := parse.Parse([]tarball.File{f}, baseURL, upsellURL)
			require.NoError(t, err)

			require.Equal(t, 1, len(ref.Modules))
			gotContent := ref.Modules[0].Sections[0].Directives[0].Prose.ToMarkdown()
			gotSyntax := ref.Modules[0].Sections[0].Directives[0].Syntax.ToMarkdown()

			assert.Equal(t, tc.wantContent, gotContent, "failed on `%s`", tc.content)
			assert.Equal(t, tc.wantSyntax, gotSyntax, "failed on `%s`", tc.syntax)
		})
	}
}
