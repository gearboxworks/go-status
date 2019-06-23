package status

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gearboxworks/go-status/only"
	"net/http"
	"strconv"
	"strings"
)

var NilS = (*S)(nil)
var _ Status = NilS

type S struct {
	success    bool
	warn       bool
	cause      error
	httpstatus int
	message    string
	details    string
	data       interface{}
	help       HelpTypeMap
	errorcode  int
	logas      int
}

func NewStatus(args *Args) (sts *S) {
	for range only.Once {
		sts = &S{
			success:    args.Success,
			message:    args.Message,
			httpstatus: args.HttpStatus,
			cause:      args.Cause,
			data:       args.Data,
			help: HelpTypeMap{
				AllHelp: &args.Help,
				ApiHelp: &args.ApiHelp,
				CliHelp: &args.CliHelp,
			},
		}

		if sts.message == "" && sts.cause != nil {
			sts.message = sts.cause.Error()
		}

		if !sts.success && sts.cause == nil {
			sts.cause = errors.New(sts.message)
		}

		if sts.httpstatus == 0 {
			sts.httpstatus = http.StatusInternalServerError
		}

		if *sts.help[AllHelp] == "" {
			help := ContactSupportHelp()
			sts.help[AllHelp] = &help
		}

		if *sts.help[ApiHelp] == "" {
			sts.help[ApiHelp] = sts.help[AllHelp]
		}

		if *sts.help[CliHelp] == "" {
			sts.help[CliHelp] = sts.help[AllHelp]
		}

	}
	return sts
}

///////////// BOOL METHOD(S) ///////////////

func (me *S) IsSuccess() bool {
	return me.success
}

func (me *S) IsError() bool {
	return !me.success
}

func (me *S) IsWarn() bool {
	return me.warn
}

///////////// CHANGE BEHAVIOR METHOD(S) ///////////////

func (me *S) LogAs() int {
	return me.logas
}

///////////// ACTION METHOD(S) ///////////////

func (me *S) Log() {
	Logger.Log(me)
}

///////////// PROPERTY METHOD(S) ///////////////

func (me *S) String() string {
	return me.message
}

func (me *S) Warn() (w string) {
	for range only.Once {
		w = ""
		if !me.warn {
			break
		}
		w = me.message
	}
	return w
}

func (me *S) Error() string {
	return me.message
}

func (me *S) Cause() error {
	return me.cause
}

func (me *S) Message() string {
	return me.message
}

func (me *S) Detail() string {
	return me.details
}

func (me *S) Help() string {
	return me.GetHelp(AllHelp)
}

func (me *S) Data() (data interface{}) {
	return me.data
}

func (me *S) HttpStatus() int {
	return me.httpstatus
}

func (me *S) ErrorCode() int {
	return me.errorcode
}

func (me *S) AllHelp() string {
	return me.GetHelp(AllHelp)
}

func (me *S) ApiHelp() string {
	return me.GetHelp(ApiHelp)
}

func (me *S) CliHelp() string {
	return me.GetHelp(CliHelp)
}

///////////// AGGREGATE METHOD(S) ///////////////

func (me *S) Json() (js []byte) {
	for range only.Once {
		var err error
		js, err = json.Marshal(me.PropertyStringMap())
		if err != nil {
			js = ConvertErrorToJson(err, "property string map")
		}
	}
	return js
}

func (me *S) LongMessage() (lm string) {
	for range only.Once {
		lm = me.message
		if me.cause == nil {
			break
		}
		sts, ok := me.cause.(Status)
		if !ok {
			lm = fmt.Sprintf("%s; %s", lm, me.cause.Error())
			break
		}
		lm = fmt.Sprintf("%s; %s", lm, sts.LongMessage())
	}
	return lm
}

func (me *S) LongFullText() (lft string) {
	for range only.Once {
		lft = me.FullText()
		if me.cause == nil {
			break
		}
		sts, ok := me.cause.(Status)
		if !ok {
			lft = fmt.Sprintf("%s; %s", lft, me.cause.Error())
		}
		lft = fmt.Sprintf("%s; %s", lft, sts.LongFullText())
	}
	return lft
}

