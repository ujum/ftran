package data

type ResourceLog struct {
	Source  string
	Target  string
	Skipped bool
	Error   error
}

func NewResourceLog(nestedSourceDir string, nestedTargetDir string, skipped bool, err error) *ResourceLog {
	return &ResourceLog{
		Source:  nestedSourceDir,
		Target:  nestedTargetDir,
		Skipped: skipped,
		Error:   err,
	}
}
