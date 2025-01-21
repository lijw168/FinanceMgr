//go:build release
// +build release

package log

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// log level, from low to high, more high means more serious
const (
	LevelTrace = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

const (
	Ltime  = 1 << iota //time format "2006/01/02 15:04:05"
	Lfile              //file.go:123
	Llevel             //[Trace|Debug|Info...]
)

var LevelName [6]string = [6]string{"Trace", "Debug", "Info", "Warn", "Error", "Fatal"}

const TimeFormat = "2006/01/02 15:04:05"

const maxBufPoolSize = 16

type Logger struct {
	sync.Mutex

	level int
	flag  int

	handler Handler

	quit     chan struct{}
	quitDone chan struct{}
	msg      chan []byte

	bufs [][]byte
}

// new a logger with specified handler and flag
func New(handler Handler, flag int) *Logger {
	var l = new(Logger)

	l.level = LevelInfo
	l.handler = handler

	l.flag = flag

	l.quit = make(chan struct{})

	l.quitDone = make(chan struct{})

	l.msg = make(chan []byte, 1024)

	l.bufs = make([][]byte, 0, 16)

	go l.run()

	return l
}

// new a default logger with specified handler and flag: Ltime|Lfile|Llevel
func NewDefault(handler Handler) *Logger {
	return New(handler, Ltime|Lfile|Llevel)
}

func newStdHandler() *StreamHandler {
	h, _ := NewStreamHandler(os.Stdout)
	return h
}

var std = NewDefault(newStdHandler())

func (l *Logger) run() {
	for {
		select {
		case msg := <-l.msg:
			l.handler.Write(msg)
			l.putBuf(msg)
		case <-l.quit:
			l.drainMsg()
			close(l.quitDone)
			return
		}
	}
}

func (l *Logger) drainMsg() {
	for len(l.msg) > 0 {
		msg := <-l.msg
		l.handler.Write(msg)
	}
}

func (l *Logger) popBuf() []byte {
	l.Lock()
	var buf []byte
	if len(l.bufs) == 0 {
		buf = make([]byte, 0, 1024)
	} else {
		buf = l.bufs[len(l.bufs)-1]
		l.bufs = l.bufs[0 : len(l.bufs)-1]
	}
	l.Unlock()

	return buf
}

func (l *Logger) putBuf(buf []byte) {
	l.Lock()
	if len(l.bufs) < maxBufPoolSize {
		buf = buf[0:0]
		l.bufs = append(l.bufs, buf)
	}
	l.Unlock()
}

func (l *Logger) Close() {
	select {
	case l.quit <- struct{}{}:
	case <-l.quitDone:
		return
	}

	<-l.quitDone // First time close
	l.handler.Close()
}

// set log level, any log level less than it will not log
func (l *Logger) SetLevel(level int) {
	l.level = level
}

// get log level
func (l *Logger) GetLevel() int {
	return l.level
}

// a low interface, maybe you can use it for your special log format
// but it may be not exported later......
func (l *Logger) Output(callDepth int, level int, format string, v ...interface{}) {
	if l.level > level {
		return
	}

	s := fmt.Sprintf(format, v...)
	l.output(callDepth, level, s)
}

func (l *Logger) Outputln(callDepth int, level int, v ...interface{}) {
	if l.level > level {
		return
	}
	s := fmt.Sprintln(v...)
	l.output(callDepth, level, s)
}

func (l *Logger) output(callDepth int, level int, s string) {

	buf := l.popBuf()

	if l.flag&Ltime > 0 {
		now := time.Now().Format(TimeFormat)
		buf = append(buf, '[')
		buf = append(buf, now...)
		buf = append(buf, "] "...)
	}

	if l.flag&Lfile > 0 {
		_, file, line, ok := runtime.Caller(callDepth)
		if !ok {
			file = "???"
			line = 0
		} else {
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					file = file[i+1:]
					break
				}
			}
		}

		buf = append(buf, file...)
		buf = append(buf, ':')

		buf = strconv.AppendInt(buf, int64(line), 10)
		buf = append(buf, ' ')
	}

	if l.flag&Llevel > 0 {
		buf = append(buf, '[')
		buf = append(buf, LevelName[level]...)
		buf = append(buf, "] [-] "...)
	}

	buf = append(buf, s...)

	if s[len(s)-1] != '\n' {
		buf = append(buf, '\n')
	}

	l.msg <- buf
}

