package db

import (
	"financeMgr/src/common/log"
)

var (
	gLogger log.ILog
)

func InitDao(logger *log.Logger) {
	gLogger = logger
}
