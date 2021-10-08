package tasks

import (
	"github.com/chyroc/anb/internal"
	"github.com/chyroc/anb/internal/config"
)

func RunLocalCmd(task *config.Task, cli *internal.SSH, vals config.Any) error {
	if task.Dir != "" {
		internal.PrintfWhite("\t[local_cmd] dir=%q\n", task.Dir)
	}
	for _, cmd := range task.LocalCmd.Commands {
		internal.PrintfWhite("\t[local_cmd] %q\n", cmd)
		out, err := internal.NewLocalCommand().RunCommandInPipe(joinCmd(task.Dir, cmd))
		if err != nil {
			return err
		}

		appendVals(task.ID, out, vals)
	}

	return nil
}
