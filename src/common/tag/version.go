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
	ZBS_BUILD_VERSION  string // commit
	ZBS_BUILD_TIME     string
	GO_VERSION         string // go version
	WEB_SERVER_VERSION string

	// ZBS_CLIENT_VER    string
	// ZBS_COMMON_VER    string
	// ZBS_GATEWAY_VER   string
	// ZBS_SCHEDULER_VER string
	// ZBS_SERVER_VER    string
	// ZBS_STORAGE_VER   string
	// ZBS_WORKER_VER    string
	// ZBS_OPENAPI_VER   string

	showVersion = flag.Bool("v", false, "show build version and time")
)

type VersionInfo struct {
	BuildTime        string
	GoVersion        string
	BuildVersion     string
	WebServerVersion string
	// ZbsClientVersion    string
	// ZbsCommonVersion    string
	// ZbsGatewayVersion   string
	// ZbsSchedulerVersion string
	// ZbsServerVersion    string
	// ZbsStorageVersion   string
	// ZbsWorkerVersion    string
	// ZbsOpenApiVersion   string
}

func CheckAndShowVersion() bool {
	if *showVersion {
		log.Println("BuildTime\t", ZBS_BUILD_TIME)
		log.Println("GoVersion\t", GO_VERSION)
		log.Println("BuildVersion\t", ZBS_BUILD_VERSION)
		log.Println("BuildVersion\t", WEB_SERVER_VERSION)
		// log.Println("ZbsClientVersion\t", ZBS_CLIENT_VER)
		// log.Println("ZbsCommonVersion\t", ZBS_COMMON_VER)
		// log.Println("ZbsGatewayVersion\t", ZBS_GATEWAY_VER)
		// log.Println("ZbsSchedulerVersion\t", ZBS_SCHEDULER_VER)
		// log.Println("ZbsServerVersion\t", ZBS_SERVER_VER)
		// log.Println("ZbsStorageVersion\t", ZBS_STORAGE_VER)
		// log.Println("ZbsWorkerVersion\t", ZBS_WORKER_VER)
		// log.Println("ZbsOpenApiVersion\t", ZBS_OPENAPI_VER)
		return true
	}
	return false
}

const (
	ZBS_WEB_VERSION_PATH = "/version"
)

func ShowVersionHandler(w http.ResponseWriter, r *http.Request) {
	var buf []byte
	var err error
	info := VersionInfo{
		BuildVersion:     ZBS_BUILD_VERSION,
		BuildTime:        ZBS_BUILD_TIME,
		GoVersion:        GO_VERSION,
		WebServerVersion: WEB_SERVER_VERSION,
		// ZbsClientVersion:    ZBS_CLIENT_VER,
		// ZbsCommonVersion:    ZBS_COMMON_VER,
		// ZbsGatewayVersion:   ZBS_GATEWAY_VER,
		// ZbsSchedulerVersion: ZBS_SCHEDULER_VER,
		// ZbsServerVersion:    ZBS_SERVER_VER,
		// ZbsStorageVersion:   ZBS_STORAGE_VER,
		// ZbsWorkerVersion:    ZBS_WORKER_VER,
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
	url := fmt.Sprintf("%s%s", host, ZBS_WEB_VERSION_PATH)
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
