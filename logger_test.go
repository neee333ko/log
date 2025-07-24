package log

import (
	"context"
	"sync"
	"testing"
)

func Test_V(t *testing.T) {
	logger1 := V(0)
	logger1.Info("Test V0")

	opts := InitOptions(WithLevel("debug"))
	logger2 := New(opts).V(1)
	logger2.Info("Test V1")
}

func Test_WithValues(t *testing.T) {
	logger := New(nil)
	logger.Info("msg without values")

	logger = logger.WithValues("musician", "kanye")
	logger.Info("msg with values")
}

func Test_WithName(t *testing.T) {
	logger := New(nil)
	logger.Info("msg without name")

	logger = logger.WithNamed("nekologger")
	logger.Info("msg with name")
}

func Test_Context(t *testing.T) {
	logger := New(nil).WithValues("username", "kanye")

	gc := context.Background()
	vc := logger.WithContext(gc)

	var wg sync.WaitGroup
	wg.Add(1)

	go func(ctx context.Context) {
		defer wg.Done()

		logger := FromContext(ctx)
		logger.Info("msg with context")
	}(vc)

	vc2 := context.WithValue(gc, KeyRequestID, "123456")
	logger2 := New(nil).L(vc2)

	logger2.Info("msg with context from L")

	wg.Wait()
}
