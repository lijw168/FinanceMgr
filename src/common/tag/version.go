package tag

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	FINANCE_BUILD_VERSION string // commit
	FINANCE_BUILD_TIME    string
	GO_VERSION            string // go version
	WEB_SERVER_VERSION    string
	showVersion           = flag.Bool("v", false, "show build version and time")
)

type VersionInfo struct {
	BuildTime        string
	GoVersion        string
	BuildVersion     string
	WebServerVersion string
}

func CheckAndShowVersion() bool {
	if *showVersion {
		log.Println("BuildTime\t", FINANCE_BUILD_TIME)
		log.Println("GoVersion\t", GO_VERSION)
		log.Println("BuildVersion\t", FINANCE_BUILD_VERSION)
		log.Println("analysis-server\t", WEB_SERVER_VERSION)
		return true
	}
	return false
}

const (
	VERSION_PATH = "/version"
)

func ShowVersionHandler(w http.ResponseWriter, r *http.Request) {
	var buf []byte
	var err error
	info := VersionInfo{
		BuildVersion:     FINANCE_BUILD_VERSION,
		BuildTime:        FINANCE_BUILD_TIME,
		GoVersion:        GO_VERSION,
		WebServerVersion: WEB_SERVER_VERSION,
	}
	if buf, err = json.Marshal(info); err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(buf)
	}
}

func QueryVersion(host string) (*VersionInfo, error) {
	var err error
	var buf []byte
	var rsp *http.Response
	url := fmt.Sprintf("%s%s", host, VERSION_PATH)
	if rsp, err = http.Get(url); err != nil {
		return nil, err
	}

	if rsp.Body == nil {
		return nil, errors.New("empty response")
	}

	if buf, err = ioutil.ReadAll(rsp.Body); err != nil {
		return nil, err
	}

	defer rsp.Body.Close()

	var info VersionInfo
	if err = json.Unmarshal(buf, &info); err != nil {
		return nil, err
	}
	return &info, nil
}
