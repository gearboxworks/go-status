package status

import (
	"encoding/json"
	"fmt"
	"github.com/gearboxworks/go-status/only"
	"log"
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
	additional string
	data       interface{}
	help       HelpTypeMap
	errorcode  int
	logto      int
}

func (me *S) GetFullDetails() (fd string) {
	for range only.Once {
		fd = me.GetFullMessage()
		if me.additional != "" {
			fd = fmt.Sprintf("%s; %s", fd, me.additional)
		}
		if d, ok := me.data.(string); ok {
			fd = fmt.Sprintf("%s; %s", fd, d)
		}
		if d, ok := me.data.(string); ok {
			fd = fmt.Sprintf("%s; %s", fd, d)
		}
		fh := me.GetFullHelp()
		if fh != "" {
			fd = fmt.Sprintf("%s; %s", fd, fh)
		}
	}
	return fd
}

func (me *S) Log() {
	if Logger == nil {
		log.Fatal("status.Logger is nil")
	}
	for range only.Once {
		if me == nil {
			break
		}
		if me.logto == FatalLog {
			Logger.Fatal(me.Message())
		}
		if me.logto == ErrorLog {
			Logger.Error(me.Message())
			break
		}
		if me.logto == WarnLog {
			Logger.Warn(me.Message())
			break
		}
		if me.logto == DebugLog {
			Logger.Debug(me.Message())
			break
		}
		if IsError(me) {
			Logger.Error(me.Message())
			break
		}
		if IsWarn(me) {
			Logger.Warn(me.Message())
			break
		}
		if IsSuccess(me) {
			Logger.Debug(me.Message())
			break
		}
	}
}

func (me *S) IsWarn() bool {
	return me.warn
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

func (me *S) SetWarn(bool) Status {
	me.success = true
	me.warn = true
	return me
}

func (me *S) Json() []byte {
	js, _ := json.Marshal(&jsonS{
		Message: me.message,
		Help:    *me.help[ApiHelp],
		Data:    me.data,
	})
	return js
}

func (me *S) Error() string {
	return me.message
}

func (me *S) IsSuccess() bool {
	return me.success
}

func (me *S) IsError() bool {
	return !me.success
}

func (me *S) LogTo() int {
	return me.logto
}

func (me *S) Cause() error {
	return me.cause
}

func (me *S) Message() string {
	return me.message
}

func (me *S) Additional() string {
	return me.additional
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

func (me *S) GetFullMessage() (fm string) {
	for range only.Once {
		fm = me.message
		if me.cause == nil {
			break
		}
		sts, ok := me.cause.(Status)
		if !ok {
			fm = fmt.Sprintf("%s; %s", fm, me.cause.Error())
		}
		fm = fmt.Sprintf("%s; %s", fm, sts.GetFullMessage())
	}
	return fm
}

func (me *S) FullError() (err error) {
	msg := me.message
	c := me.cause
	for {
		var ok bool
		c, ok = c.(error)
		if !ok {
			break
		}
		s := c.Error()
		if s != msg {
			msg = fmt.Sprintf("%s; %s", s, msg)
		}
		sts, ok := c.(Status)
		if !ok {
			break
		}
		c = sts.Cause()
	}
	return fmt.Errorf(msg)
}

func (me *S) SetLogTo(logto int) Status {
	me.logto = logto
	return me
}

func (me *S) SetSuccess(success bool) Status {
	me.success = success
	return me
}

func (me *S) SetMessage(msg string, args ...interface{}) Status {
	me.message = fmt.Sprintf(msg, args...)
	return me
}

func (me *S) SetAdditional(details string, args ...interface{}) Status {
	me.additional = fmt.Sprintf(details, args...)
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
	return me.SetHelp(AllHelp, help, args)
}

func (me *S) SetApiHelp(help string, args ...interface{}) Status {
	return me.SetHelp(CliHelp, help, args)
}

func (me *S) SetCliHelp(help string, args ...interface{}) Status {
	return me.SetHelp(CliHelp, help, args)
}

func (me *S) GetFullHelp() (h string) {
	for range only.Once {
		h = me.GetAllHelp()
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

func (me *S) GetHelp(helptype HelpType) string {
	h, _ := me.help[helptype]
	return *h
}

func (me *S) GetAllHelp() string {
	return me.GetHelp(AllHelp)
}

func (me *S) GetApiHelp() string {
	return me.GetHelp(ApiHelp)
}

func (me *S) GetCliHelp(helptype HelpType) string {
	return me.GetHelp(CliHelp)
}

func (me *S) String() string {
	s := me.message
	if me.cause != nil && me.cause.Error() != me.message {
		s = fmt.Sprintf("%s: %s", s, me.cause.Error())
	}
	return s
}
