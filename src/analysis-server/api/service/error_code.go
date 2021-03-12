package service

const (
	ErrNull = ""
)

const (
	ErrSystem      = "system"
	ErrAccSub      = "accountSubject"
	ErrCompany     = "company"
	ErrOperator    = "operator"
	ErrVoucherInfo = "voucherInfo"
	ErrVoucher     = "voucher"
	ErrTenant      = "tenant"
	ErrHost        = "host"
	ErrPool        = "pool"
	ErrRack        = "rack"
	ErrDisk        = "disk"
	ErrRg          = "rg"
	ErrProxy       = "proxy"
	ErrVolume      = "volume"
	ErrSnapshot    = "snapshot"
	ErrQuota       = "quota"
	ErrVolumeAtt   = "volumeAttachment"
	ErrTask        = "task"
	ErrVolumeTask  = "volumeTask"
	ErrAdminTask   = "adminTask"
	ErrRT          = "replicaTask"
	ErrSpeedLimit  = "speedlimit"
	ErrCacheList   = "cachelist"
)

const (
	ErrError          = "error"
	ErrMantaining     = "mantaining"
	ErrInvalid        = "invalid"
	ErrMiss           = "miss"
	ErrMalformed      = "malformed"
	ErrNotFound       = "notfound"
	ErrInuse          = "inuse"
	ErrConflict       = "conflict"
	ErrUnsupported    = "unsupported"
	ErrForbidden      = "forbidden"
	ErrLimit          = "limit"
	ErrExceeded       = "exceeded"
	ErrDuplicate      = "duplicate"
	ErrNorse          = "norse"
	ErrStatusDismatch = "volume status is dismatch"
	ErrParam          = "invalidParam"
	ErrNotAllowed     = "notallowed"
	ErrValue          = "invalid value type"
	ErrChangeContent  = "change content"
	ErrEmpty          = "empty"
	ErrMaintenance    = "maintenance"
	ErrDismatch       = "dismatch"
	ErrBusy           = "busy"
)

const (
	ErrId            = "id"
	ErrIds           = "ids"
	ErrName          = "name"
	ErrPasswd        = "password"
	ErrVouMon        = "voucher_month"
	ErrVouRecSub     = "voucher_record_subject"
	ErrVouRecDebit   = "voucher_record_debit"
	ErrVouRecCredit  = "voucher_record_crebit"
	ErrVoucherRecord = "voucher_record"
	ErrField         = "field"
	ErrOd            = "direction" // order direction
	ErrType          = "type"
	ErrFixedIp       = "fixed_ip"
	ErrMgmtIp        = "mgmt_ip"
	ErrDataIp        = "underlay"
	ErrAdminStaus    = "admin_status"
	ErrTenantId      = "tenant_id"
	ErrTenantIds     = "tenant_ids"
	ErrSubLevel      = "subjectLevel"
	ErrSize          = "size"
	ErrOffset        = "offset"
	ErrVolumeId      = "volume_id"
	ErrSnapshotId    = "snapshot_id"
	ErrSnapshotName  = "snapshot_name"
	ErrHostIp        = "host_ip"
	ErrDescription   = "description"
	ErrSubjectName   = "subjectName"
	ErrInstanceUuid  = "instance_uuid"
	ErrInstanceType  = "instance_type"
	ErrImageFormat   = "image_format"
	ErrImageId       = "image_id"
	ErrImageLocation = "image_location"
	ErrImageType     = "image_type"
	ErrPublicImage   = "public_image"
	ErrImageHash     = "image_hash"
	ErrRegion        = "region"
)

/*
 * C: Error Code
 * R: Resource Type (vpc, subnet)
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
	Code() int
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

func (e *CcErrObj) Code() int {
	return e.C
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
