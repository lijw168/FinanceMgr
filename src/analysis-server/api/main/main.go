package main

import (
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
	exitCh = make(chan bool)
)

func interceptSignal() {
	daemonExitCh := make(chan os.Signal, 1)
	signal.Notify(daemonExitCh, syscall.SIGTERM, syscall.SIGQUIT,
		syscall.SIGINT, syscall.SIGHUP)
	go func() {
		for {
			sig := <-daemonExitCh
			fmt.Printf("time:%v;the sig is %s\n", time.Now(), sig.String())
			saveIdResource()
			handler.GAccessTokenH.QuitExpirationCheckService()
			break
		}
		exitCh <- true
	}()
}

func saveIdResource() {
	ccErr := service.GIdInfoService.WriteIdResourceToDb()
	if ccErr != nil {
		ccErr := service.GIdInfoService.WriteIdResourceToDb()
		if ccErr != nil {
			fmt.Printf("time:%v;WriteIdResourceToDb,it is twice to fail,ErrInfo:%s", time.Now(), ccErr.Error())
		}
	}
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
			fmt.Println("[Init] http server exit, error: ", err)
		}
	}()
}

func handlerInit(httpRouter *url.UrlRouter, logger *log.Logger, apiServerConf *cfg.ApiServerConf) error {
	//var err error
	// if serverConf.IsUserApiServer() {
	// 	err = initUserApiServer(serverConf.UserServerCfg, logger, httpRouter, copySnpCfg)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	logger.LogInfo("init user api server")
	// }
	err := initApiServer(apiServerConf.MysqlConf, logger, httpRouter)
	if err != nil {
		return err
	}
	logger.LogInfo("init api server")
	return nil
}

func initApiServer(mysqlConf *config.MysqlConf, logger *log.Logger, httpRouter *url.UrlRouter) error {
	_db, err := config.MysqlInstance{Conf: mysqlConf, Logger: logger}.NewMysqlInstance()
	if err != nil {
		fmt.Println("[Init] Create Db connection error: ", err)
		return err
	}
	//初始化ID Resource
	idInfoDao := &db.IDInfoDao{Logger: logger}
	service.GIdInfoService.InitIdInfoService(logger, idInfoDao, _db)
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
		Db:              _db}

	registerYearBalance(logger, httpRouter, _db)
	registerVoucherTemplate(logger, httpRouter, _db)
	registerComGroup(logger, httpRouter, companygroupDao, _db)
	registerCompany(logger, httpRouter, comService)
	registerAccSub(logger, httpRouter, companyDao, voucherRecordDao, _db)
	registerOptAndAuthenHandler(logger, httpRouter, comService, _db)
	registerResAndVoucherHandler(logger, httpRouter, companyDao, voucherRecordDao, _db)
	registerMenuHandler(logger, httpRouter, _db)
	return nil
}

func main() {

	if utils.SetLimit() != nil {
		fmt.Println("[Init] set max open files failed")
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
		fmt.Println("[Init] parse config", *apiServerCfgFile, "err: ", err)
		return
	}
	if err = apiServerConf.CheckValid(); err != nil {
		fmt.Println("[init] checkconfig err", err)
		return
	}

	logger, err := config.LogFac{Logconf: apiServerConf.LogConf}.NewLogger()
	if err != nil {
		fmt.Println("[Init] new logger err: ", err)
		return
	}
	url.InitCommonUrlRouter(logger, nil)
	httpRouter := url.NewUrlRouter(logger)
	err = handlerInit(httpRouter, logger, apiServerConf)
	if err != nil {
		fmt.Println("[Init] Handler register error: ", err)
		return
	}
	interceptSignal()
	startServer(httpRouter, apiServerConf.ServerConf)
	waitDaemonExit()
	logger.Close()
}
