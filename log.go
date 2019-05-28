package status

import (
	"fmt"
	"github.com/gearboxworks/go-status/only"
	"log"
	"os"
)

var Logger MsgLogger

func init() {
	Logger = &L{}
}

var NilL = (*L)(nil)
var _ MsgLogger = NilL

type L struct{}

func (me *L) Debug(msg Msg) {
	fmt.Printf("[DEBUG] %s\n", msg)
}

func (me *L) Warn(msg Msg) {
	fmt.Printf("[WARN] %s\n", msg)
}

func (me *L) Error(msg Msg) {
	fmt.Printf("[ERROR] %s\n", msg)
}

func (me *L) Fatal(msg Msg) {
	fmt.Printf("[FATAL] %s\n", msg)
	os.Exit(1)
}

func (me *L) Log(err error) {
	for range only.Once {
		if err == nil {
			break
		}
		sts, ok := err.(Status)
		if !ok {
			Logger.Error(err.Error())
			break
		}
		if me == nil {
			log.Fatal("status.Logger is nil\n")
		}
		if sts.LogAs() == FatalLog {
			Logger.Fatal(sts.Message())
		}
		if sts.LogAs() == ErrorLog {
			Logger.Error(sts.Message())
			break
		}
		if sts.LogAs() == WarnLog {
			Logger.Warn(sts.Message())
			break
		}
		if sts.LogAs() == DebugLog {
			Logger.Debug(sts.Message())
			break
		}
		if IsError(sts) {
			Logger.Error(sts.Message())
			break
		}
		if IsWarn(sts) {
			Logger.Warn(sts.Message())
			break
		}
		if IsSuccess(sts) {
			Logger.Debug(sts.Message())
			break
		}
	}
}
