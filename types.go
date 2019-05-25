package status

type jsonS struct {
	Message string      `json:"message"`
	Help    string      `json:"help,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type Status interface {
	Cause() error
	Data() interface{}
	Additional() string
	Error() string
	ErrorCode() int
	FullError() error
	GetFullDetails() string
	GetFullMessage() string
	GetHelp(HelpType) string
	Help() string
	HttpStatus() int
	IsError() bool
	IsSuccess() bool
	IsWarn() bool
	Log()
	Message() string
	SetCause(error) Status
	SetData(interface{}) Status
	SetAdditional(string, ...interface{}) Status
	SetErrorCode(int) Status
	SetHelp(HelpType, string, ...interface{}) Status
	SetHttpStatus(int) Status
	SetMessage(string, ...interface{}) Status
	SetOtherHelp(HelpTypeMap) Status
	SetSuccess(bool) Status
	SetWarn(bool) Status
	Warn() string
}

type Args struct {
	Success    bool
	Warning    bool
	Help       string
	ApiHelp    string
	CliHelp    string
	OtherHelp  HelpTypeMap
	Message    string
	HttpStatus int
	Cause      error
	Data       interface{}
}

type SuccessInspector interface {
	IsSuccess() bool
}

type (
	Msg = string
)
type MsgLogger interface {
	Debug(Msg)
	Warn(Msg)
	Fatal(Msg)
}
