package status

import (
	"encoding/json"
	"fmt"
	"github.com/gearboxworks/go-status/only"
	"net/http"
)

func Log(err error) {
	Logger.Log(err)
}

func HttpStatus(err error) int {
	sts, ok := err.(Status)
	if !ok {
		return 0
	}
	return sts.HttpStatus()
}

func Message(err error) string {
	sts, ok := err.(Status)
	if !ok {
		return ""
	}
	return sts.Message()
}

func Cause(err error) error {
	sts, ok := err.(Status)
	if !ok {
		return err
	}
	return sts.Cause()
}

func Help(err error, helptype ...HelpType) string {
	sts, ok := err.(Status)
	if !ok {
		return ""
	}
	if len(helptype) == 0 {
		helptype = []HelpType{AllHelp}
	}
	return sts.GetHelp(helptype[0])
}

func IsError(err error) bool {
	sts, ok := err.(Status)
	if !ok {
		return err != nil
	}
	return sts.IsError()
}

func IsSuccess(err error) bool {
	sts, ok := err.(Status)
	if !ok {
		return err == nil
	}
	return !sts.IsError() || sts.HttpStatus() == 0
}

func IsWarn(err error) bool {
	sts, ok := err.(Status)
	if !ok {
		return err == nil
	}
	return sts.IsWarn()
}

func Wrap(err error, args ...*Args) Status {
	var _args *Args
	if len(args) == 0 || args[0] == nil {
		_args = &Args{}
	} else {
		_args = args[0]
	}
	_args.Cause = err
	if _args.Message == "" {
		_args.Message = err.Error()
	}
	return NewStatus(_args)
}

func SimpleSuccess() Status {
	return Success("everything a-ok")
}

func Success(msg string, args ...interface{}) Status {
	return NewStatus(&Args{
		Success:    true,
		Message:    fmt.Sprintf(msg, args...),
		HttpStatus: http.StatusOK,
	})
}

func Fail(args ...*Args) Status {
	var _args *Args
	if len(args) == 0 || args[0] == nil {
		_args = &Args{}
	} else {
		_args = args[0]
	}
	_args.Success = false
	return NewStatus(_args)
}

func Warn(msg string, args ...interface{}) Status {
	return Success(msg, args...).SetWarn(true)
}

func YourBad(msg string, args ...interface{}) Status {
	return Fail(&Args{
		Message:    fmt.Sprintf(msg, args...),
		HttpStatus: http.StatusBadRequest,
	})
}

func OurBad(msg string, args ...interface{}) Status {
	return Fail(&Args{
		Message:    fmt.Sprintf(msg, args...),
		HttpStatus: http.StatusInternalServerError,
	})
}

func ContactSupportHelp() string {
	return "contact support"
}

func ConvertStringMapToJson(psm PropertyStringMap) (js []byte) {
	for range only.Once {
		var err error
		js, err = json.Marshal(psm)
		if err == nil {
			break
		}
		js = ConvertErrorToJson(err, "property string map")
	}
	return js
}

func ConvertErrorToJson(err error, what string) (js []byte) {
	js, err = json.Marshal(err)
	if err != nil {
		js, _ = json.Marshal(errObj{
			Error: "cannot unmarshal " + what,
			Cause: err.Error(),
		})
	}
	return js
}
