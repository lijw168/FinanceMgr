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

	"analysis-server/api/cfg"
	"analysis-server/api/db"
	"analysis-server/api/handler"
	"analysis-server/api/service"
	aUtils "analysis-server/api/utils"
)

func main() {

	if utils.SetLimit() != nil {
		fmt.Println("[Init] set max open files failed")
		return
	}

	// all in one
	var apiServerCfgFile = flag.String("c",
		"/etc/analysis/api_server.cfg", "Server config file name")

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

	if err = startServer(httpRouter, apiServerConf.ServerConf); err != nil {
		fmt.Println("[Init] http server exit, error: ", err)
	}
	logger.Close()
	return
}

func startServer(router *url.UrlRouter, serverConf *cfg.ServerConf) error {

	runtime.GOMAXPROCS(serverConf.Cores)

	http.Handle(serverConf.BaseUrl, router)

	return http.ListenAndServe(":"+strconv.Itoa(serverConf.Port), nil)
}

func handlerInit(httpRouter *url.UrlRouter, logger *log.Logger, apiServerConf *cfg.ApiServerCfg) error {
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
	_db, err := config.MysqlInstance{Conf: smysqlConf, Logger: logger}.NewMysqlInstance()
	if err != nil {
		fmt.Println("[Init] Create Db connection error: ", err)
		return err
	}
	/* Dao */
	idInfoDao := &db.IDInfoDao{Logger: logger}
	companyDao := &db.CompanyDao{Logger: logger}
	accSubDao := &db.AccSubDao{Logger: logger}
	optInfoDao := &db.OperatorInfoDao{Logger: logger}
	voucherInfoDao := &db.VoucherInfoDao{Logger: logger}
	voucherRecordDao := &db.VoucherRecordDao{Logger: logger}
	/*service*/
	idInfoService := &service.IDInfoService{Logger: logger, IdInfoDao: idInfoDao, Db: _db}
	idInfoView := idInfoService.GetIdInfo()
	genSubIdInfo := aUtils.NewGenIdInfo(idInfoView.SubjectID)
	genComIdInfo := aUtils.NewGenIdInfo(idInfoView.CompanyID)
	genVouIdInfo := aUtils.NewGenIdInfo(idInfoView.VoucherID)
	genVouRecIdInfo := aUtils.NewGenIdInfo(idInfoView.VoucherRecordID)

	accSubService := &service.AccountSubService{Logger: logger, AccSubDao: accSubDao, Db: _db, GenSubId: genSubIdInfo}
	comService := &service.CompanyService{Logger: logger, CompanyDao: companyDao, Db: _db, GenComId: genComIdInfo}
	optInfoService := &service.OperatorInfoService{Logger: logger, OptInfoDao: optInfoDao, Db: _db}
	vouInfoService := &service.VoucherInfoService{Logger: logger, VInfoDao: voucherInfoDao, Db: _db}
	voucherService := &service.VoucherService{Logger: logger, VRecordDao: voucherRecordDao, VInfoDao: voucherInfoDao,
		GenRecordId: genVouRecIdInfo, GenVoucherId: genVouIdInfo, Db: _db}
	vouRecordService := &service.VoucherRecordService{Logger: logger, VRecordDao: voucherRecordDao, GenRecordId: genVouRecIdInfo, Db: _db}
	//handlers
	accSubHandlers := &handler.AccountSubHandlers{Logger: logger, AccSubService: accSubService}
	comHandlers := &handler.CompanyHandlers{Logger: logger, ComService: comService}
	optInfoHandlers := &handler.OperatorInfoHandlers{Logger: logger, ComService: comService, OptInfoService: optInfoService}
	voucherHandlers := &handler.VoucherHandlers{Logger: logger, Vis: vouInfoService, Vs: voucherService, Vrs: vouRecordService}

	httpRouter.RegisterFunc("CreateAccSub", accSubHandlers.CreateAccSub)
	httpRouter.RegisterFunc("DeleteAccSub", accSubHandlers.DeleteAccSub)
	httpRouter.RegisterFunc("ListAccSub", accSubHandlers.ListAccSub)
	httpRouter.RegisterFunc("GetAccSub", accSubHandlers.GetAccSub)

	httpRouter.RegisterFunc("CreateCompany", comHandlers.CreateCompany)
	httpRouter.RegisterFunc("DeleteCompany", comHandlers.DeleteCompany)
	httpRouter.RegisterFunc("GetCompany", comHandlers.GetCompany)
	httpRouter.RegisterFunc("Listcompany", comHandlers.Listcompany)
	httpRouter.RegisterFunc("UpdateCompany", comHandlers.UpdateCompany)

	httpRouter.RegisterFunc("CreateOperator", optInfoHandlers.CreateOperator)
	httpRouter.RegisterFunc("DeleteOperator", optInfoHandlers.DeleteOperator)
	httpRouter.RegisterFunc("GetOperatorInfo", optInfoHandlers.GetOperatorInfo)
	httpRouter.RegisterFunc("ListOperatorInfo", optInfoHandlers.ListOperatorInfo)
	httpRouter.RegisterFunc("UpdateOperator", optInfoHandlers.UpdateOperator)

	httpRouter.RegisterFunc("CreateVoucher", voucherHandlers.CreateVoucher)
	httpRouter.RegisterFunc("CreateVoucherRecords", voucherHandlers.CreateVoucherRecords)
	httpRouter.RegisterFunc("DeleteVoucher", voucherHandlers.DeleteVoucher)
	httpRouter.RegisterFunc("DeleteVoucherRecord", voucherHandlers.DeleteVoucherRecord)
	httpRouter.RegisterFunc("GetVoucherInfo", voucherHandlers.GetVoucherInfo)
	httpRouter.RegisterFunc("GetVoucher", voucherHandlers.GetVoucher)
	httpRouter.RegisterFunc("ListVoucherInfo", voucherHandlers.ListVoucherInfo)
	httpRouter.RegisterFunc("ListVoucherRecords", voucherHandlers.ListVoucherRecords)
	httpRouter.RegisterFunc("UpdateVoucherRecord", voucherHandlers.UpdateVoucherRecord)
	return nil
}
