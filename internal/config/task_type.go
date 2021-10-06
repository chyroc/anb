package config

type TaskType string

const (
	TaskTypeCopy     TaskType = "copy"
	TaskTypeCmd      TaskType = "cmd"
	TaskTypeLocalCmd TaskType = "local_cmd"
)

func (r *ConfigTask) TaskType() TaskType {
	if r.Copy != nil {
		return TaskTypeCopy
	}
	if r.Cmd != nil {
		return TaskTypeCmd
	}
	if r.LocalCmd != nil {
		return TaskTypeLocalCmd
	}
	return ""
}
