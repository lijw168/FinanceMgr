package constant

/* for api response code
 0:success
-1:
-2~-100:common error
-101~-200:account subject error
*/

// common
const (
	CodeNoLogin     = -2
	CodeNoAuthority = -3
)

// account subject
const (
	CodeAccSubNotExist   = -101
	CodeInvalAccSubLevel = -102
	CodeInvalAccSubDir   = -103
	CodeInvalAccSubType  = -104
)

// company information
const (
	CodeComInfoNotExist = -201
)

// operator information
const (
	CodeOptInfoNotExist = -301
	CodeUserNameWrong   = -302
	CodePasswdWrong     = -303
)

// voucher information
const (
	CodeVoucherInfoNotExist     = -401
	CodeVoucherRecordNotExist   = -402
	CodeVoucherTemplateNotExist = -403
)

// the begin of year balance
const (
	CodeYearBalanceNotExist = -501
)

// IdInfo
const (
	CodeIdInfoNotExist = -601
)

// companyGroup
const (
	CodeCompanyGroupNotExist = -701
)
