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
