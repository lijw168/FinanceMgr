package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	//_ "net/http/pprof"
	"runtime"
	"strconv"

	"common/config"
	"common/log"
	"common/tag"
	"common/url"
	"common/utils"
	"os"
	"os/signal"
	"syscall"

	"analysis-server/api/cfg"
	"analysis-server/api/db"
	"analysis-server/api/handler"
	"analysis-server/api/service"
	//aUtils "analysis-server/api/utils"
)

var (
	exitCh = make(chan bool)
)

func interceptSignal() {
	daemonExitCh := make(chan os.Signal)
	signal.Notify(daemonExitCh, syscall.SIGTERM, syscall.SIGKILL,
		syscall.SIGQUIT, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		for {
			sig := <-daemonExitCh
			fmt.Printf("the sig is %s\n", sig.String())
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
			fmt.Printf("WriteIdResourceToDb,it is twice to fail,ErrInfo:%s", ccErr.Error())
		}
	}
}

func waitDaemonExit() {
	<-exitCh
	time.Sleep(3 * time.Second)
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
		fmt.Println("[Init] Handler registe error: ", err)
		return
	}
	interceptSignal()
	startServer(httpRouter, apiServerConf.ServerConf)
	waitDaemonExit()
	logger.Close()
	return
}

func startServer(router *url.UrlRouter, serverConf *cfg.ServerConf) {

	runtime.GOMAXPROCS(serverConf.Cores)

	http.Handle(serverConf.BaseUrl, router)
	go func() {
		if err := http.ListenAndServe(":"+strconv.Itoa(serverConf.Port), nil); err != nil {
			fmt.Println("[Init] http server exit, error: ", err)
		}
	}()
	return
}

