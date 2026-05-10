package main

import (
	"financeMgr/src/analysis-server/api/db"
	"financeMgr/src/analysis-server/api/handler"
	"financeMgr/src/analysis-server/api/service"
	"financeMgr/src/common/url"
)

func registerYearBalance(httpRouter *url.UrlRouter) {
	yearBalanceDao := &db.YearBalanceDao{}
	yearBalService := &service.YearBalanceService{YearBalDao: yearBalanceDao}
	yearBalHandlers := &handler.YearBalHandlers{YearBalService: yearBalService}
	httpRouter.RegisterFunc("GetYearBalance", yearBalHandlers.GetYearBalance)
	httpRouter.RegisterFunc("GetAccSubYearBalValue", yearBalHandlers.GetAccSubYearBalValue)
	httpRouter.RegisterFunc("CreateYearBalance", yearBalHandlers.CreateYearBalance)
	httpRouter.RegisterFunc("BatchCreateYearBalance", yearBalHandlers.BatchCreateYearBalance)
	httpRouter.RegisterFunc("UpdateYearBalance", yearBalHandlers.UpdateYearBalance)
	httpRouter.RegisterFunc("BatchUpdateBals", yearBalHandlers.BatchUpdateBals)
	httpRouter.RegisterFunc("DeleteYearBalance", yearBalHandlers.DeleteYearBalance)
	httpRouter.RegisterFunc("BatchDeleteYearBalance", yearBalHandlers.BatchDeleteYearBalance)
	httpRouter.RegisterFunc("ListYearBalance", yearBalHandlers.ListYearBalance)
	httpRouter.RegisterFunc("AnnualClosing", yearBalHandlers.AnnualClosing)
	httpRouter.RegisterFunc("CancelAnnualClosing", yearBalHandlers.CancelAnnualClosing)
	httpRouter.RegisterFunc("GetAnnualClosingStatus", yearBalHandlers.GetAnnualClosingStatus)
}

// register voucher template
func registerVoucherTemplate(httpRouter *url.UrlRouter) {
	voucherTempDao := &db.VoucherTemplateDao{}
	voucherTempService := &service.VoucherTemplateService{VTemplateDao: voucherTempDao}
	voucherTempHandlers := &handler.VoucherTemplateHandlers{VoucherTempService: voucherTempService}
	httpRouter.RegisterFunc("CreateVoucherTemplate", voucherTempHandlers.CreateVoucherTemplate)
	httpRouter.RegisterFunc("DeleteVoucherTemplate", voucherTempHandlers.DeleteVoucherTemplate)
	httpRouter.RegisterFunc("GetVoucherTemplate", voucherTempHandlers.GetVoucherTemplate)
	httpRouter.RegisterFunc("ListVoucherTemplate", voucherTempHandlers.ListVoucherTemplate)
}

// register company group
func registerComGroup(httpRouter *url.UrlRouter, comGroupDao *db.CompanyGroupDao) {
	comGroupService := &service.CompanyGroupService{ComGroupDao: comGroupDao}
	comGroupHandlers := &handler.CompanyGroupHandlers{ComGroupService: comGroupService}
	httpRouter.RegisterFunc("CreateCompanyGroup", comGroupHandlers.CreateCompanyGroup)
	httpRouter.RegisterFunc("DeleteCompanyGroup", comGroupHandlers.DeleteCompanyGroup)
	httpRouter.RegisterFunc("GetCompanyGroup", comGroupHandlers.GetCompanyGroup)
	httpRouter.RegisterFunc("ListCompanyGroup", comGroupHandlers.ListCompanyGroup)
	httpRouter.RegisterFunc("UpdateCompanyGroup", comGroupHandlers.UpdateCompanyGroup)
}

// register companyHander
func registerCompany(httpRouter *url.UrlRouter, comService *service.CompanyService) {
	comHandlers := &handler.CompanyHandlers{ComService: comService}
	httpRouter.RegisterFunc("CreateCompany", comHandlers.CreateCompany)
	httpRouter.RegisterFunc("DeleteCompany", comHandlers.DeleteCompany)
	httpRouter.RegisterFunc("GetCompany", comHandlers.GetCompany)
	httpRouter.RegisterFunc("ListCompany", comHandlers.ListCompany)
	httpRouter.RegisterFunc("UpdateCompany", comHandlers.UpdateCompany)
	httpRouter.RegisterFunc("AssociatedCompanyGroup", comHandlers.AssociatedCompanyGroup)
	httpRouter.RegisterFunc("ListCompanyAccountYearInfo", comHandlers.ListCompanyAccountYearInfo)
}

