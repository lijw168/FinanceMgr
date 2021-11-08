package utils

//user status
const (
	InvalidStatus = 0 //invalid status
	Online        = 1 // user online status
	Offline       = 2 // user offline	status
)

//voucher status
const (
	NoAuditVoucher = 1 // 该凭证未审核
	InvalidVoucher = 2 // 该凭证作废
	AuditedVoucher = 3 // 该凭证已审核
)

//company group status
const (
	InvalidCompanyGroup = 0 // 该公司组无效
	ValidCompanyGroup   = 1 // 该公司组有效
)
