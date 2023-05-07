package url

import (
	"net/http"

	_ "net/http/pprof"

	"financeMgr/src/common/log"
	"financeMgr/src/common/tag"
)

func InitCommonUrlRouter(l *log.Logger, funcs map[string]http.HandlerFunc) {
	// http://host:port/common?action=
	// setloglevel&level=[trace|debug|info|warn|error|fatal]
	// version
	// getloglevel
	router := NewUrlRouter(l)
	router.RegisterFunc("setloglevel", router.SetLogLevel)
	router.RegisterFunc("getloglevel", router.GetLogLevel)
	router.RegisterFunc("version", tag.ShowVersionHandler)

	for key, funcPtr := range funcs {
		router.RegisterFunc(key, funcPtr)
	}

	http.Handle("/common", router)
}