// register account subject
func registerAccSub(httpRouter *url.UrlRouter, comDao *db.CompanyDao,
	voucherRecordDao *db.VoucherRecordDao) {
	accSubDao := &db.AccSubDao{}
	accSubService := &service.AccountSubService{
		AccSubDao:  accSubDao,
		CompanyDao: comDao,
		VRecordDao: voucherRecordDao}
	accSubHandlers := &handler.AccountSubHandlers{AccSubService: accSubService}
	httpRouter.RegisterFunc("CreateAccSub", accSubHandlers.CreateAccSub)
	httpRouter.RegisterFunc("DeleteAccSub", accSubHandlers.DeleteAccSub)
	httpRouter.RegisterFunc("ListAccSub", accSubHandlers.ListAccSub)
	httpRouter.RegisterFunc("GetAccSub", accSubHandlers.GetAccSub)
	httpRouter.RegisterFunc("UpdateAccSub", accSubHandlers.UpdateAccSub)
	httpRouter.RegisterFunc("QueryAccSubReference", accSubHandlers.QueryAccSubReference)
	httpRouter.RegisterFunc("CopyAccSubTemplate", accSubHandlers.CopyAccSubTemplate)
	httpRouter.RegisterFunc("GenerateAccSubTemplate", accSubHandlers.GenerateAccSubTemplate)
}

// register operatorHander and authenHandler
func registerOptAndAuthenHandler(httpRouter *url.UrlRouter, comService *service.CompanyService) {
	optInfoDao := &db.OperatorInfoDao{}
	loginInfoDao := &db.LoginInfoDao{}
	optInfoService := &service.OperatorInfoService{OptInfoDao: optInfoDao}
	authService := &service.AuthenService{LogInfoDao: loginInfoDao, OptInfoDao: optInfoDao}
	optInfoHandlers := &handler.OperatorInfoHandlers{ComService: comService, OptInfoService: optInfoService}
	authHandlers := &handler.AuthenHandlers{AuthService: authService, ComService: comService, OptInfoService: optInfoService}
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
	handler.GAccessTokenH.InitAccessTokenHandler(authService, optInfoService)
	httpRouter.LoginCheck = handler.GAccessTokenH.LoginCheck
	httpRouter.InterfaceAuthorityCheck = handler.GAuthManaged.InterfaceAuthorityCheck
	//用户登录的过期检查服务
	//go handler.GAccessTokenH.ExpirationCheck()

}

// register resHander and voucherHandler
func registerResAndVoucherHandler(httpRouter *url.UrlRouter, comDao *db.CompanyDao, voucherRecordDao *db.VoucherRecordDao) {
	//voucher
	voucherInfoDao := &db.VoucherInfoDao{}
	//voucherRecordDao := &db.VoucherRecordDao{}
	vouDao := &db.VoucherDao{}
	vouInfoService := &service.VoucherInfoService{VInfoDao: voucherInfoDao}
	voucherService := &service.VoucherService{VRecordDao: voucherRecordDao, VInfoDao: voucherInfoDao, VouDao: vouDao}
	vouRecordService := &service.VoucherRecordService{VRecordDao: voucherRecordDao}
	//resource
	//resService := &service.ResouceInfoService{CompanyDao: comDao}
	//resHandlers := &handler.ResourceInfoHandlers{ResService: resService}
	//httpRouter.RegisterFunc("InitResourceInfo", resHandlers.InitResourceInfo)
	//voucher
	voucherHandlers := &handler.VoucherHandlers{Vis: vouInfoService, Vs: voucherService, Vrs: vouRecordService}
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
	httpRouter.RegisterFunc("ListVoucherInfoWithAuxCondition", voucherHandlers.ListVoucherInfoWithAuxCondition)
	httpRouter.RegisterFunc("GetMaxNumOfMonth", voucherHandlers.GetMaxNumOfMonth)
	httpRouter.RegisterFunc("UpdateVoucherInfo", voucherHandlers.UpdateVoucherInfo)
	httpRouter.RegisterFunc("BatchAuditVouchers", voucherHandlers.BatchAuditVouchers)
	httpRouter.RegisterFunc("CalculateAccumulativeMoney", voucherHandlers.CalculateAccumulativeMoney)
	httpRouter.RegisterFunc("BatchCalcAccuMoney", voucherHandlers.BatchCalcAccuMoney)
	httpRouter.RegisterFunc("CalcAccountOfPeriod", voucherHandlers.CalcAccountOfPeriod)
	httpRouter.RegisterFunc("GetNoAuditedVoucherInfoCount", voucherHandlers.GetNoAuditedVoucherInfoCount)
}

// register menuHandler
func registerMenuHandler(httpRouter *url.UrlRouter) {
	menuInfoDao := &db.MenuInfoDao{}
	menuService := &service.MenuInfoService{MenuDao: menuInfoDao}
	menuHandlers := &handler.MenuInfoHandlers{MenuService: menuService}
	httpRouter.RegisterFunc("ListMenuInfo", menuHandlers.ListMenuInfo)
}
