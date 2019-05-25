package status

import (
	"fmt"
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
	fmt.Printf("[DEBUG] %s", msg)
}
func (me *L) Warn(msg Msg) {
	fmt.Printf("[WARN] %s", msg)
}
func (me *L) Error(msg Msg) {
	fmt.Printf("[ERROR] %s", msg)
}
func (me *L) Fatal(msg Msg) {
	fmt.Printf("[FATAL] %s", msg)
	os.Exit(1)
}
