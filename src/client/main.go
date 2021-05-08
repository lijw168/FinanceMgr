package main

import (
	"client/service"
	"common/log"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func initLogger(fileName string, fileMaxSize, filCount int) (*log.Logger, error) {
	var h log.Handler
	var err error
	h, err = log.NewRotatingFileHandler(fileName, fileMaxSize, filCount)
	if err != nil {
		fmt.Printf("new log handler err: %v\r\n", err.Error())
		return nil, err
	}
	logger := log.NewDefault(h)
	//logger.SetLevel(log.LevelInfo)
	logger.SetLevel(log.LevelDebug)
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
	var pStrServerHost = flag.String("h", "127.0.0.1", "gateway server host")
	var pTimeout = flag.Uint64("p", 30, "timeout")
	//tcp server information
	var pListenPort = flag.Int("l", 9999, "tcp listen port")
	flag.Parse()

	// if flag.NFlag() < 3 {
	// 	flag.Usage()
	// 	return
	// }
	logger, err := initLogger(*logFileName, *logFileSize, *logFileCount)
	if err != nil {
		fmt.Printf("Init logger err: %v \r\n", err.Error())
		return
	}
	logger.LogInfo("login service is beginning")
	proxy := service.Proxy{}
	proxy.Init(*pListenPort, *pServerPort, *pStrServerHost, *pTimeout, logger)
	proxy.StartTcpService()
	logger.LogInfo("login quit")
}
