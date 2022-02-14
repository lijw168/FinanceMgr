package main

import (
	"analysis-server/api/db"
	"analysis-server/api/handler"
	"analysis-server/api/service"
	"common/log"
	"common/url"
	"database/sql"
)

// func registerYearBalance(logger *log.Logger, httpRouter *url.UrlRouter, _db *sql.DB) {
// 	yearBalanceDao := &db.YearBalanceDao{Logger: logger}
// 	yearBalService := &service.YearBalanceService{Logger: logger, YearBalDao: yearBalanceDao, Db: _db}
// 	yearBalHandlers := &handler.YearBalHandlers{Logger: logger, YearBalService: yearBalService}
// 	httpRouter.RegisterFunc("GetYearBalance", yearBalHandlers.GetYearBalance)
// 	httpRouter.RegisterFunc("CreateYearBalance", yearBalHandlers.CreateYearBalance)
// 	httpRouter.RegisterFunc("UpdateYearBalance", yearBalHandlers.UpdateYearBalance)
// 	httpRouter.RegisterFunc("DeleteYearBalance", yearBalHandlers.DeleteYearBalance)
// }

//register voucher template
func registerVoucherTemplate(logger *log.Logger, httpRouter *url.UrlRouter, _db *sql.DB) {
	voucherTempDao := &db.VoucherTemplateDao{Logger: logger}
	voucherTempService := &service.VoucherTemplateService{Logger: logger, VTemplateDao: voucherTempDao, Db: _db}
	voucherTempHandlers := &handler.VoucherTemplateHandlers{Logger: logger, VoucherTempService: voucherTempService}
	httpRouter.RegisterFunc("CreateVoucherTemplate", voucherTempHandlers.CreateVoucherTemplate)
	httpRouter.RegisterFunc("DeleteVoucherTemplate", voucherTempHandlers.DeleteVoucherTemplate)
	httpRouter.RegisterFunc("GetVoucherTemplate", voucherTempHandlers.GetVoucherTemplate)
	httpRouter.RegisterFunc("ListVoucherTemplate", voucherTempHandlers.ListVoucherTemplate)
}

//register company group
func registerComGroup(logger *log.Logger, httpRouter *url.UrlRouter, comGroupDao *db.CompanyGroupDao,
	_db *sql.DB) {
	comGroupService := &service.CompanyGroupService{Logger: logger, ComGroupDao: comGroupDao, Db: _db}
	comGroupHandlers := &handler.CompanyGroupHandlers{Logger: logger, ComGroupService: comGroupService}
	httpRouter.RegisterFunc("CreateCompanyGroup", comGroupHandlers.CreateCompanyGroup)
	httpRouter.RegisterFunc("DeleteCompanyGroup", comGroupHandlers.DeleteCompanyGroup)
	httpRouter.RegisterFunc("GetCompanyGroup", comGroupHandlers.GetCompanyGroup)
	httpRouter.RegisterFunc("ListCompanyGroup", comGroupHandlers.ListCompanyGroup)
	httpRouter.RegisterFunc("UpdateCompanyGroup", comGroupHandlers.UpdateCompanyGroup)
}

//register companyHander
func registerCompany(logger *log.Logger, httpRouter *url.UrlRouter, comService *service.CompanyService,
	_db *sql.DB) {
	comHandlers := &handler.CompanyHandlers{Logger: logger, ComService: comService}
	httpRouter.RegisterFunc("CreateCompany", comHandlers.CreateCompany)
	httpRouter.RegisterFunc("DeleteCompany", comHandlers.DeleteCompany)
	httpRouter.RegisterFunc("GetCompany", comHandlers.GetCompany)
	httpRouter.RegisterFunc("ListCompany", comHandlers.ListCompany)
	httpRouter.RegisterFunc("UpdateCompany", comHandlers.UpdateCompany)
	httpRouter.RegisterFunc("AssociatedCompanyGroup", comHandlers.AssociatedCompanyGroup)
}

