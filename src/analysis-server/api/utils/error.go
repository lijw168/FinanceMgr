package utils

const (
	ErrNull = ""
)

const (
	ErrSystem = "system"
	ErrFilter = "filter"
	ErrOrder  = "order"
	ErrTenant = "tenant"
	ErrHost   = "host"
	ErrPool   = "pool"
)

const (
	ErrError     = "error"
	ErrInvalid   = "invalid"
	ErrMiss      = "miss"
	ErrMalformed = "malformed"
)

const (
	ErrId          = "id"
	ErrIds         = "ids"
	ErrName        = "name"
	ErrTenantId    = "tenant_id"
	ErrField       = "field"
	ErrDefaultPool = "default_pool"
)

/*
 * R: Resource Type (vpc, subnet)
 * T: Error Type (invalid, notfound)
 * P: Parameter (name, cidr)
 * D: Detail
 */
type CcErrObj struct {
	R string
	T string
	P string
	D string
}

type CcError interface {
	Error() string
	Detail() string
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

func NewError(r, t, p, d string) CcError {
	return &CcErrObj{r, t, p, d}
}

func NewSysErr(err error) CcError {
	return NewError(ErrSystem, ErrError, ErrNull, err.Error())
}
