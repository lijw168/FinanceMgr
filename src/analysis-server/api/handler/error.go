package handler

const (
	InvalidParams         = 4000
	Unsupport             = 4600
	NotFound              = 5000
	ResourceOperatefailed = 5100
	InternalError         = 6000
)

// InvalidParams
const (
	IncompleteParam         = 4001
	InvalidCidr             = 4002
	InvalidName             = 4003
	DuplicateVpcName        = 4004
	DuplicateSubnetName     = 4005
	DuplicateRouteTableName = 4006
	DuplicateAclName        = 4007
	InvalidDescription      = 4020
	PriorityConflict        = 4050
)

// QuotaExceed
const (
	VpcQuotaExceed               = 4401
	SubnetQuotaExceed            = 4402
	RouteTableQuotaExceed        = 4403
	RouteQuotaExceed             = 4404
	AclQuotaExceed               = 4405
	AclRuleQuotaExceed           = 4406
	SecurityGroupQuotaExceed     = 4407
	SecurityGroupRuleQuotaExceed = 4408
	FloatingipQuotaExceed        = 4409
)

// NotFound

const (
	VpcNotFound           = 5001
	SubnetNotFound        = 5002
	RouteTableNotFound    = 5003
	AclNotFound           = 5004
	AclRuleNotFound       = 5005
	HostNotFound          = 5006
	PortNotFound          = 5007
	SecurityGroupNotFound = 5008
	RouteNotFound         = 5009
)

// ResourceOperatefailed
const (
	VpcInUse                  = 5101
	SubnetInUse               = 5102
	RouteTableInUse           = 5103
	AclInUse                  = 5104
	SecurityGroupInUse        = 5105
	CountOfSgsExceededPerPort = 5106
)