//register account subject
func registerAccSub(logger *log.Logger, httpRouter *url.UrlRouter, comDao *db.CompanyDao,
	voucherRecordDao *db.VoucherRecordDao, _db *sql.DB) {
	accSubDao := &db.AccSubDao{Logger: logger}
	accSubService := &service.AccountSubService{
		Logger:     logger,
		AccSubDao:  accSubDao,
		CompanyDao: comDao,
		VRecordDao: voucherRecordDao,
		Db:         _db}
	accSubHandlers := &handler.AccountSubHandlers{Logger: logger, AccSubService: accSubService}
	httpRouter.RegisterFunc("CreateAccSub", accSubHandlers.CreateAccSub)
	httpRouter.RegisterFunc("DeleteAccSub", accSubHandlers.DeleteAccSub)
	httpRouter.RegisterFunc("ListAccSub", accSubHandlers.ListAccSub)
	httpRouter.RegisterFunc("ListYearBalance", accSubHandlers.ListYearBalance)
	httpRouter.RegisterFunc("GetAccSub", accSubHandlers.GetAccSub)
	httpRouter.RegisterFunc("GetYearBalance", accSubHandlers.GetYearBalance)
	httpRouter.RegisterFunc("UpdateAccSub", accSubHandlers.UpdateAccSub)
	httpRouter.RegisterFunc("UpdateYearBalance", accSubHandlers.UpdateYearBalance)
	httpRouter.RegisterFunc("QueryAccSubReference", accSubHandlers.QueryAccSubReference)
}

//register operatorHander and authenHandler
func registerOptAndAuthenHandler(logger *log.Logger, httpRouter *url.UrlRouter,
	comService *service.CompanyService, _db *sql.DB) {
	optInfoDao := &db.OperatorInfoDao{Logger: logger}
	loginInfoDao := &db.LoginInfoDao{Logger: logger}
	optInfoService := &service.OperatorInfoService{Logger: logger, OptInfoDao: optInfoDao, Db: _db}
	authService := &service.AuthenService{Logger: logger, LogInfoDao: loginInfoDao, OptInfoDao: optInfoDao, Db: _db}
	optInfoHandlers := &handler.OperatorInfoHandlers{Logger: logger, ComService: comService, OptInfoService: optInfoService}
	authHandlers := &handler.AuthenHandlers{Logger: logger, AuthService: authService, ComService: comService, OptInfoService: optInfoService}
	httpRouter.RegisterFunc("CreateOperator", optInfoHandlers.CreateOperator)
	httpRouter.RegisterFunc("DeleteOperator", optInfoHandlers.DeleteOperator)
	httpRouter.RegisterFunc("GetOperatorInfo", optInfoHandlers.GetOperatorInfo)
	httpRouter.RegisterFunc("ListOperatorInfo", optInfoHandlers.ListOperatorInfo)
	httpRouter.RegisterFunc("UpdateOperator", optInfoHandlers.UpdateOperator)
	httpRouter.RegisterFunc("Login", authHandlers.Login)
	httpRouter.RegisterFunc("Logout", authHandlers.Logout)
	httpRouter.RegisterFunc("StatusCheckout", authHandlers.StatusCheckout)
	httpRouter.RegisterFunc("ListLoginInfo", authHandlers.ListLoginInfo)
	//检查是否登录
	handler.GAccessTokenH.InitAccessTokenHandler(authService, optInfoService, logger)
	httpRouter.LoginCheck = handler.GAccessTokenH.LoginCheck
	httpRouter.InterfaceAuthorityCheck = handler.GAuthManaged.InterfaceAuthorityCheck
	//用户登录的过期检查服务
	go handler.GAccessTokenH.ExpirationCheck()

}

