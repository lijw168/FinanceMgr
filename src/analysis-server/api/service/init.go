package service

import "time"

var GIdInfoService = NewIDInfoService()

var (
	quit chan bool
)

// 写一个个定时器，每隔interval将当前的id资源写入数据库，防止服务异常退出导致id资源丢失
func StartIdResourcePersistence(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-quit:
				ccErr := GIdInfoService.WriteIdResourceToDb()
				if ccErr != nil {
					GIdInfoService.logger.LogError("before quit,WriteIdResourceToDb,ErrInfo:", ccErr.Error())
				}
				GIdInfoService.logger.LogInfo("StartIdResourcePersistence receive quit signal, exit StartIdResourcePersistence goroutine")
				return
			case <-ticker.C:
				ccErr := GIdInfoService.WriteIdResourceToDb()
				if ccErr != nil {
					GIdInfoService.logger.LogError("WriteIdResourceToDb,it is twice to fail,ErrInfo:", ccErr.Error())
				}
			}
		}
	}()
}

func StopIdResourcePersistence() {
	quit <- true
}