func (me *S) FullTextError() (err error) {
	return fmt.Errorf(me.LongFullText())
}

func (me *S) LongError() (err error) {
	return fmt.Errorf(me.LongMessage())
}

func (me *S) FullHelp() (h string) {
	for range only.Once {
		h = me.AllHelp()
		if h == "" {
			break
		}
		m := make(map[string]bool, 1)
		m[h] = true
		for ht, hh := range me.help {
			if *hh == "" {
				continue
			}
			if _, ok := m[*hh]; ok {
				continue
			}
			h = fmt.Sprintf("%s; [%s] %s", h, strings.ToUpper(string(ht)), *hh)
		}
	}
	return h
}

func (me *S) PropertyStringMap() (m PropertyStringMap) {
	var status Property
	if me.warn {
		status = "warning"
	} else if me.success {
		status = "success"
	} else {
		status = "failure"
	}

	var cause Property
	if err, ok := me.cause.(Status); ok {
		cause = Property(ConvertStringMapToJson(err.PropertyStringMap()))
	} else if me.cause != nil {
		cause = err.Error()
	}
	data, err := json.Marshal(me.data)
	if err != nil {
		data = []byte(ConvertErrorToJson(err, "data"))
	}
	m = PropertyStringMap{
		"status":     status,
		"cause":      cause,
		"message":    me.message,
		"details":    me.details,
		"help":       me.FullHelp(),
		"data":       Property(data),
		"httpstatus": strconv.Itoa(me.httpstatus),
		"errorcode":  strconv.Itoa(me.errorcode),
	}
	for ht, h := range me.help {
		ht = fmt.Sprintf("%s_help", ht)
		m[ht] = *h
	}
	return m
}

func (me *S) FullText() (fd string) {
	for range only.Once {
		fd = me.LongMessage()
		if me.details != "" {
			fd = fmt.Sprintf("%s; %s", fd, me.details)
		}
		if d, ok := me.data.(string); ok {
			fd = fmt.Sprintf("%s; %s", fd, d)
		}
		fh := me.FullHelp()
		if fh != "" {
			fd = fmt.Sprintf("%s; %s", fd, fh)
		}
	}
	return fd
}

///////////// GET METHOD(S) ///////////////

func (me *S) GetHelp(helptype HelpType) string {
	h, _ := me.help[helptype]
	return *h
}

///////////// SET METHODS ///////////////

func (me *S) SetLogAs(logAs int) Status {
	me.logas = logAs
	return me
}

func (me *S) SetSuccess(success bool) Status {
	me.success = success
	return me
}

func (me *S) SetWarn(bool) Status {
	me.success = true
	me.warn = true
	return me
}

func (me *S) SetMessage(msg string, args ...interface{}) Status {
	me.message = fmt.Sprintf(msg, args...)
	return me
}

func (me *S) SetDetail(details string, args ...interface{}) Status {
	me.details = fmt.Sprintf(details, args...)
	return me
}

func (me *S) SetHttpStatus(httpstatus int) Status {
	me.httpstatus = httpstatus
	return me
}

func (me *S) SetErrorCode(code int) Status {
	me.errorcode = code
	return me
}

func (me *S) SetData(data interface{}) Status {
	me.data = data
	return me
}

func (me *S) SetCause(err error) Status {
	me.cause = err
	return me
}

func (me *S) SetOtherHelp(help HelpTypeMap) Status {
	for t, h := range help {
		me.help[t] = h
	}
	return me
}

func (me *S) SetHelp(helptype HelpType, help string, args ...interface{}) Status {
	if len(args) > 0 {
		help = fmt.Sprintf(help, args...)
	}
	me.help[helptype] = &help
	if helptype == AllHelp {
		for t := range me.help {
			me.help[t] = &help
		}
	}
	return me
}

func (me *S) SetAllHelp(help string, args ...interface{}) Status {
	return me.SetHelp(AllHelp, help, args...)
}

func (me *S) SetApiHelp(help string, args ...interface{}) Status {
	return me.SetHelp(CliHelp, help, args...)
}

func (me *S) SetCliHelp(help string, args ...interface{}) Status {
	return me.SetHelp(CliHelp, help, args...)
}
