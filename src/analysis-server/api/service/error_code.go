package service

const (
	ErrNull = ""
)

const (
	ErrSystem          = "system"
	ErrAccSub          = "accountSubject"
	ErrCompany         = "company"
	ErrComGroup        = "companyGroup"
	ErrOperator        = "operator"
	ErrVoucherInfo     = "voucherInfo"
	ErrVoucherTemplate = "voucherTemplate"
	ErrVoucher         = "voucher"
	ErrOrder           = "order"
	ErrLogin           = "login"
	ErrLogout          = "logout"
	ErrIdInfo          = "IdInfo"
	ErrMenuInfo        = "MenuInfo"
	ErrResInfo         = "ResourceInfo"
	ErrYearBalance     = "yearBalance"
)

const (
	ErrError         = "error"
	ErrMantaining    = "mantaining"
	ErrInvalid       = "invalid"
	ErrMiss          = "miss"
	ErrMalformed     = "malformed"
	ErrNotFound      = "not found"
	ErrInuse         = "inuse"
	ErrConflict      = "conflict"
	ErrUnsupported   = "unsupported"
	ErrForbidden     = "forbidden"
	ErrLimit         = "limit"
	ErrExceeded      = "exceeded"
	ErrDuplicate     = "duplicate"
	ErrParam         = "invalidParam"
	ErrNotAllowed    = "notallowed"
	ErrValue         = "invalid value type"
	ErrChangeContent = "change content"
	ErrEmpty         = "empty"
	ErrMaintenance   = "maintenance"
	ErrDismatch      = "dismatch"
)

const (
	ErrId                = "id"
	ErrIds               = "ids"
	ErrName              = "name"
	ErrCommonId          = "commonId"
	ErrCompanyId         = "companyId"
	ErrCookie            = "cookie"
	ErrPasswd            = "password"
	ErrVouYear           = "voucher_year"
	ErrVouMon            = "voucher_month"
	ErrVouFiller         = "voucher_filler"
	ErrVouAuditor        = "voucher_auditor"
	ErrVouRecSub         = "voucher_record_subject"
	ErrVouRecDebit       = "voucher_record_debit"
	ErrVouRecCredit      = "voucher_record_crebit"
	ErrVoucherData       = "voucher debit or crebit data"
	ErrVouSummary        = "voucher_summary"
	ErrVoucherRecord     = "voucher_record"
	ErrVoucherTemplateID = "voucher_template_id"
	ErrField             = "field"
	ErrOd                = "direction" // order direction
	ErrType              = "type"
	ErrSubLevel          = "subjectLevel"
	ErrSubdir            = "subjectDirection"
	ErrSubStyle          = "subjectStyle"
	ErrSize              = "size"
	ErrOffset            = "offset"
	ErrSubjectName       = "subjectName"
	ErrRole              = "role"
	ErrStatus            = "status"
	ErrAttachParam       = "attach parameter"
	ErrNoAuthority       = "no authority"
	ErrAccountPeriod     = "account period"
	ErrYear              = "year"
	ErrBalance           = "year balance"
)

const (
	ErrRecordExist    = "该记录已存在"
	ErrFiledDuplicate = "存在重复的字段"
)

/*
 * C: Error Code
 * R: Resource Type
 * T: Error Type (invalid, notfound)
 * P: Parameter (name, cidr)
 * D: Detail
 */
type CcErrObj struct {
	C int
	R string
	T string
	P string
	D string
}

type CcError interface {
	Error() string
	Detail() string
	GetCode() int
	SetCode(errCode int)
}

func (e *CcErrObj) Error() string {
	s := e.R + "." + e.T
	if e.P != ErrNull {
		s += "." + e.P
	}
	return s
}

func (e *CcErrObj) Detail() string {
	return e.D
}

func (e *CcErrObj) GetCode() int {
	return e.C
}

func (e *CcErrObj) SetCode(errCode int) {
	e.C = errCode
}

func NewError(r, t, p, d string) CcError {
	return &CcErrObj{-1, r, t, p, d}
}

func NewCcError(c int, r, t, p, d string) CcError {
	return &CcErrObj{c, r, t, p, d}
}

func NewSysErr(err error) CcError {
	return NewError(ErrSystem, ErrError, ErrNull, err.Error())
}
