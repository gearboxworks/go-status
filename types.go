package status

type jsonS struct {
	Message string      `json:"message"`
	Help    string      `json:"help,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type Status interface {
	IsSuccess() bool
	IsWarning() bool
	IsError() bool
	Error() string
	Data() interface{}
	Help() string
	Detail() string
	GetHelp(HelpType) string
	HttpStatus() int
	Message() string
	Warning() string
	Cause() error
	FullError() error
	ErrorCode() int
	SetData(interface{}) Status
	SetDetail(string, ...interface{}) Status
	SetCause(error) Status
	SetSuccess(bool) Status
	SetWarning(bool) Status
	SetMessage(string, ...interface{}) Status
	SetHttpStatus(int) Status
	SetHelp(HelpType, string, ...interface{}) Status
	SetOtherHelp(HelpTypeMap) Status
	SetErrorCode(int) Status
	GetString() (string, Status)
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
