package tasks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_setOutput(t *testing.T) {
	as := assert.New(t)

	t.Run("", func(t *testing.T) {
		res := parseOutput(`::set-output name=SELECTED_COLOR::green`)
		as.Equal("green", res["SELECTED_COLOR"])
	})

	t.Run("", func(t *testing.T) {
		res := parseOutput(`::set-output name=SELECTED_COLOR::green space`)
		as.Equal("green space", res["SELECTED_COLOR"])
	})
}
