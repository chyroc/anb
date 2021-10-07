package tasks

import (
	"fmt"

	"github.com/chyroc/anb/internal"
	"github.com/chyroc/anb/internal/config"
	"github.com/chyroc/anb/internal/task_command"
)

func RunCmd(task *config.Task, cli *internal.SSHCommand, vals config.Any) error {
	for _, cmd := range task.Cmd.Commands {
		internal.PrintfWhite("\t[cmd] %q\n", cmd)
		out, err := cli.Run(cmd)
		if err != nil {
			return err
		}
		fmt.Print(out)

		if output := task_command.ParseOutput(out); len(output) > 0 {
			if _, ok := vals["$tasks"].(config.Any)[task.ID]; !ok {
				vals["$tasks"].(config.Any)[task.ID] = config.Any{}
			}
			for k, v := range output {
				vals["$tasks"].(config.Any)[task.ID].(config.Any)[k] = v
			}
		}
	}
	return nil
}
