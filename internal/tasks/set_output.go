package tasks

import (
	"regexp"
	"strings"

	"github.com/chyroc/anb/internal/config"
)

// https://docs.github.com/en/actions/learn-github-actions/workflow-commands-for-github-actions
// echo "::workflow-command parameter1={data},parameter2={data}::{command value}"

// ::set-output name=SELECTED_COLOR::green
func parseOutput(s string) config.Any {
	res := config.Any{}
	for _, v := range strings.Split(s, "\n") {
		vv := regSetOutput.FindStringSubmatch(v)
		if len(vv) == 3 {
			res[vv[1]] = vv[2]
		}
	}

	return res
}

var regSetOutput = regexp.MustCompile(`^::set-output name=(.*?)::(.*?)$`)
