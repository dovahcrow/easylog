package logs

import (
	"testing"
)

func Test_xiang(t *testing.T) {
	Info(`nothing`)
	Trace(`nothing`)
	Debug(`nothing`)
	Warn(`nothing`)
	Critical(`nothing`)
	Error("nothing")
}
