package config

type TaskType string

const (
	TaskTypeUpload TaskType = "upload"
	TaskTypeCmd    TaskType = "cmd"
	TaskTypeLocalCmd TaskType = "local_cmd"
)

func (r *Task) TaskType() TaskType {
	if r.Upload != nil {
		return TaskTypeUpload
	}
	if r.Cmd != nil {
		return TaskTypeCmd
	}
	if r.LocalCmd != nil {
		return TaskTypeLocalCmd
	}
	return ""
}
