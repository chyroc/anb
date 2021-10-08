package tasks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_splitWithColon(t *testing.T) {
	as := assert.New(t)

	tests := []struct {
		name  string
		args  string
		want1 string
		want2 string
	}{
		{name: "1", args: "${{GITHUB_TOKEN}}:${{GITHUB_TOKEN}}", want1: "${{GITHUB_TOKEN}}", want2: "${{GITHUB_TOKEN}}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, got2 := splitWithColon(tt.args)
			as.Equal(tt.want1, got1)
			as.Equal(tt.want2, got2)
		})
	}
}
