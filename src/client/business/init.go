package business

import (
	"financeMgr/src/analysis-server/sdk"
	"financeMgr/src/common/log"
)

var (
	logger *log.Logger
	cSdk   *sdk.CcSdk
)

func InitBusiness(pLog *log.Logger, verbose bool, domain string, timeout uint64) {
	logger = pLog
	cSdk = &sdk.CcSdk{}
	cSdk.Verbose = verbose
	cSdk.Domain = domain
	cSdk.Timeout = timeout
	cSdk.Logger = logger
	cSdk.Setup()
}
