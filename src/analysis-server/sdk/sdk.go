package sdk

import (
	"net/http"
	"time"

	"analysis-server/sdk/mgr"
	"analysis-server/sdk/util"
)

type CcSdk struct {
	Domain string
	//Tenant  string
	Verbose bool
	Admin   bool
	TraceId string
	Timeout uint64
	mgr.AccSub
	mgr.Company
	mgr.Operator
	mgr.Voucher
}

func (c *CcSdk) Setup() {
	util.Domain = c.Domain
	//util.Tenant = c.Tenant
	util.Verbose = c.Verbose
	util.Admin = c.Admin
	util.Client = new(http.Client)
	util.Client.Timeout = time.Duration(c.Timeout) * time.Millisecond
	util.TraceId = c.TraceId
}
