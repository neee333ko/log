package klog

import (
	"flag"

	"go.uber.org/zap"
	"k8s.io/klog"
)

func InitKlog(l *zap.Logger) {
	fs := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(fs)
	defer klog.Flush()

	klog.SetOutputBySeverity("INFO", &infoLogger{logger: l})
	klog.SetOutputBySeverity("WARNING", &warnLogger{logger: l})
	klog.SetOutputBySeverity("ERROR", &errorLogger{logger: l})
	klog.SetOutputBySeverity("FATAL", &fatalLogger{logger: l})

	_ = fs.Set("skip_headers", "true")
	_ = fs.Set("logtostderr", "false")

}

type infoLogger struct {
	logger *zap.Logger
}

func (l *infoLogger) Write(p []byte) (n int, err error) {
	l.logger.Info(string(p[:len(p)-1]))

	return len(p), nil
}

type warnLogger struct {
	logger *zap.Logger
}

func (l *warnLogger) Write(p []byte) (n int, err error) {
	l.logger.Warn(string(p[:len(p)-1]))

	return len(p), nil
}

type errorLogger struct {
	logger *zap.Logger
}

func (l *errorLogger) Write(p []byte) (n int, err error) {
	l.logger.Error(string(p[:len(p)-1]))

	return len(p), nil
}

type fatalLogger struct {
	logger *zap.Logger
}

func (l *fatalLogger) Write(p []byte) (n int, err error) {
	l.logger.Fatal(string(p[:len(p)-1]))

	return len(p), nil
}
