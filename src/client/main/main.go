package main

import (
	//"bytes"
	"client/business"
	"client/util"
	"common/log"
	"encoding/binary"
	"fmt"
	//"net"
	"net/http"
	"os"
	//"runtime/debug"
	"strings"
	"time"
	"unsafe"
)

/*
#include <stdlib.h>
*/
import "C"

var (
	logger      *log.Logger
	quitCheckCh chan bool
	auth        *business.Authen
	profPort    string
)

func main() {

	//for test
	// iServerPort := 9999
	// pServerHost := "47.100.210.38"
	// var uTimeout uint64 = 30
	// pLogFileName := "service.log"
	// iLogFileSize := 2097152
	// iLogFileCount := 20
	// iRes := InitProxy(iServerPort, pServerHost, uTimeout, pLogFileName, iLogFileSize, iLogFileCount)
	// if iRes != 0 {
	// 	fmt.Printf("InitProxy,failed\r\n")
	// }
	// var pReqData *C.char = C.CString("abcdef")
	// var pOut *C.char = ProcessClientRequest(util.AccSubCreate, pReqData)
	// if pOut != nil {
	// 	fmt.Printf("ProcessClientRequest,pOut:%s\r\n", pOut)
	// 	FreeCBuf(pOut)
	// }
	// FreeCBuf(pReqData)
}

func initDebugPort() {
	http.HandleFunc("/setLogLevel", setLogLevel)
	//http.HandleFunc("/", response)
	//init profile server
	go func() {
		logger.LogInfo("profPort:", profPort)
		fmt.Println(http.ListenAndServe(":"+profPort, nil))
		os.Exit(0)
	}()
}

func initLogger(fileName string, fileMaxSize, filCount, iLogLevel int) (*log.Logger, error) {
	var h log.Handler
	var err error
	h, err = log.NewRotatingFileHandler(fileName, fileMaxSize, filCount)
	if err != nil {
		fmt.Printf("new log handler err: %v\r\n", err.Error())
		return nil, err
	}
	logger := log.NewDefault(h)
	logger.SetLevel(iLogLevel)
	return logger, nil
}

//export InitProxy
func InitProxy(iServerPort C.int, pServerHost *C.char, uTimeout C.int,
	pLogFileName *C.char, iLogFileSize, iLogFileCount C.int) C.int {
	strServerHost := C.GoString(pServerHost)
	logFileName := C.GoString(pLogFileName)
	serverPort := int(iServerPort)
	logFileSize := int(iLogFileSize)
	logFileCount := int(iLogFileCount)
	domain := fmt.Sprintf("http://%s:%d/analysis_server", strServerHost, serverPort)
	verbose := true
	var err error
	logger, err = initLogger(logFileName, logFileSize, logFileCount, log.LevelDebug)
	if err != nil {
		fmt.Printf("Init logger err: %v \r\n", err.Error())
		return -1
	}
	business.InitBusiness(logger, verbose, domain, uint64(uTimeout))
	profPort = "20000"
	initDebugPort()
	auth = new(business.Authen)
	quitCheckCh = make(chan bool)
	return 0
}

// func InitProxy(iServerPort int, strServerHost string, uTimeout uint64,
// 	logFileName string, logFileSize, logFileCount int) int {
// 	domain := fmt.Sprintf("http://%s:%d/analysis_server", strServerHost, iServerPort)
// 	verbose := true
// 	var err error
// 	logger, err = initLogger(logFileName, logFileSize, logFileCount, log.LevelDebug)
// 	if err != nil {
// 		fmt.Printf("Init logger err: %v \r\n", err.Error())
// 		return -1
// 	}
// 	business.InitBusiness(logger, verbose, domain, uint64(uTimeout))
// 	profPort = "20000"
// 	initDebugPort()
// 	auth = new(business.Authen)
// 	return 0
// }

//export GetUserStatus
func GetUserStatus() C.int {
	return C.int(auth.GetUserStatus())
}

func makeByteSlice(n int) []byte {
	p := C.malloc(C.size_t(n))
	return ((*[1 << 31]byte)(p))[0:n:n]
}