//register resHander and voucherHandler
func registerResAndVoucherHandler(logger *log.Logger, httpRouter *url.UrlRouter, comDao *db.CompanyDao,
	voucherRecordDao *db.VoucherRecordDao, _db *sql.DB) {
	//voucher
	voucherInfoDao := &db.VoucherInfoDao{Logger: logger}
	//voucherRecordDao := &db.VoucherRecordDao{Logger: logger}
	vouDao := &db.VoucherDao{Logger: logger}
	vouInfoService := &service.VoucherInfoService{Logger: logger, VInfoDao: voucherInfoDao, Db: _db}
	voucherService := &service.VoucherService{Logger: logger, VRecordDao: voucherRecordDao, VInfoDao: voucherInfoDao, VouDao: vouDao, Db: _db}
	vouRecordService := &service.VoucherRecordService{Logger: logger, VRecordDao: voucherRecordDao, Db: _db}
	//resource
	resService := &service.ResouceInfoService{Logger: logger, VInfoDao: voucherInfoDao, CompanyDao: comDao, Db: _db}
	resHandlers := &handler.ResourceInfoHandlers{Logger: logger, ResService: resService}
	voucherHandlers := &handler.VoucherHandlers{Logger: logger, Vis: vouInfoService, Vs: voucherService, Vrs: vouRecordService}
	httpRouter.RegisterFunc("InitResourceInfo", resHandlers.InitResourceInfo)
	//voucher
	httpRouter.RegisterFunc("CreateVoucher", voucherHandlers.CreateVoucher)
	httpRouter.RegisterFunc("UpdateVoucher", voucherHandlers.UpdateVoucher)
	httpRouter.RegisterFunc("DeleteVoucher", voucherHandlers.DeleteVoucher)
	httpRouter.RegisterFunc("ArrangeVoucher", voucherHandlers.ArrangeVoucher)
	// httpRouter.RegisterFunc("CreateVoucherRecords", voucherHandlers.CreateVoucherRecords)
	// httpRouter.RegisterFunc("DeleteVoucherRecord", voucherHandlers.DeleteVoucherRecord)
	// httpRouter.RegisterFunc("DeleteVoucherRecords", voucherHandlers.DeleteVoucherRecords)
	httpRouter.RegisterFunc("ListVoucherRecords", voucherHandlers.ListVoucherRecords)
	//httpRouter.RegisterFunc("UpdateVoucherRecordByID", voucherHandlers.UpdateVoucherRecordByID)
	httpRouter.RegisterFunc("GetVoucherInfo", voucherHandlers.GetVoucherInfo)
	httpRouter.RegisterFunc("GetVoucher", voucherHandlers.GetVoucher)
	httpRouter.RegisterFunc("GetLatestVoucherInfo", voucherHandlers.GetLatestVoucherInfo)
	httpRouter.RegisterFunc("ListVoucherInfo", voucherHandlers.ListVoucherInfo)
	httpRouter.RegisterFunc("ListVoucherInfoByMulCondition", voucherHandlers.ListVoucherInfoByMulCondition)
	httpRouter.RegisterFunc("GetMaxNumOfMonth", voucherHandlers.GetMaxNumOfMonth)
	httpRouter.RegisterFunc("UpdateVoucherInfo", voucherHandlers.UpdateVoucherInfo)
	httpRouter.RegisterFunc("BatchAuditVouchers", voucherHandlers.BatchAuditVouchers)
	httpRouter.RegisterFunc("CalculateAccumulativeMoney", voucherHandlers.CalculateAccumulativeMoney)

}

//register menuHandler
func registerMenuHandler(logger *log.Logger, httpRouter *url.UrlRouter, _db *sql.DB) {
	menuInfoDao := &db.MenuInfoDao{Logger: logger}
	menuService := &service.MenuInfoService{Logger: logger, MenuDao: menuInfoDao, Db: _db}
	menuHandlers := &handler.MenuInfoHandlers{Logger: logger, MenuService: menuService}
	httpRouter.RegisterFunc("ListMenuInfo", menuHandlers.ListMenuInfo)
}
