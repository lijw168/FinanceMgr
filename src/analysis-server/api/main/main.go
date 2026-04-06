package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"time"

	_ "net/http/pprof"
	//"runtime"
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
	//aUtils "financeMgr/src/analysis-server/api/utils"
)

var (
	exitCh  = make(chan bool)
	gLogger *log.Logger
)

func interceptSignal() {
	daemonExitCh := make(chan os.Signal, 1)
	signal.Notify(daemonExitCh, syscall.SIGTERM, syscall.SIGQUIT,
		syscall.SIGINT, syscall.SIGHUP)
	go func() {
		for {
			sig := <-daemonExitCh
			gLogger.LogInfo("receive signal: ", sig.String())
			service.StopIdResourcePersistence()
			handler.GAccessTokenH.QuitExpirationCheckService()
			break
		}
		exitCh <- true
	}()
}

func waitDaemonExit() {
	<-exitCh
	time.Sleep(3 * time.Second)
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

func handlerInit(httpRouter *url.UrlRouter, logger *log.Logger, dbIns *sql.DB) error {
	//var err error
	// if serverConf.IsUserApiServer() {
	// 	err = initUserApiServer(serverConf.UserServerCfg, logger, httpRouter, copySnpCfg)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	logger.LogInfo("init user api server")
	// }
	err := initApiServer(dbIns, logger, httpRouter)
	if err != nil {
		return err
	}
	logger.LogInfo("init api server")
	return nil
}

func initApiServer(dbIns *sql.DB, logger *log.Logger, httpRouter *url.UrlRouter) error {
	//初始化ID Resource
	idInfoDao := &db.IDInfoDao{Logger: logger}
	service.GIdInfoService.InitIdInfoService(logger, idInfoDao, dbIns)
	ccErr := service.GIdInfoService.InitIdResource()
	if ccErr != nil {
		return ccErr
	}
	//下面的几个变量是公共使用的部分
	/* Dao */
	companyDao := &db.CompanyDao{Logger: logger}
	companygroupDao := &db.CompanyGroupDao{Logger: logger}
	voucherRecordDao := &db.VoucherRecordDao{Logger: logger}
	/*service*/
	comService := &service.CompanyService{
		Logger:          logger,
		CompanyDao:      companyDao,
		CompanyGroupDao: companygroupDao,
		Db:              dbIns}

	registerYearBalance(logger, httpRouter, dbIns)
	registerVoucherTemplate(logger, httpRouter, dbIns)
	registerComGroup(logger, httpRouter, companygroupDao, dbIns)
	registerCompany(logger, httpRouter, comService)
	registerAccSub(logger, httpRouter, companyDao, voucherRecordDao, dbIns)
	registerOptAndAuthenHandler(logger, httpRouter, comService, dbIns)
	registerResAndVoucherHandler(logger, httpRouter, companyDao, voucherRecordDao, dbIns)
	registerMenuHandler(logger, httpRouter, dbIns)
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
	err = handlerInit(httpRouter, gLogger, dbIns)
	if err != nil {
		fmt.Println("[Main] Handler register error: ", err)
		return
	}
	interceptSignal()
	service.StartIdResourcePersistence(time.Duration(apiServerConf.ServerConf.SynDuration) * time.Minute)
	//start server
	startServer(httpRouter, apiServerConf.ServerConf)
	waitDaemonExit()
	gLogger.Close()
	dbIns.Close()
	fmt.Println("[Main] analysis server exit")
}