// convert struct param to json string
func (l *Logger) OutputJson(level int, v interface{}) {
	if l.level > level {
		return
	}

	buf := l.popBuf()

	reqBody, err := json.Marshal(v)
	if err != nil {
		fmt.Printf("json marshal err(%s)\n", err.Error())
		return
	}

	buf = append(buf, reqBody...)

	if buf[len(buf)-1] != '\n' {
		buf = append(buf, '\n')
	}

	l.msg <- buf
}

func (l *Logger) Trace(format string, v ...interface{}) {
}

func (l *Logger) TraceJson(v interface{}) {
}

func (l *Logger) Debug(format string, v ...interface{}) {
}

func (l *Logger) DebugJson(v interface{}) {
}

func (l *Logger) LogDebug(v ...interface{}) {
}

func (l *Logger) LogInfo(v ...interface{}) {
	l.Outputln(3, LevelInfo, v...)
}

func (l *Logger) LogTrace(v ...interface{}) {
}

func (l *Logger) LogWarn(v ...interface{}) {
	l.Outputln(3, LevelWarn, v...)
}

func (l *Logger) LogError(v ...interface{}) {
	l.Outputln(3, LevelError, v...)
}

func (l *Logger) LogFatal(v ...interface{}) {
	l.Outputln(3, LevelFatal, v...)
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.Output(3, LevelInfo, format, v...)
}

func (l *Logger) InfoJson(v interface{}) {
	l.OutputJson(LevelInfo, v)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.Output(3, LevelWarn, format, v...)
}

func (l *Logger) WarnJson(v interface{}) {
	l.OutputJson(LevelWarn, v)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.Output(3, LevelError, format, v...)
}

func (l *Logger) ErrorJson(v interface{}) {
	l.OutputJson(LevelError, v)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.Output(3, LevelFatal, format, v...)
}

func (l *Logger) FatalJson(v interface{}) {
	l.OutputJson(LevelFatal, v)
}

func SetLevel(level int) {
	std.SetLevel(level)
}

func Trace(format string, v ...interface{}) {
}

func LogTrace(v ...interface{}) {
}

func Debug(format string, v ...interface{}) {
}

func Info(format string, v ...interface{}) {
	std.Output(3, LevelInfo, format, v...)
}

func Warn(format string, v ...interface{}) {
	std.Output(3, LevelWarn, format, v...)
}

func Error(format string, v ...interface{}) {
	std.Output(3, LevelError, format, v...)
}

func Fatal(format string, v ...interface{}) {
	std.Output(3, LevelFatal, format, v...)
}

func (l *Logger) OutputContext(calldepth int, ctx context.Context, level int, format string, v ...interface{}) {
	if l.level > level {
		return
	}

	buf := l.popBuf()
	if l.flag&Ltime > 0 {
		now := time.Now().Format(TimeFormat)
		buf = append(buf, '[')
		buf = append(buf, now...)
		buf = append(buf, "] "...)
	}

	if l.flag&Lfile > 0 {
		_, file, line, ok := runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		} else {
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					file = file[i+1:]
					break
				}
			}
		}

		buf = append(buf, file...)
		buf = append(buf, ':')

		buf = strconv.AppendInt(buf, int64(line), 10)
		buf = append(buf, ' ')
	}

	if l.flag&Llevel > 0 {
		buf = append(buf, '[')
		buf = append(buf, LevelName[level]...)
		buf = append(buf, "] "...)
	}
	traceId := "-"
	if ctx != nil {
		traceId, _ = ctx.Value("Trace_id").(string)
	}

	traceIdArea := fmt.Sprintf("[%s] ", traceId)
	buf = append(buf, traceIdArea...)

	s := fmt.Sprintf(format, v...)
	buf = append(buf, s...)

	if s[len(s)-1] != '\n' {
		buf = append(buf, '\n')
	}

	l.msg <- buf
}

func (l *Logger) TraceContext(ctx context.Context, format string, v ...interface{}) {
}

func (l *Logger) DebugContext(ctx context.Context, format string, v ...interface{}) {
}

func (l *Logger) WarnContext(ctx context.Context, format string, v ...interface{}) {
	l.OutputContext(2, ctx, LevelWarn, format, v...)
}

func (l *Logger) InfoContext(ctx context.Context, format string, v ...interface{}) {
	l.OutputContext(2, ctx, LevelInfo, format, v...)
}

func (l *Logger) ErrorContext(ctx context.Context, format string, v ...interface{}) {
	l.OutputContext(2, ctx, LevelError, format, v...)
}

func (l *Logger) FatalContext(ctx context.Context, format string, v ...interface{}) {
	l.OutputContext(2, ctx, LevelFatal, format, v...)
}
