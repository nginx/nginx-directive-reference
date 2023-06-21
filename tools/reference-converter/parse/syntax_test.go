package parse_test

import (
	"encoding/xml"
	"fmt"
	"testing"

	"github.com/nginxinc/ampex-apps/tools/reference-converter/parse"
	"github.com/stretchr/testify/require"
)

func TestSyntax(t *testing.T) {
	t.Parallel()
	testcases := map[string]struct {
		XML  string
		want string
	}{
		"enum": {
			XML:  "<literal>enumA</literal> | <literal>enumB</literal>",
			want: "`enumA` | `enumB`",
		},
		"arg": {
			XML:  "<value>arg</value>",
			want: "*`arg`*",
		},
		"args": {
			XML:  "<value>argA</value> <value>argB</value>",
			want: "*`argA`* *`argB`*",
		},
		"optional flags": {
			XML:  "[<literal>flagA</literal>]\n[<literal>flagB</literal>]\n[<literal>flagC</literal>]",
			want: "[`flagA`]\n[`flagB`]\n[`flagC`]",
		},
		"arg with optional flag": {
			XML:  "<value>arg</value> [<literal>flag</literal>]",
			want: "*`arg`* [`flag`]",
		},
		"arg or flag": {
			XML:  "<value>argA</value> | <value>argB</value> | <literal>flag</literal>",
			want: "*`argA`* | *`argB`* | `flag`",
		},
		"arg with optional flag or flag": {
			XML:  "<value>arg</value> [<literal>flagA</literal>] | <literal>flagB</literal>",
			want: "*`arg`* [`flagA`] | `flagB`",
		},
		"enum with optional flag": {
			XML:  "<literal>enumA</literal> | <literal>enumB</literal> [<literal>flag</literal>]",
			want: "`enumA` | `enumB` [`flag`]",
		},
		"arg with named options": {
			XML:  "<value>arg</value> [<literal>opt</literal>=<value>val</value>]",
			want: "*`arg`* [`opt`=*`val`*]",
		},
	}
	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			type doc struct {
				Got parse.Syntax `xml:"syntax"`
			}

			var d doc
			err := xml.Unmarshal([]byte(fmt.Sprintf("<doc><syntax>%s</syntax></doc>", tc.XML)), &d)
			require.NoError(t, err)

			require.Equal(t, tc.want, d.Got.Content)
		})
	}
}
