package utils

//user status
const (
	InvalidStatus = 0 //invalid status
	Online        = 1 // user online status
	Offline       = 2 // user offline	status
)

//voucher record status
const (
	NoAuditVoucher = 1 // 该凭证未审核
	InvalidVoucher = 2 // 该凭证作废
	AuditedVoucher = 3 // 该凭证已审核
)
