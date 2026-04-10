package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"

	"errors"
	_ "net/http/pprof"
	"strconv"

	"financeMgr/src/common/config"
	"financeMgr/src/common/log"
	"financeMgr/src/common/tag"
	"financeMgr/src/common/url"
	"financeMgr/src/common/utils"
	"os"
	"os/signal"
	"syscall"

	"financeMgr/src/analysis-server/api/cfg"
	"financeMgr/src/analysis-server/api/db"
	"financeMgr/src/analysis-server/api/handler"
	"financeMgr/src/analysis-server/api/service"
)

var (
	gLogger *log.Logger
)

func interceptSignal(ws *sync.WaitGroup) {
	daemonExitCh := make(chan os.Signal, 1)
	signal.Notify(daemonExitCh, syscall.SIGTERM, syscall.SIGQUIT,
		syscall.SIGINT, syscall.SIGHUP)
	for {
		sig := <-daemonExitCh
		close(daemonExitCh)
		gLogger.LogInfo("receive signal: ", sig.String())
		service.StopIdResourcePersistence()
		handler.GAccessTokenH.QuitExpirationCheckService()
		break
	}
	ws.Done()
}

func startServer(router *url.UrlRouter, serverConf *cfg.ServerConf) {
	//runtime.GOMAXPROCS(serverConf.Cores)
	http.Handle(serverConf.BaseUrl, router)
	go func() {
		if err := http.ListenAndServe(":"+strconv.Itoa(serverConf.Port), nil); err != nil {
			fmt.Println("[Main] http server exit, error: ", err)
		}
	}()
}

func startBusinessWorker(ws *sync.WaitGroup, syncDuration int) {
	ws.Add(1)
	go interceptSignal(ws)
	//用户登录的过期检查服务
	ws.Add(1)
	go handler.GAccessTokenH.ExpirationCheck(ws)
	ws.Add(1)
	go service.StartIdResourcePersistence(time.Duration(syncDuration)*time.Minute, ws)
}

func releaseBusinessResource() {
	service.ReleaseResourceInService()
}

func registerHandler(httpRouter *url.UrlRouter, logger *log.Logger, dbIns *sql.DB) error {
	//下面的几个变量是公共使用的部分
	/* Dao */
	db.InitDao(logger)
	companyDao := &db.CompanyDao{}
	companygroupDao := &db.CompanyGroupDao{}
	voucherRecordDao := &db.VoucherRecordDao{}
	/*service*/
	//初始化API service
	ccErr := service.InitService(logger, dbIns)
	if ccErr != nil {
		return errors.New("InitService err: " + ccErr.Detail())
	}
	comService := &service.CompanyService{
		CompanyDao:      companyDao,
		CompanyGroupDao: companygroupDao}

	//注册handler
	handler.InitHandler(logger)
	registerYearBalance(httpRouter)
	registerVoucherTemplate(httpRouter)
	registerComGroup(httpRouter, companygroupDao)
	registerCompany(httpRouter, comService)
	registerAccSub(httpRouter, companyDao, voucherRecordDao)
	registerOptAndAuthenHandler(httpRouter, comService)
	registerResAndVoucherHandler(httpRouter, companyDao, voucherRecordDao)
	registerMenuHandler(httpRouter)
	return nil
}

func main() {

	if utils.SetLimit() != nil {
		fmt.Println("[Main] set max open files failed")
		return
	}

	// all in one
	var apiServerCfgFile = flag.String("c",
		"/etc/analysis/web_server.cfg", "Server config file name")

	flag.Parse()
	if tag.CheckAndShowVersion() {
		return
	}

	apiServerConf, err := cfg.ParseApiServerConfig(apiServerCfgFile)
	if err != nil {
		fmt.Println("[Main] parse config", *apiServerCfgFile, "err: ", err)
		return
	}
	if err = apiServerConf.CheckValid(); err != nil {
		fmt.Println("[init] checkconfig err", err)
		return
	}

	gLogger, err = config.LogFac{Logconf: apiServerConf.LogConf}.NewLogger()
	if err != nil {
		fmt.Println("[Main] new logger err: ", err)
		return
	}
	//initialize the common url router
	url.InitCommonUrlRouter(gLogger, nil)
	httpRouter := url.NewUrlRouter(gLogger)
	//initialize the database connection
	dbIns, err := config.MysqlInstance{Conf: apiServerConf.MysqlConf}.NewMysqlInstance(gLogger)
	if err != nil {
		fmt.Println("[Main] Create Db connection error: ", err)
		return
	}
	//register the handle
	err = registerHandler(httpRouter, gLogger, dbIns)
	if err != nil {
		fmt.Println("[Main] Handler register error: ", err)
		return
	}
	ws := &sync.WaitGroup{}
	startBusinessWorker(ws, apiServerConf.ServerConf.SynDuration)
	//start server
	startServer(httpRouter, apiServerConf.ServerConf)
	ws.Wait()
	gLogger.LogInfo("[Main] analysis server is exiting...")
	//release resource
	releaseBusinessResource()
	dbIns.Close()
	gLogger.Close()
	fmt.Println("[Main] analysis server exit")
}