//export FreeByteSlice
func FreeByteSlice(p []byte) {
	C.free(unsafe.Pointer(&p[0]))
}

//export ProcessClientRequest
func ProcessClientRequest(iOpCode int, reqParamBuf []byte) []byte {
	var dataBuf []byte
	if isConvertToUtf8(iOpCode) {
		//logger.LogDebug("before convertion ,operation code:", iOpCode, "param data:", string(reqParamBuf))
		var err error
		if dataBuf, err = util.GBKToUTF8(reqParamBuf); err != nil {
			logger.LogError("covert gbk to utf8 failed")
			return nil
		}
		//logger.LogDebug("after convertion ,operation code:", iOpCode, "param data:", string(dataBuf))
	} else {
		dataBuf = make([]byte, len(reqParamBuf))
		copy(dataBuf, reqParamBuf)
	}
	var resultData []byte
	logger.LogDebug("ProcessClientRequest begin ,operation code:", iOpCode, "param data:", dataBuf)
	switch iOpCode {
	case util.QuitApp:
		errCode := quitApp()
		resultData = respOptResWithoutData(errCode)
		//quitCheckCh <- true
		close(quitCheckCh)
		time.Sleep(1 * time.Second)
		//os.Exit(0)
	case util.Heartbeat:
		resultData = respHeartbeatInfo()
	default:
		if iOpCode != util.UserLogin {
			if auth.GetUserStatus() != util.Online {
				resultData = respOptResWithoutData(util.ErrOffline)
				break
			}
		}
		switch {
		case iOpCode >= util.UserLogin && iOpCode <= util.OperatorUpdate:
			resultData = processOperator(iOpCode, dataBuf)
		case iOpCode >= util.CompanyCreate && iOpCode <= util.CompanyUpdate:
			resultData = processCompany(iOpCode, dataBuf)
		case iOpCode >= util.AccSubCreate && iOpCode <= util.AccSubUpdate:
			resultData = processAccSub(iOpCode, dataBuf)
		case iOpCode >= util.VoucherCreate && iOpCode <= util.VouTemplateList:
			resultData = processVoucher(iOpCode, dataBuf)
		case iOpCode >= util.YearBalanceCreate && iOpCode <= util.YearBalanceUpdate:
			resultData = processYearBalance(iOpCode, dataBuf)
		case iOpCode == util.MenuInfoList:
			resultData = processMenu(iOpCode, dataBuf)
		default:
			logger.LogError("opcode is mistake,the mistake operation code is: \r\n", iOpCode)
		}
	}
	if resultData == nil {
		return nil
	}
	logger.LogInfo("ProcessClientRequest finish ,operation code:", iOpCode, "result data:", resultData)
	//c malloc memory ,the caller free the memory
	cBuf := makeByteSlice(len(resultData))
	copy(cBuf, resultData)
	return cBuf
}

