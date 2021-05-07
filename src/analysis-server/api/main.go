package main

import (
	"flag"
	"fmt"
	"net/http"

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
	aUtils "analysis-server/api/utils"
)

var (
	gGenSubIdInfo    *aUtils.GenIdInfo
	gGenComIdInfo    *aUtils.GenIdInfo
	gGenVouIdInfo    *aUtils.GenIdInfo
	gGenVouRecIdInfo *aUtils.GenIdInfo
	gIdInfoService   *service.IDInfoService
	exitCh           = make(chan bool)
)

func updateIdInfo() {
	subId := gGenSubIdInfo.GetId()
	comId := gGenComIdInfo.GetId()
	vouId := gGenVouIdInfo.GetId()
	vouRecId := gGenVouRecIdInfo.GetId()
	updateFields := make(map[string]interface{})
	updateFields["SubjectID"] = subId
	updateFields["CompanyID"] = comId
	updateFields["VoucherID"] = vouId
	updateFields["VoucherRecordID"] = vouRecId
	gIdInfoService.UpdateIdInfo(updateFields)
}
func interceptSignal() {
	daemonExitCh := make(chan os.Signal)
	signal.Notify(daemonExitCh, syscall.SIGTERM, syscall.SIGKILL,
		syscall.SIGQUIT, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		for {
			sig := <-daemonExitCh
			fmt.Printf("the sig is %s\n", sig.String())
			updateIdInfo()
			break
		}
		exitCh <- true
	}()
}

func waitDaemonExit() {
	<-exitCh
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
	gIdInfoService = new(service.IDInfoService)
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
	gIdInfoService.Logger = logger
	gIdInfoService.IdInfoDao = idInfoDao
	gIdInfoService.Db = _db
	/*service*/
	idInfoService := &service.IDInfoService{Logger: logger, IdInfoDao: idInfoDao, Db: _db}
	idInfoView, err2 := idInfoService.GetIdInfo()
	if err2 != nil {
		fmt.Println("[Init] GetIdInfo failed,error: ", err2.Error())
		return err2
	}
	gGenSubIdInfo, err = aUtils.NewGenIdInfo(idInfoView.SubjectID)
	if err != nil {
		fmt.Println("[Init] initialize subjectID ,failed. error: ", err)
		return err
	}
	gGenComIdInfo, err = aUtils.NewGenIdInfo(idInfoView.CompanyID)
	if err != nil {
		fmt.Println("[Init] initialize companyID ,failed. error: ", err)
		return err
	}
	gGenVouIdInfo, err = aUtils.NewGenIdInfo(idInfoView.VoucherID)
	if err != nil {
		fmt.Println("[Init] initialize voucherID ,failed. error: ", err)
		return err
	}
	gGenVouRecIdInfo, err = aUtils.NewGenIdInfo(idInfoView.VoucherRecordID)
	if err != nil {
		fmt.Println("[Init] initialize voucherRecordID ,failed. error: ", err)
		return err
	}

	accSubService := &service.AccountSubService{Logger: logger, AccSubDao: accSubDao, Db: _db, GenSubId: gGenSubIdInfo}
	comService := &service.CompanyService{Logger: logger, CompanyDao: companyDao, Db: _db, GenComId: gGenComIdInfo}
	optInfoService := &service.OperatorInfoService{Logger: logger, OptInfoDao: optInfoDao, Db: _db}
	authService := &service.AuthenService{Logger: logger, LogInfoDao: loginInfoDao, OptInfoDao: optInfoDao, Db: _db}
	vouInfoService := &service.VoucherInfoService{Logger: logger, VInfoDao: voucherInfoDao, Db: _db}
	voucherService := &service.VoucherService{Logger: logger, VRecordDao: voucherRecordDao, VInfoDao: voucherInfoDao,
		GenRecordId: gGenVouRecIdInfo, GenVoucherId: gGenVouIdInfo, Db: _db}
	vouRecordService := &service.VoucherRecordService{Logger: logger, VRecordDao: voucherRecordDao, GenRecordId: gGenVouRecIdInfo, Db: _db}
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
	//httpRouter.RegisterFunc("GetLoginInfo", authHandlers.GetLoginInfo)
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
	//检查是否登录
	httpRouter.CheckoutCall = handler.AccessToken.Checkout
	return nil
}
