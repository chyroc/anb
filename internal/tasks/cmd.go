package tasks

import (
	"fmt"

	"github.com/chyroc/anb/internal"
	"github.com/chyroc/anb/internal/config"
)

func RunCmd(task *config.Task, cli *internal.SSH, vals config.Any) error {
	if task.Dir != "" {
		internal.PrintfWhite("\t[cmd] dir=%q\n", task.Dir)
	}
	for _, cmd := range task.Cmd.Commands {
		internal.PrintfWhite("\t[cmd] %q\n", cmd)
		out, err := cli.RunInPipe(joinCmd(task.Dir, cmd))
		if err != nil {
			return err
		}

		appendVals(task.ID, out, vals)
	}
	return nil
}

func appendVals(taskID string, msg string, vals config.Any) {
	if output := parseOutput(msg); len(output) > 0 {
		if _, ok := vals["$tasks"].(config.Any)[taskID]; !ok {
			vals["$tasks"].(config.Any)[taskID] = config.Any{}
		}
		for k, v := range output {
			vals["$tasks"].(config.Any)[taskID].(config.Any)[k] = v
		}
	}
}

func joinCmd(dir, cmd string) string {
	if dir == "" {
		return cmd
	}
	return fmt.Sprintf("cd %s; %s", dir, cmd)
}
