package status

type LogType = int

const (
	FatalLog LogType = iota + 1
	ErrorLog
	WarnLog
	DebugLog
)
