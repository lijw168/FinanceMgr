package util

const SplitStr string = ";"

//operation code
const (
	QuitApp = 1 + iota
	Heartbeat
	UserLogin
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
	UserOnline  = 1 // user login
	UserOffline = 2 // user logout
	InvalidUser = 3
)

const (
	SetLogLevel = "/setLogLevel"
)
