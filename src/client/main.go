package main

import (
	"financeMgr/src/client/service"
	"financeMgr/src/common/log"
	"flag"
	"fmt"
	_ "net/http/pprof"
	"os"
	"path/filepath"
)

func initLogger(fileName string, fileMaxSize, filCount, iLogLevel int) (*log.Logger, error) {
	var h log.Handler
	var err error
	h, err = log.NewRotatingFileHandler(fileName, fileMaxSize, filCount)
	if err != nil {
		fmt.Printf("new log handler err: %v\r\n", err.Error())
		return nil, err
	}
	logger := log.NewDefault(h)
	logger.SetLevel(iLogLevel)
	return logger, nil
}

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("get file path,failed,err:%v\r\n", err.Error())
		return
	}
	//log information
	strDefaultLogFile := dir + ".log"
	var logFileName = flag.String("f", strDefaultLogFile, "log file name")
	var logFileCount = flag.Int("n", 20, "number of log files")
	var logFileSize = flag.Int("s", 20480, "log file size")
	//server information
	var pServerPort = flag.Int("p", 7500, "gateway server port")
	var pServerHost = flag.String("a", "47.100.210.38", "gateway server host")
	//var pServerHost = flag.String("a", "192.168.148.190", "gateway server host")
	var pTimeout = flag.Uint64("t", 3000, "timeout")
	//tcp server information
	var pListenPort = flag.Int("l", 9999, "tcp listen port")
	flag.Parse()

	// if flag.NFlag() < 3 {
	// 	flag.Usage()
	// 	return
	// }
	logger, err := initLogger(*logFileName, *logFileSize, *logFileCount, log.LevelDebug)
	if err != nil {
		fmt.Printf("Init logger err: %v \r\n", err.Error())
		return
	}
	logger.LogInfo("proxy service is beginning")
	proxy := service.Proxy{}
	proxy.Init(*pListenPort, *pServerPort, *pServerHost, *pTimeout, logger)
	proxy.StartTcpService()
	logger.LogInfo("the proxy service has been to end")
}
