package url

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"financeMgr/src/common/log"
	"financeMgr/src/common/utils"

	"github.com/gorilla/mux"
)

type RestUrlRouter struct {
	mux      *mux.Router
	Logger   *log.Logger
	hostname string
}

func (p *RestUrlRouter) Init(path string) *RestUrlRouter {
	if p.mux == nil {
		p.mux = mux.NewRouter().StrictSlash(true)
	}
	if len(p.hostname) == 0 {
		p.hostname, _ = os.Hostname()
	}
	if strings.Index(path, "/") == 0 && strings.LastIndex(path, "/")+1 == len(path) {
		p.mux = p.mux.PathPrefix(path).Subrouter()
	}
	return p
}

func (p *RestUrlRouter) AddOjbect(path string) *mux.Router {
	return p.mux.PathPrefix(path).Subrouter()
}

func (p *RestUrlRouter) GetRouter() *mux.Router {
	return p.mux
}

func (p *RestUrlRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			p.Logger.ErrorContext(r.Context(), "[RestUrlRouter/handler] [method: %v, url: %v, remote_addr:%v, panic: %v, x-forwarded-for: %v]", r.Method,
				r.URL.RequestURI(), r.RemoteAddr, err, r.Header.Get("X-Forwarded-For"))
			panic(err)
		}
	}()
	requestId := r.Header.Get("x-jcloud-request-id")
	if requestId == "" {
		requestId = utils.Uuid()
	}
	w.Header().Set("Trace-Id", requestId)
	// add Trace-Id
	ctx := r.Context()
	ctx = context.WithValue(ctx, "Trace-Id", requestId)
	r = r.WithContext(ctx)
	start := time.Now()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Requst-Host-Trace", p.hostname)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		p.Logger.ErrorContext(r.Context(), "[url/handler] [method: %v, url: %v, remote_addr:%v, x-forwarded-for: %v]", r.Method,
			r.URL.RequestURI(), r.RemoteAddr, r.Header.Get("X-Forwarded-For"))
		http.NotFound(w, r)
		return
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	defer func(elapsed time.Duration) {
		p.Logger.InfoContext(r.Context(), "[RestUrlRouter/handler] [method: %v, url: %v, body: %v, remote_addr:%v, x-forwarded-for: %v, elapsed: %v]", r.Method,
			r.URL.RequestURI(), strings.Replace(string(body), "\n", " ", -1), r.RemoteAddr, r.Header.Get("X-Forwarded-For"), elapsed)
	}(time.Since(start))

	p.mux.NotFoundHandler = UrlUnMatchHandler()
	p.mux.ServeHTTP(w, r)
}

var UrlUmatchErrStr = "{\"requestId\": \"%s\",\"error\":{\"code\":%d,\"message\":\"url not found.\",\"status\":\"NOT_FOUND\",\"detail\":[]}}"

func UrlUnFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	requestId := ""
	ctx := r.Context()
	if ctx.Value("Trace-Id") != nil {
		requestId = ctx.Value("Trace-Id").(string)
	}
	w.Write([]byte(fmt.Sprintf(UrlUmatchErrStr, requestId, http.StatusNotFound)))
}

func UrlUnMatchHandler() http.Handler { return http.HandlerFunc(UrlUnFound) }
