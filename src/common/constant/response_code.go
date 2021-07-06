package constant

/* for api response code
 0:success
-1:
-2~-100:common error
-101~-200:account subject error
*/

//common
const (
	CodeNoLogin = -2
)

//account subject
const (
	CodeAccSubNotExist   = -101
	CodeInvalAccSubLevel = -102
)

//company information
const (
	CodeComInfoNotExist = -201
)

//operator information
const (
	CodeOptInfoNotExist = -301
	CodeUserNameWrong   = -302
	CodePasswdWrong     = -303
)

//voucher information
const (
	CodeVoucherInfoNotExist = -401
)

//voucher record
const (
	CodeVoucherRecordNotExist = -501
)

//IdInfo
const (
	CodeIdInfoNotExist = -601
)