func processOperator(iOpCode int, dataBuf []byte) []byte {
	//var strRes string
	//var err error
	var optGate business.OperatorGateway
	switch iOpCode {
	case util.UserLogin:
		errCode := auth.UserLogin(dataBuf)
		if errCode == util.ErrNull {
			go onLineLoopCheck()
		}
		return respAuthResInfo(errCode)
	case util.LoginInfoList:
		resData, errCode := auth.ListLoginInfo(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.UserLogout:
		errCode := auth.Logout()
		quitCheckCh <- true
		return respAuthResInfo(errCode)
	case util.OperatorCreate:
		resData, errCode := optGate.CreateOperator(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.OperatorList:
		resData, errCode := optGate.ListOperatorInfo(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.OperatorShow:
		resData, errCode := optGate.GetOperatorInfo(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.OperatorDel:
		errCode := optGate.DeleteOperator(dataBuf)
		return respOptResWithoutData(errCode)
	case util.OperatorUpdate:
		errCode := optGate.UpdateOperator(dataBuf)
		return respOptResWithoutData(errCode)
	default:
		logger.LogError("opcode is mistake,the mistake operation code is: \r\n", iOpCode)
		panic("bug")
		//break
	}

	return nil
}

func processCompany(iOpCode int, dataBuf []byte) []byte {
	//var strRes string
	var comGate business.CompanyGateway
	switch iOpCode {
	case util.CompanyCreate:
		resData, errCode := comGate.CreateCompany(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.CompanyList:
		resData, errCode := comGate.ListCompany(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.CompanyDel:
		errCode := comGate.DeleteCompany(dataBuf)
		return respOptResWithoutData(errCode)
	case util.CompanyShow:
		resData, errCode := comGate.GetCompany(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.CompanyUpdate:
		errCode := comGate.UpdateCompany(dataBuf)
		return respOptResWithoutData(errCode)
	case util.InitResourceInfo:
		resData, errCode := comGate.InitResourceInfo(auth.OperatorID)
		return respOptResWithData(iOpCode, resData, errCode)
	default:
		logger.LogError("opcode is mistake,the mistake operation code is: \r\n", iOpCode)
		panic("bug")
	}
	return nil
}

func processAccSub(iOpCode int, dataBuf []byte) []byte {
	//var strRes string
	var accSubGate business.AccSubGateway
	switch iOpCode {
	case util.AccSubCreate:
		resData, errCode := accSubGate.CreateAccSub(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.AccSubReferenceQuery:
		resData, errCode := accSubGate.QueryAccSubReference(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.AccSubList:
		resData, errCode := accSubGate.ListAccSub(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.AccSubShow:
		resData, errCode := accSubGate.GetAccSub(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.AccSubDel:
		errCode := accSubGate.DeleteAccSub(dataBuf)
		return respOptResWithoutData(errCode)
	// case util.YearBalanceShow:
	// 	resData, errCode := accSubGate.GetYearBalance(dataBuf)
	// 	respOptResWithData(iOpCode, resData, errCode)
	// case util.YearBalanceUpdate:
	// 	errCode := accSubGate.UpdateYearBalance(dataBuf)
	// 	respOptResWithoutData(errCode)
	// case util.YearBalanceList:
	// 	resData, errCode := accSubGate.ListYearBalance(dataBuf)
	// 	respOptResWithData(iOpCode, resData, errCode)
	case util.CopyAccSubTemplate:
		resData, errCode := accSubGate.CopyAccSubTemplate(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.AccSubUpdate:
		errCode := accSubGate.UpdateAccSub(dataBuf)
		return respOptResWithoutData(errCode)
	default:
		logger.LogError("opcode is mistake,the mistake operation code is: \r\n", iOpCode)
		panic("bug")
	}
	return nil
}

func processVoucher(iOpCode int, dataBuf []byte) []byte {
	//var strRes string
	var voucherGate business.VoucherGateway
	switch iOpCode {
	case util.VoucherCreate:
		resData, errCode := voucherGate.CreateVoucher(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.VoucherUpdate:
		resData, errCode := voucherGate.UpdateVoucher(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.VoucherDel:
		errCode := voucherGate.DeleteVoucher(dataBuf)
		return respOptResWithoutData(errCode)
	case util.VoucherShow:
		resData, errCode := voucherGate.GetVoucher(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.VoucherArrange:
		errCode := voucherGate.ArrangeVoucher(dataBuf)
		return respOptResWithoutData(errCode)
	case util.VouInfoShow:
		resData, errCode := voucherGate.GetVoucherInfo(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.VouInfoList:
		resData, errCode := voucherGate.ListVoucherInfo(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.VouInfoListByMulCon:
		resData, errCode := voucherGate.ListVoucherInfoByMulCondition(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.VouInfoListLatest:
		resData, errCode := voucherGate.GetLatestVoucherInfo(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.VouInfoMaxNumOfMonth:
		resData, errCode := voucherGate.GetMaxNumOfMonth(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.BatchAuditVouchers:
		errCode := voucherGate.BatchAuditVouchers(dataBuf)
		return respOptResWithoutData(errCode)
	case util.VouInfoUpdate:
		errCode := voucherGate.UpdateVoucherInfo(dataBuf)
		return respOptResWithoutData(errCode)
	// case util.VouRecordCreate:
	// 	resData, errCode := voucherGate.CreateVoucherRecords(dataBuf)
	// 	return respOptResWithData(iOpCode, resData, errCode)
	// case util.VouRecordDel:
	// 	errCode := voucherGate.DeleteVoucherRecord(dataBuf)
	// 	return respOptResWithoutData(errCode)
	// case util.VouRecordsDel:
	// 	errCode := voucherGate.DeleteVoucherRecords(dataBuf)
	// 	return respOptResWithoutData(errCode)
	case util.VouRecordList:
		resData, errCode := voucherGate.ListVoucherRecords(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	// case util.VouRecordUpdate:
	// 	errCode := voucherGate.UpdateVoucherRecordByID(dataBuf)
	// 	return respOptResWithoutData(errCode)
	case util.CalculateAccuMoney:
		resData, errCode := voucherGate.CalculateAccumulativeMoney(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.BatchCalcAccuMoney:
		resData, errCode := voucherGate.BatchCalcAccuMoney(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.CalcAccountOfPeriod:
		resData, errCode := voucherGate.CalcAccountOfPeriod(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.VouTemplateCreate:
		resData, errCode := voucherGate.CreateVoucherTemplate(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.VouTemplateDel:
		errCode := voucherGate.DeleteVoucherTemplate(dataBuf)
		return respOptResWithoutData(errCode)
	case util.VouTemplateShow:
		resData, errCode := voucherGate.GetVoucherTemplate(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.VouTemplateList:
		resData, errCode := voucherGate.ListVoucherTemplate(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	default:
		logger.LogError("opcode is mistake,the mistake operation code is: \r\n", iOpCode)
		panic("bug")
	}
	return nil
}

func processMenu(iOpCode int, dataBuf []byte) []byte {
	//var strRes string
	var menuGate business.MenuInfoGateway
	switch iOpCode {
	case util.MenuInfoList:
		resData, errCode := menuGate.ListMenuInfo(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	default:
		logger.LogError("opcode is mistake,the mistake operation code is: \r\n", iOpCode)
		panic("bug")
	}
	return nil
}

func processYearBalance(iOpCode int, dataBuf []byte) []byte {
	//var strRes []byte
	var yearBalGate business.YearBalGateway
	switch iOpCode {
	case util.YearBalanceCreate:
		errCode := yearBalGate.CreateYearBalance(dataBuf)
		return respOptResWithoutData(errCode)
	case util.YearBalanceBatchCreate:
		errCode := yearBalGate.BatchCreateYearBalance(dataBuf)
		return respOptResWithoutData(errCode)
	case util.YearBalanceDel:
		errCode := yearBalGate.DeleteYearBalance(dataBuf)
		return respOptResWithoutData(errCode)
	case util.YearBalanceShow:
		resData, errCode := yearBalGate.GetYearBalance(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.YearBalanceList:
		resData, errCode := yearBalGate.ListYearBalance(dataBuf)
		return respOptResWithData(iOpCode, resData, errCode)
	case util.YearBalanceBatchUpdate:
		errCode := yearBalGate.BatchUpdateYearBalance(dataBuf)
		return respOptResWithoutData(errCode)
	case util.YearBalanceUpdate:
		errCode := yearBalGate.UpdateYearBalance(dataBuf)
		return respOptResWithoutData(errCode)
	default:
		logger.LogError("opcode is mistake,the mistake operation code is: \r\n", iOpCode)
		panic("bug")
	}
	return nil
}

func quitApp() int {
	//fmt.Println("quitApp,begin")
	logger.LogInfo("quitApp,begin")
	return util.ErrNull
}

//login/logout information;user errCode + status
func respAuthResInfo(errCode int) []byte {
	dataBuf := make([]byte, 8)
	binary.LittleEndian.PutUint32(dataBuf[0:4], uint32(errCode))
	binary.LittleEndian.PutUint32(dataBuf[4:], uint32(auth.GetUserStatus()))
	return dataBuf
}

//errCode
func respOptResWithoutData(errCode int) []byte {
	dataBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(dataBuf[0:4], uint32(errCode))
	return dataBuf
}

//errCode + data
func respOptResWithData(iOperationCode int, resData []byte, errCode int) []byte {
	var dataBuf []byte
	iOpCode := int(iOperationCode)
	if isConvertToGbk(iOpCode) {
		tmpBuf := make([]byte, 0)
		if errCode == util.ErrNull {
			var err error
			if tmpBuf, err = util.UTF8ToGBK(resData); err != nil {
				errCode = util.ErrUtf8ToGbkFailed
				tmpBuf = tmpBuf[0:0]
			}
		}
		iSize := int32(4 + len(tmpBuf))
		dataBuf = make([]byte, iSize)
		binary.LittleEndian.PutUint32(dataBuf[0:4], uint32(errCode))
		if errCode == util.ErrNull {
			copy(dataBuf[4:], tmpBuf)
		}
	} else {
		if errCode != util.ErrNull {
			resData = resData[0:0]
		}
		iSize := int32(4 + len(resData))
		dataBuf = make([]byte, iSize)
		binary.LittleEndian.PutUint32(dataBuf[0:4], uint32(errCode))
		if errCode == util.ErrNull {
			copy(dataBuf[4:], resData)
		}
	}
	return dataBuf
}

func isConvertToGbk(iOpCode int) bool {
	bRet := false
	switch iOpCode {
	case util.CompanyCreate:
		fallthrough
	case util.AccSubCreate:
		fallthrough
	case util.OperatorCreate:
	default:
		bRet = true
	}
	return bRet
}

func isConvertToUtf8(iOperationCode int) bool {
	bRet := false
	switch iOperationCode {
	case util.AccSubShow:
		fallthrough
	case util.AccSubDel:
		fallthrough
	case util.UserLogout:
		fallthrough
	case util.CompanyDel:
		fallthrough
	case util.CompanyShow:
		fallthrough
	case util.OperatorShow:
		fallthrough
	case util.OperatorDel:
		fallthrough
	case util.VoucherDel:
		fallthrough
	case util.VoucherShow:
		fallthrough
	// case util.VouRecordDel:
	// 	fallthrough
	case util.VouInfoShow:
		break
	case util.AccSubReferenceQuery:
		break
	case util.YearBalanceList:
		break
	case util.YearBalanceShow:
		break
	case util.YearBalanceUpdate:
		break
	case util.YearBalanceBatchUpdate:
		break
	case util.YearBalanceBatchCreate:
		break
	case util.YearBalanceCreate:
		break
	case util.YearBalanceDel:
		break
	default:
		bRet = true
		break
	}
	return bRet
}

//user status
func respHeartbeatInfo() []byte {
	dataBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(dataBuf[0:4], uint32(auth.GetUserStatus()))
	return dataBuf
}

func onLineLoopCheck() {
	logger.LogInfo("onLineLoopCheck,begin")
	for {
		select {
		case <-quitCheckCh:
			goto end
		case <-time.Tick(time.Second * 30):
			if auth.GetUserStatus() == util.Online {
				auth.OnlineCheck()
			} else {
				logger.LogInfo("It's going to quit onLineLoopCheck,because the user status is not online")
				goto end
			}
		}
	}
end:
	logger.LogInfo("onLineLoopCheck,end")
	return
}

//http://127.0.0.1:8888/setLogLevel?level=Info
func setLogLevel(w http.ResponseWriter, r *http.Request) {
	levelVal := r.FormValue("level")
	if len(levelVal) == 0 {
		w.Write([]byte("lack the log's level value"))
		return
	}
	logLevel := ""
	for index := 0; index < len(log.LevelName); index++ {
		if strings.Contains(levelVal, log.LevelName[index]) {
			logger.SetLevel(index)
			logLevel = log.LevelName[index]
		}
	}
	if logLevel == "" {
		w.Write([]byte("log's level value is invalid"))
		return
	}
	logLevelInfo := fmt.Sprintf("log's level has been to %s level", logLevel)
	w.Write([]byte(logLevelInfo))
	return
}
