package url

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"encoding/json"
	cons "financeMgr/src/common/constant"
	"financeMgr/src/common/log"
	"financeMgr/src/common/message"
	"financeMgr/src/common/utils"
)

func NewUrlRouter(l *log.Logger) *UrlRouter {
	return &UrlRouter{
		Logger:  l,
		Handler: make(map[string]http.Handler),
		HF:      make(map[string]http.HandlerFunc),
	}
}

type LoginCheckHandle func(action string, r *http.Request) (bool, string, error)
type InterfaceAuthorityCheckHandle func(action, accessToken string) (bool, error)

type UrlRouter struct {
	Handler                 map[string]http.Handler
	HF                      map[string]http.HandlerFunc
	Logger                  *log.Logger
	hostname                string
	LoginCheck              LoginCheckHandle
	InterfaceAuthorityCheck InterfaceAuthorityCheckHandle
}

func (p *UrlRouter) RegisterFunc(action string, handler http.HandlerFunc) *UrlRouter {
	p.HF[action] = handler
	if len(p.hostname) == 0 {
		p.hostname, _ = os.Hostname()
	}
	return p
}

func (p *UrlRouter) Register(action string, handler http.Handler) *UrlRouter {
	p.Handler[action] = handler
	if len(p.hostname) == 0 {
		p.hostname, _ = os.Hostname()
	}
	return p
}

func (p *UrlRouter) checkRequest(action string, r *http.Request) (bool, message.ResponseParam) {
	rsp := message.ResponseParam{}
	//login check
	if p.LoginCheck != nil {
		bIsPass, accessToken, err := p.LoginCheck(action, r)
		if !bIsPass {
			if err == nil {
				rsp.Code = cons.CodeNoLogin
				rsp.Message = "please login first."
			} else {
				rsp.Code = -1
				rsp.Message = err.Error()
			}
			return bIsPass, rsp
		}
		//api interface authority check
		if p.InterfaceAuthorityCheck != nil {
			bIsPass, err = p.InterfaceAuthorityCheck(action, accessToken)
			if !bIsPass {
				rsp.Code = cons.CodeNoAuthority
				rsp.Message = err.Error()
				return bIsPass, rsp
			}
		}
	}
	return true, rsp
}

