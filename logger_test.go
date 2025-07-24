package log

import (
	"testing"
)

func Test_V(t *testing.T) {
	logger1 := V(0)
	logger1.Info("Test V0")

	opts := InitOptions(WithLevel("debug"))
	logger2 := New(opts).V(1)
	logger2.Info("Test V1")
}
