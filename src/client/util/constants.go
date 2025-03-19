package util

const SplitStr string = ";"

// operation code
const (
	QuitApp = 1 + iota
	Heartbeat
	Login
	LoginInfoList
	Logout
	OperatorCreate
	OperatorList
	OperatorShow
	OperatorDel
	OperatorUpdate
	CompanyCreate
	CompanyList
	CompanyShow
	CompanyDel
	InitResourceInfo
	CompanyUpdate
	AccSubCreate
	AccSubReferenceQuery
	AccSubList
	AccSubShow
	AccSubDel
	CopyAccSubTemplate
	AccSubUpdate
	VoucherCreate
	VoucherUpdate
	VoucherDel
	VoucherShow
	VoucherArrange
	VouInfoShow
	VouInfoList
	VouInfoListWithAuxCond
	VouInfoListLatest
	VouInfoMaxNumOfMonth
	BatchAuditVouchers
	VouInfoUpdate
	VoucherInfoNoAuditedShow
	// VouRecordCreate
	// VouRecordDel
	// VouRecordsDel
	VouRecordList
	// VouRecordUpdate
	CalculateAccuMoney
	BatchCalcAccuMoney
	CalcAccountOfPeriod
	VouTemplateCreate
	VouTemplateDel
	VouTemplateShow
	VouTemplateList
	YearBalanceCreate
	YearBalanceBatchCreate
	YearBalanceList
	YearBalanceShow
	YearBalanceDel
	YearBalanceBatchDel
	AccSubYearBalValueShow
	AnnualClosing
	CancelAnnualClosing
	AnnualClosingStatusShow
	BatchUpdateBals
	YearBalanceUpdate
	MenuInfoList
)

const (
	InvalidStatus = 0 //invalid status
	Online        = 1 // user online status
	Offline       = 2 // user offline	status
)

const (
	SetLogLevel = "/setLogLevel"
)
