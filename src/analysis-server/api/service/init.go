package service

import (
	"database/sql"
	"sync"
	"time"

	dbp "financeMgr/src/analysis-server/api/db"
	// "financeMgr/src/analysis-server/model"
	"financeMgr/src/common/log"
)

var gIdInfoService = NewIDInfoService()

var (
	quitCh                chan bool
	gCreateVoucherTableCh chan int
	gLogger               *log.Logger
	gDb                   *sql.DB
)

func InitService(logger *log.Logger, db *sql.DB) CcError {
	gLogger = logger
	gDb = db
	quitCh = make(chan bool)
	gCreateVoucherTableCh = make(chan int, 5)
	idInfoDao := &dbp.IDInfoDao{}
	gIdInfoService.InitIdInfoService(idInfoDao)
	return gIdInfoService.InitIdResource()
}

func ReleaseResourceInService() {
	close(gCreateVoucherTableCh)
}

// 写一个个定时器，每隔interval将当前的id资源写入数据库，防止服务异常退出导致id资源丢失
func StartIdResourcePersistence(interval time.Duration, ws *sync.WaitGroup) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case iYear := <-gCreateVoucherTableCh:
			err := CreateYearVoucherTable(iYear)
			if err != nil {
				gLogger.LogError("CreateYearVoucherTable,ErrInfo:", err.Error())
			}
		case <-quitCh:
			ccErr := gIdInfoService.WriteIdResourceToDb()
			if ccErr != nil {
				gLogger.LogError("before quit,WriteIdResourceToDb,ErrInfo:", ccErr.Error())
			}
			gLogger.LogInfo("StartIdResourcePersistence receive quit signal, exit StartIdResourcePersistence goroutine")
			goto end
		case <-ticker.C:
			ccErr := gIdInfoService.WriteIdResourceToDb()
			if ccErr != nil {
				gLogger.LogError("WriteIdResourceToDb,it is twice to fail,ErrInfo:", ccErr.Error())
			}
		}
	}
end:
	ws.Done()
}

func StopIdResourcePersistence() {
	close(quitCh)
}