func (p *UrlRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			p.Logger.ErrorContext(r.Context(), "[url/handler] [method: %v, url: %v, remote_addr:%v, panic: %v, x-forwarded-for: %v, stack %s]", r.Method,
				r.URL.RequestURI(), r.RemoteAddr, err, r.Header.Get("X-Forwarded-For"), string(debug.Stack()))
			panic(err)
		}
	}()

	// 最新的是使用open-api的方式
	// Client-Token 幂等
	// Request-Id
	// 中间层请求是Request-Id
	// 主机过来的Trace-Id
	traceIds, ok := r.Header["Client-Token"]
	if !ok {
		if traceIds, ok = r.Header["Request-Id"]; !ok {
			traceIds, _ = r.Header["Trace-Id"]
		}
	}

	var traceId string
	if len(traceIds) < 1 || traceIds[0] == "" {
		traceId = utils.Uuid()
		p.Logger.LogTrace("[url/handler] [generate traceid]", traceId)
		r.Header.Set("Trace-Id", traceId)
		w.Header().Set("Trace-Id", traceId)
	} else {
		traceId = traceIds[0]
	}
	// add trace_id
	ctx := r.Context()
	ctx = context.WithValue(ctx, "Trace-Id", traceId)
	r = r.WithContext(ctx)
	start := time.Now()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Requst-Trace", p.hostname)
	action := r.FormValue("Action")
	if len(action) == 0 {
		action = r.FormValue("action")
	}
	if action == "" {
		p.Logger.ErrorContext(r.Context(), "[url/handler] [method: %v, url: %v, remote_addr:%v, action: %v, x-forwarded-for: %v]", r.Method,
			r.URL.RequestURI(), r.RemoteAddr, action, r.Header.Get("X-Forwarded-For"))
		w.Write([]byte(message.ErrReqParam))
		return
	}

	bIsPass, rsp := p.checkRequest(action, r)
	if !bIsPass {
		jsonRsp, err := json.Marshal(&rsp)
		if err != nil {
			p.Logger.ErrorContext(ctx, "[UrlRouter/ServeHTTP] [Marshal,failed, error: %s]", err.Error())
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonRsp)
		return
	}

	var Ihandler http.Handler
	var IhandlerFunc http.HandlerFunc
	if p.Handler != nil {
		Ihandler, _ = p.Handler[action]
	}
	if Ihandler == nil && p.HF != nil {
		IhandlerFunc, _ = p.HF[action]
	}

	if IhandlerFunc == nil && Ihandler == nil {
		p.Logger.InfoContext(r.Context(), "[url/handler] [method: %v, url: %v, remote_addr:%v, action: %v, x-forwarded-for: %v]", r.Method,
			r.URL.RequestURI(), r.RemoteAddr, action, r.Header.Get("X-Forwarded-For"))
		w.Write([]byte(message.ErrNoAction))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		p.Logger.ErrorContext(r.Context(), "[url/handler] [method: %v, url: %v, remote_addr:%v, action: %v, x-forwarded-for: %v]", r.Method,
			r.URL.RequestURI(), r.RemoteAddr, action, r.Header.Get("X-Forwarded-For"))
		return
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	p.Logger.InfoContext(r.Context(), "[url/handler] [method: %v, url: %v, body: %v, remote_addr:%v, action: %v, header : %v ]", r.Method,
		r.URL.RequestURI(), strings.Replace(string(body), "\n", " ", -1), r.RemoteAddr, action, r.Header)

	if IhandlerFunc != nil {
		IhandlerFunc(w, r)
	} else if Ihandler != nil {
		Ihandler.ServeHTTP(w, r)
	} else {
		p.Logger.ErrorContext(r.Context(), "[url/handler]")
		return
	}
	elapsed := time.Since(start)
	p.Logger.InfoContext(r.Context(), "[url/handler] elapsed: %v]", elapsed)
}

func (p *UrlRouter) SetLogLevel(w http.ResponseWriter, r *http.Request) {
	level := r.FormValue("level")
	if level == "" {
		p.Logger.ErrorContext(r.Context(), "[url/handler] [level: empty]")
		w.Write([]byte(message.ErrReqParam))
		return
	}

	var logLevel int
	switch strings.ToLower(level) {
	case "trace":
		logLevel = log.LevelTrace
	case "debug":
		logLevel = log.LevelDebug
	case "info":
		logLevel = log.LevelInfo
	case "warn":
		logLevel = log.LevelWarn
	case "error":
		logLevel = log.LevelError
	case "fatal":
		logLevel = log.LevelFatal
	default:
		p.Logger.ErrorContext(r.Context(), "[url/handler] [level: %v]", level)
		w.Write([]byte(message.ErrReqParam))
		return
	}
	p.Logger.InfoContext(r.Context(), "[url/handler] [newLevel: %v]", level)
	p.Logger.SetLevel(logLevel)
}

func (p *UrlRouter) GetLogLevel(w http.ResponseWriter, r *http.Request) {
	level := p.Logger.GetLevel()

	var logLevel string
	switch level {
	case log.LevelTrace:
		logLevel = "trace"
	case log.LevelDebug:
		logLevel = "debug"
	case log.LevelInfo:
		logLevel = "info"
	case log.LevelWarn:
		logLevel = "warn"
	case log.LevelError:
		logLevel = "error"
	case log.LevelFatal:
		logLevel = "fatal"
	default:
		p.Logger.ErrorContext(r.Context(), "[unknown loglevel: %v]", level)
		w.Write([]byte(message.ErrReqParam))
		return
	}
	w.Write([]byte(logLevel))
}
