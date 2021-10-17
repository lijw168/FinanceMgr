package util

const SplitStr string = ";"

//operation code
const (
	QuitApp = 1 + iota
	Heartbeat
	UserLogin
	LoginInfoList
	UserLogout
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
	AccSubUpdate
	VoucherCreate
	VoucherUpdate
	VoucherDel
	VoucherShow
	VoucherArrange
	VouInfoShow
	VouInfoList
	VouInfoListByMulCon
	VouInfoListLatest
	VouInfoMaxNumOfMonth
	BatchAuditVouchers
	VouInfoUpdate
	VouRecordCreate
	VouRecordDel
	VouRecordsDel
	VouRecordList
	VouRecordUpdate
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
