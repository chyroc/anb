package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitShellCommand(t *testing.T) {
	tests := []struct {
		name string
		args string
		want []string
	}{
		{name: "1", args: `a b c`, want: []string{"a", "b", "c"}},
		{name: "2", args: `a b "c d e"`, want: []string{"a", "b", "c d e"}},
		{name: "3", args: `a b " c d e "`, want: []string{"a", "b", " c d e "}},
		{name: "4", args: `  a   b   c `, want: []string{"a", "b", "c"}},
		{name: "5", args: `a b 'c'`, want: []string{"a", "b", "c"}},
		{name: "6", args: `a "b" 'c'`, want: []string{"a", "b", "c"}},
		{name: "7", args: `a "b " 'c '`, want: []string{"a", "b ", "c "}},
		{name: "8", args: `\`, want: []string{`\`}},
		{name: "9", args: "\\", want: []string{`\`}},
		{name: "10", args: `\ ""`, want: []string{`\`, ""}},
		{name: "11", args: `\ "" "c "`, want: []string{`\`, "", "c "}},
		{name: "12", args: `\ "" "c\ "`, want: []string{`\`, "", "c\\ "}},
		{name: "12", args: "'a' 'b c' '\"'", want: []string{`a`, `b c`, `"`}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, SplitShellCommand(tt.args), tt.name+" / "+tt.args)
		})
	}
}
