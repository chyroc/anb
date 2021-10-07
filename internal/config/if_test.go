package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIfExpr_IsTrue(t *testing.T) {
	as := assert.New(t)

	tests := []struct {
		name       string
		r          IfExpr
		val        Any
		want       bool
		ErrContain string
	}{
		// {name: "", r: "", want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.IsTrue(Any{})
			if tt.ErrContain != "" {
				as.NotNil(err)
				as.Contains(err.Error(), tt.ErrContain)
			} else {
				as.Nil(err)
				as.Equal(tt.want, got)
			}
		})
	}
}

func Test_RunBoolScript(t *testing.T) {
	as := assert.New(t)

	chTestDir()

	tests := []struct {
		name       string
		script     string
		val        Any
		want       bool
		ErrContain string
	}{
		{name: "1", script: "false", want: false},
		{name: "2", script: "true", want: true},

		{name: "3", script: `exist("./README.md")`, want: true},
		{name: "4", script: `!exist("./README.md")`, want: false},
		{name: "5", script: `exist("./fake.data")`, want: false},
		{name: "6", script: `!exist("./fake.data")`, want: true},

		{name: "7", script: `exist("./README.md") && exist("./LICENSE")`, want: true},
		{name: "8", script: `exist("./README.md") && exist("./fake.data")`, want: false},
		{name: "9", script: `exist("./README.md") || exist("./fake.data")`, want: true},
		{name: "10", script: `(exist("./README.md") || exist("./fake.data")) && exist("./README.md")`, want: true},
		{name: "11", script: `(exist("./README.md") || exist("./fake.data")) && exist("./fake.data")`, want: false},

		{name: "12", script: "1 == 1", want: true},
		{name: "13", script: "1 == 2", want: false},
		{name: "14", script: `"x" == "x"`, want: true},
		{name: "15", script: `"x" == "y"`, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExecBoolScript(tt.script, tt.val)
			if tt.ErrContain != "" {
				as.NotNil(err, tt.name)
				as.Contains(err.Error(), tt.ErrContain, tt.name)
			} else {
				as.Nil(err, tt.name)
				as.Equal(tt.want, got, tt.name)
			}
		})
	}
}