func handlerInit(httpRouter *url.UrlRouter, logger *log.Logger, apiServerConf *cfg.ApiServerConf) error {
	var err error
	// if serverConf.IsUserApiServer() {
	// 	err = initUserApiServer(serverConf.UserServerCfg, logger, httpRouter, copySnpCfg)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	logger.LogInfo("init user api server")
	// }
	err = initApiServer(apiServerConf.MysqlConf, logger, httpRouter)
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
	/* Dao */
	idInfoDao := &db.IDInfoDao{Logger: logger}
	companyDao := &db.CompanyDao{Logger: logger}
	accSubDao := &db.AccSubDao{Logger: logger}
	optInfoDao := &db.OperatorInfoDao{Logger: logger}
	loginInfoDao := &db.LoginInfoDao{Logger: logger}
	voucherInfoDao := &db.VoucherInfoDao{Logger: logger}
	voucherRecordDao := &db.VoucherRecordDao{Logger: logger}
	//初始化ID Resource
	service.GIdInfoService.InitIdInfoService(logger, idInfoDao, _db)
	ccErr := service.GIdInfoService.InitIdResource()
	if ccErr != nil {
		return ccErr
	}
	/*service*/
	accSubService := &service.AccountSubService{Logger: logger, AccSubDao: accSubDao, Db: _db}
	comService := &service.CompanyService{Logger: logger, CompanyDao: companyDao, Db: _db}
	optInfoService := &service.OperatorInfoService{Logger: logger, OptInfoDao: optInfoDao, Db: _db}
	authService := &service.AuthenService{Logger: logger, LogInfoDao: loginInfoDao, OptInfoDao: optInfoDao, Db: _db}
	vouInfoService := &service.VoucherInfoService{Logger: logger, VInfoDao: voucherInfoDao, Db: _db}
	voucherService := &service.VoucherService{Logger: logger, VRecordDao: voucherRecordDao, VInfoDao: voucherInfoDao, Db: _db}
	vouRecordService := &service.VoucherRecordService{Logger: logger, VRecordDao: voucherRecordDao, Db: _db}
	//handlers
	accSubHandlers := &handler.AccountSubHandlers{Logger: logger, AccSubService: accSubService}
	comHandlers := &handler.CompanyHandlers{Logger: logger, ComService: comService}
	optInfoHandlers := &handler.OperatorInfoHandlers{Logger: logger, ComService: comService, OptInfoService: optInfoService}
	voucherHandlers := &handler.VoucherHandlers{Logger: logger, Vis: vouInfoService, Vs: voucherService, Vrs: vouRecordService}
	authHandlers := &handler.AuthenHandlers{Logger: logger, AuthService: authService, ComService: comService, OptInfoService: optInfoService}

	httpRouter.RegisterFunc("CreateAccSub", accSubHandlers.CreateAccSub)
	httpRouter.RegisterFunc("DeleteAccSub", accSubHandlers.DeleteAccSub)
	httpRouter.RegisterFunc("ListAccSub", accSubHandlers.ListAccSub)
	httpRouter.RegisterFunc("GetAccSub", accSubHandlers.GetAccSub)
	httpRouter.RegisterFunc("UpdateAccSub", accSubHandlers.UpdateAccSub)

	httpRouter.RegisterFunc("CreateCompany", comHandlers.CreateCompany)
	httpRouter.RegisterFunc("DeleteCompany", comHandlers.DeleteCompany)
	httpRouter.RegisterFunc("GetCompany", comHandlers.GetCompany)
	httpRouter.RegisterFunc("ListCompany", comHandlers.ListCompany)
	httpRouter.RegisterFunc("UpdateCompany", comHandlers.UpdateCompany)

	httpRouter.RegisterFunc("CreateOperator", optInfoHandlers.CreateOperator)
	httpRouter.RegisterFunc("DeleteOperator", optInfoHandlers.DeleteOperator)
	httpRouter.RegisterFunc("GetOperatorInfo", optInfoHandlers.GetOperatorInfo)
	httpRouter.RegisterFunc("ListOperatorInfo", optInfoHandlers.ListOperatorInfo)
	httpRouter.RegisterFunc("UpdateOperator", optInfoHandlers.UpdateOperator)

	httpRouter.RegisterFunc("Login", authHandlers.Login)
	httpRouter.RegisterFunc("Logout", authHandlers.Logout)
	httpRouter.RegisterFunc("StatusCheckout", authHandlers.StatusCheckout)
	httpRouter.RegisterFunc("ListLoginInfo", authHandlers.ListLoginInfo)

	httpRouter.RegisterFunc("CreateVoucher", voucherHandlers.CreateVoucher)
	httpRouter.RegisterFunc("CreateVoucherRecords", voucherHandlers.CreateVoucherRecords)
	httpRouter.RegisterFunc("DeleteVoucher", voucherHandlers.DeleteVoucher)
	httpRouter.RegisterFunc("DeleteVoucherRecord", voucherHandlers.DeleteVoucherRecord)
	httpRouter.RegisterFunc("GetVoucherInfo", voucherHandlers.GetVoucherInfo)
	httpRouter.RegisterFunc("GetVoucher", voucherHandlers.GetVoucher)
	httpRouter.RegisterFunc("ListVoucherInfo", voucherHandlers.ListVoucherInfo)
	httpRouter.RegisterFunc("ListVoucherRecords", voucherHandlers.ListVoucherRecords)
	httpRouter.RegisterFunc("UpdateVoucherRecord", voucherHandlers.UpdateVoucherRecord)
	httpRouter.RegisterFunc("VoucherAudit", voucherHandlers.VoucherAudit)
	//检查是否登录
	handler.GAccessTokenH.InitAccessTokenHandler(authService, optInfoService, logger)
	httpRouter.LoginCheck = handler.GAccessTokenH.LoginCheck
	httpRouter.InterfaceAuthorityCheck = handler.GAuthManaged.InterfaceAuthorityCheck
	//用户登录的过期检查服务
	go handler.GAccessTokenH.ExpirationCheck()
	return nil
}
