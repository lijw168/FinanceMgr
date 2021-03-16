package constant

/* for api response code
 0:success
-1:
-2~-100:volume error
-101~-200:snapshot error
*/

//volume
const (
	CodeVolNotExist    = -2
	CodeVolBusy        = -3
	CodeVolNotEnough   = -4
	CodeVolInvalType   = -5
	CodeVolInvalSize   = -6
	CodeVolInvalResize = -7
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
)

//voucher information
const (
	CodeVoucherInfoNotExist = -401
)

//voucher record
const (
	CodeVoucherRecordNotExist = -501
)
