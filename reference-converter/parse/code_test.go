package parse_test

import (
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/nginxinc/nginx-directive-reference/reference-converter/parse"
	"github.com/stretchr/testify/require"
)

func TestCodeBlocks(t *testing.T) {
	t.Parallel()
	testcases := map[string]struct {
		XML, want string
	}{
		"literal": {XML: `<literal>lit</literal>`, want: "`lit`"},
		"var":     {XML: `<var>var</var>`, want: "`var`"},
		"command": {XML: `<command>cmd</command>`, want: "`cmd`"},
		"path":    {XML: `<path>path</path>`, want: "`path`"},
		"c-def":   {XML: `<c-def>A=1</c-def>`, want: "`A=1`"},
		"c-func":  {XML: `<c-func>fn</c-func>`, want: "`fn()`"},
		"value":   {XML: `<value>val</value>`, want: "*`val`*"},
	}
	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			type doc struct {
				Got parse.Paragraph `xml:"para"`
			}

			var d doc
			err := xml.Unmarshal([]byte(fmt.Sprintf("<doc><para>%s</para></doc>", tc.XML)), &d)
			require.NoError(t, err)

			got := d.Got.ToMarkdown()

			require.Equal(t, tc.want, got)
		})
	}
}
