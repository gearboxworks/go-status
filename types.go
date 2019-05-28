package status

type PropertyStringMap map[Property]string
type Property = string

type Status interface {
	Cause() error
	Data() interface{}
	Detail() string
	Error() string
	ErrorCode() int
	FullText() string
	FullTextError() error
	GetHelp(HelpType) string
	Help() string
	HttpStatus() int
	IsError() bool
	IsSuccess() bool
	IsWarn() bool
	Log()
	LogAs() LogType
	LongError() error
	LongFullText() string
	LongMessage() string
	Message() string
	SetCause(error) Status
	SetData(interface{}) Status
	SetDetail(string, ...interface{}) Status
	SetErrorCode(int) Status
	SetHelp(HelpType, string, ...interface{}) Status
	SetHttpStatus(int) Status
	SetLogAs(LogType) Status
	SetMessage(string, ...interface{}) Status
	SetOtherHelp(HelpTypeMap) Status
	SetSuccess(bool) Status
	SetWarn(bool) Status
	PropertyStringMap() PropertyStringMap
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
	Log(error)
	Debug(Msg)
	Warn(Msg)
	Error(Msg)
	Fatal(Msg)
}

type errObj struct {
	Error string `json:"error"`
	Cause string `json:"cause"`
}
