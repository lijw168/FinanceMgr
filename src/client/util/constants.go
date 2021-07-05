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
	CompanyUpdate
	AccSubCreate
	AccSubList
	AccSubShow
	AccSubDel
	AccSubUpdate
	VoucherCreate
	VoucherDel
	VoucherShow
	VouInfoShow
	VouInfoList
	VouRecordCreate
	VouRecordDel
	VouRecordList
	VouRecordUpdate
)

const (
	InvalidStatus = 0 //invalid status
	Online        = 1 // user online status
	Offline       = 2 // user offline	status
)

const (
	SetLogLevel = "/setLogLevel"
)
