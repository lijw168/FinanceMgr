package log

import (
	"context"
)

type ILog interface {
	Trace(format string, v ...interface{})
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	Fatal(format string, v ...interface{})

	TraceJson(v interface{})
	DebugJson(v interface{})
	InfoJson(v interface{})
	WarnJson(v interface{})
	ErrorJson(v interface{})
	FatalJson(v interface{})

	LogTrace(v ...interface{})
	LogDebug(v ...interface{})
	LogInfo(v ...interface{})
	LogWarn(v ...interface{})
	LogError(v ...interface{})
	LogFatal(v ...interface{})

	TraceContext(ctx context.Context, format string, v ...interface{})
	DebugContext(ctx context.Context, format string, v ...interface{})
	WarnContext(ctx context.Context, format string, v ...interface{})
	InfoContext(ctx context.Context, format string, v ...interface{})
	ErrorContext(ctx context.Context, format string, v ...interface{})
	FatalContext(ctx context.Context, format string, v ...interface{})

	Close()
	SetLevel(level int)
}

// FIXME: fix Daliy to  Daily
func NewDaliyRotateLog(baseName string) (ILog, error) {
	timehandler, err := NewTimeRotatingFileHandler(baseName, WhenDay, 1)
	if err != nil {
		return nil, err
	}
	l := NewDefault(timehandler)
	return l, nil
}

// func NewFakeLog() (ILog, error) {
// 	return &FakeLog{}, nil
// }

// type FakeLog struct {
// }

// func (f *FakeLog) Trace(format string, v ...interface{}) {}
// func (f *FakeLog) Debug(format string, v ...interface{}) {}
// func (f *FakeLog) Info(format string, v ...interface{})  {}
// func (f *FakeLog) Warn(format string, v ...interface{})  {}
// func (f *FakeLog) Error(format string, v ...interface{}) {}
// func (f *FakeLog) Fatal(format string, v ...interface{}) {}
// func (f *FakeLog) TraceJson(v interface{})               {}
// func (f *FakeLog) DebugJson(v interface{})               {}
// func (f *FakeLog) InfoJson(v interface{})                {}
// func (f *FakeLog) WarnJson(v interface{})                {}
// func (f *FakeLog) ErrorJson(v interface{})               {}
// func (f *FakeLog) FatalJson(v interface{})               {}
// func (f *FakeLog) LogTrace(v ...interface{})             {}
// func (f *FakeLog) LogDebug(v ...interface{})             {}
// func (f *FakeLog) LogInfo(v ...interface{})              {}
// func (f *FakeLog) LogWarn(v ...interface{})              {}
// func (f *FakeLog) LogError(v ...interface{})             {}
// func (f *FakeLog) LogFatal(v ...interface{})             {}
// func (f *FakeLog) TraceContext(ctx context.Context, format string,
// 	v ...interface{}) {
// }
// func (f *FakeLog) DebugContext(ctx context.Context, format string,
// 	v ...interface{}) {
// }
// func (f *FakeLog) WarnContext(ctx context.Context, format string,
// 	v ...interface{}) {
// }
// func (f *FakeLog) InfoContext(ctx context.Context, format string,
// 	v ...interface{}) {
// }
// func (f *FakeLog) ErrorContext(ctx context.Context, format string,
// 	v ...interface{}) {
// }
// func (f *FakeLog) FatalContext(ctx context.Context, format string,
// 	v ...interface{}) {
// }
// func (f *FakeLog) Close()             {}
// func (f *FakeLog) SetLevel(level int) {}
