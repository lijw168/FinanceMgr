package main

import (
	//"bytes"
	"encoding/binary"
	"financeMgr/src/client/business"
	"financeMgr/src/client/util"
	"financeMgr/src/common/log"
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
		logger.LogDebug("before convertion ,operation code:", iOpCode, "param data:", string(reqParamBuf))
		var err error
		if dataBuf, err = util.GBKToUTF8(reqParamBuf); err != nil {
			logger.Error("covert gbk to utf8 failed,the err:%s", err.Error())
			return nil
		}
		logger.LogDebug("after convertion ,operation code:", iOpCode, "param data:", string(dataBuf))
	} else {
		dataBuf = make([]byte, len(reqParamBuf))
		copy(dataBuf, reqParamBuf)
	}
	var resultData []byte
	logger.LogDebug("ProcessClientRequest begin ,operation code:", iOpCode)
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
		if iOpCode != util.Login {
			if auth.GetUserStatus() != util.Online {
				resultData = respOptResWithoutData(util.ErrOffline)
				break
			}
		}
		switch {
		case iOpCode >= util.Login && iOpCode <= util.OperatorUpdate:
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
	logger.LogInfo("ProcessClientRequest finish ,operation code:", iOpCode)
	//c malloc memory ,the caller free the memory
	cBuf := makeByteSlice(len(resultData))
	copy(cBuf, resultData)
	return cBuf
}

func processOperator(iOpCode int, dataBuf []byte) []byte {
	var optGate business.OperatorGateway
	switch iOpCode {
	case util.Login:
		errCode, errMsg := auth.UserLogin(dataBuf)
		if errCode == util.ErrNull {
			go onLineLoopCheck()
		}
		return respAuthResInfo(errCode, errMsg)
	case util.LoginInfoList:
		resData, errCode, errMsg := auth.ListLoginInfo(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.Logout:
		errCode, errMsg := auth.Logout()
		quitCheckCh <- true
		return respAuthResInfo(errCode, errMsg)
	case util.OperatorCreate:
		resData, errCode, errMsg := optGate.CreateOperator(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.OperatorList:
		resData, errCode, errMsg := optGate.ListOperatorInfo(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.OperatorShow:
		resData, errCode, errMsg := optGate.GetOperatorInfo(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.OperatorDel:
		errCode, errMsg := optGate.DeleteOperator(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithoutData(errCode)
		}
	case util.OperatorUpdate:
		errCode, errMsg := optGate.UpdateOperator(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithoutData(errCode)
		}
	default:
		logger.LogError("opcode is mistake,the mistake operation code is: \r\n", iOpCode)
		panic("bug")
	}
}

func processCompany(iOpCode int, dataBuf []byte) []byte {
	var comGate business.CompanyGateway
	switch iOpCode {
	case util.CompanyCreate:
		resData, errCode, errMsg := comGate.CreateCompany(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.CompanyList:
		resData, errCode, errMsg := comGate.ListCompany(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.CompanyDel:
		errCode, errMsg := comGate.DeleteCompany(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithoutData(errCode)
		}
	case util.CompanyShow:
		resData, errCode, errMsg := comGate.GetCompany(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.CompanyUpdate:
		errCode, errMsg := comGate.UpdateCompany(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithoutData(errCode)
		}
	case util.InitResourceInfo:
		resData, errCode := comGate.InitResourceInfo(auth.OperatorID)
		return respOptResWithData(iOpCode, resData, errCode)
	default:
		logger.Error("opcode is mistake,the mistake operation code is: %d", iOpCode)
		panic("bug")
	}
}

func processAccSub(iOpCode int, dataBuf []byte) []byte {
	var accSubGate business.AccSubGateway
	switch iOpCode {
	case util.AccSubCreate:
		resData, errCode, errMsg := accSubGate.CreateAccSub(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.AccSubReferenceQuery:
		resData, errCode, errMsg := accSubGate.QueryAccSubReference(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.AccSubList:
		resData, errCode, errMsg := accSubGate.ListAccSub(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.AccSubShow:
		resData, errCode, errMsg := accSubGate.GetAccSub(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.AccSubDel:
		errCode, errMsg := accSubGate.DeleteAccSub(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithoutData(errCode)
		}
	case util.CopyAccSubTemplate:
		resData, errCode, errMsg := accSubGate.CopyAccSubTemplate(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.AccSubUpdate:
		errCode, errMsg := accSubGate.UpdateAccSub(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithoutData(errCode)
		}
	default:
		logger.Error("opcode is mistake,the mistake operation code is: %d", iOpCode)
		panic("bug")
	}
}

func processVoucher(iOpCode int, dataBuf []byte) []byte {
	var voucherGate business.VoucherGateway
	switch iOpCode {
	case util.VoucherCreate:
		resData, errCode, errMsg := voucherGate.CreateVoucher(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.VoucherUpdate:
		resData, errCode, errMsg := voucherGate.UpdateVoucher(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.VoucherDel:
		errCode, errMsg := voucherGate.DeleteVoucher(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithoutData(errCode)
		}
	case util.VoucherShow:
		resData, errCode, errMsg := voucherGate.GetVoucher(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.VoucherArrange:
		errCode, errMsg := voucherGate.ArrangeVoucher(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithoutData(errCode)
		}
	case util.VouInfoShow:
		resData, errCode, errMsg := voucherGate.GetVoucherInfo(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.VouInfoList:
		resData, errCode, errMsg := voucherGate.ListVoucherInfo(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.VouInfoListWithAuxCond:
		resData, errCode, errMsg := voucherGate.ListVoucherInfoWithAuxCondition(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.VouInfoListLatest:
		resData, errCode, errMsg := voucherGate.GetLatestVoucherInfo(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.VouInfoMaxNumOfMonth:
		resData, errCode, errMsg := voucherGate.GetMaxNumOfMonth(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.BatchAuditVouchers:
		errCode, errMsg := voucherGate.BatchAuditVouchers(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithoutData(errCode)
		}
	case util.VouInfoUpdate:
		errCode, errMsg := voucherGate.UpdateVoucherInfo(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithoutData(errCode)
		}
	case util.VoucherInfoNoAuditedShow:
		resData, errCode, errMsg := voucherGate.GetNoAuditedVoucherInfoCount(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
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
		resData, errCode, errMsg := voucherGate.ListVoucherRecords(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	// case util.VouRecordUpdate:
	// 	errCode := voucherGate.UpdateVoucherRecordByID(dataBuf)
	// 	return respOptResWithoutData(errCode)
	case util.CalculateAccuMoney:
		resData, errCode, errMsg := voucherGate.CalculateAccumulativeMoney(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.BatchCalcAccuMoney:
		resData, errCode, errMsg := voucherGate.BatchCalcAccuMoney(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.CalcAccountOfPeriod:
		resData, errCode, errMsg := voucherGate.CalcAccountOfPeriod(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.VouTemplateCreate:
		resData, errCode, errMsg := voucherGate.CreateVoucherTemplate(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.VouTemplateDel:
		errCode, errMsg := voucherGate.DeleteVoucherTemplate(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithoutData(errCode)
		}
	case util.VouTemplateShow:
		resData, errCode, errMsg := voucherGate.GetVoucherTemplate(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.VouTemplateList:
		resData, errCode, errMsg := voucherGate.ListVoucherTemplate(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	default:
		logger.Error("opcode is mistake,the mistake operation code is: %d", iOpCode)
		panic("bug")
	}
}

func processMenu(iOpCode int, dataBuf []byte) []byte {
	var menuGate business.MenuInfoGateway
	switch iOpCode {
	case util.MenuInfoList:
		resData, errCode, errMsg := menuGate.ListMenuInfo(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	default:
		logger.Error("opcode is mistake,the mistake operation code is: %d", iOpCode)
		panic("bug")
	}
}

func processYearBalance(iOpCode int, dataBuf []byte) []byte {
	//var strRes []byte
	var yearBalGate business.YearBalGateway
	switch iOpCode {
	case util.YearBalanceCreate:
		errCode, errMsg := yearBalGate.CreateYearBalance(dataBuf)
		if errCode == util.ErrNull {
			return respOptResWithoutData(errCode)
		} else {
			return respOptResWithErrMsg(errCode, errMsg)
		}
	case util.YearBalanceBatchCreate:
		errCode, errMsg := yearBalGate.BatchCreateYearBalance(dataBuf)
		if errCode == util.ErrNull {
			return respOptResWithoutData(errCode)
		} else {
			return respOptResWithErrMsg(errCode, errMsg)
		}
	case util.YearBalanceDel:
		errCode, errMsg := yearBalGate.DeleteYearBalance(dataBuf)
		if errCode == util.ErrNull {
			return respOptResWithoutData(errCode)
		} else {
			return respOptResWithErrMsg(errCode, errMsg)
		}
	case util.YearBalanceShow:
		resData, errCode, errMsg := yearBalGate.GetYearBalance(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.YearBalanceList:
		resData, errCode, errMsg := yearBalGate.ListYearBalance(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.AccSubYearBalValueShow:
		resData, errCode, errMsg := yearBalGate.GetAccSubYearBalValue(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	case util.AnnualClosing:
		errCode, errMsg := yearBalGate.AnnualClosing(dataBuf)
		if errCode == util.ErrNull {
			return respOptResWithoutData(errCode)
		} else {
			return respOptResWithErrMsg(errCode, errMsg)
		}
	case util.CancelAnnualClosing:
		errCode, errMsg := yearBalGate.CancelAnnualClosing(dataBuf)
		if errCode == util.ErrNull {
			return respOptResWithoutData(errCode)
		} else {
			return respOptResWithErrMsg(errCode, errMsg)
		}
	case util.BatchUpdateBals:
		errCode, errMsg := yearBalGate.BatchUpdateBals(dataBuf)
		if errCode == util.ErrNull {
			return respOptResWithoutData(errCode)
		} else {
			return respOptResWithErrMsg(errCode, errMsg)
		}
	case util.YearBalanceUpdate:
		errCode, errMsg := yearBalGate.UpdateYearBalance(dataBuf)
		if errCode == util.ErrNull {
			return respOptResWithoutData(errCode)
		} else {
			return respOptResWithErrMsg(errCode, errMsg)
		}
	case util.AnnualClosingStatusShow:
		resData, errCode, errMsg := yearBalGate.GetAnnualClosingStatus(dataBuf)
		if errCode != util.ErrNull {
			return respOptResWithErrMsg(errCode, errMsg)
		} else {
			return respOptResWithData(iOpCode, resData, errCode)
		}
	default:
		logger.Error("opcode is mistake,the mistake operation code is: %d\r\n", iOpCode)
		panic("bug")
	}
	//return nil
}

func quitApp() int {
	//fmt.Println("quitApp,begin")
	logger.LogInfo("quitApp,begin")
	return util.ErrNull
}

// login/logout information;user errCode + status
func respAuthResInfo(errCode int, errMsg string) []byte {
	dataBuf := make([]byte, 8+len(errMsg))
	binary.LittleEndian.PutUint32(dataBuf[0:4], uint32(errCode))
	binary.LittleEndian.PutUint32(dataBuf[4:], uint32(auth.GetUserStatus()))
	if errCode == util.ErrNull {
		return dataBuf
	}
	if tmpBuf, err := util.UTF8ToGBK([]byte(errMsg)); err != nil {
		errCode = util.ErrUtf8ToGbkFailed
		logger.Error("UTF8ToGBK failed,err:%s\r\n", err.Error())
		copy(dataBuf[8:], errMsg)
	} else {
		copy(dataBuf[8:], tmpBuf)
	}
	return dataBuf
}

// errCode
func respOptResWithoutData(errCode int) []byte {
	dataBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(dataBuf[0:4], uint32(errCode))
	return dataBuf
}

// errCode + data
func respOptResWithData(iOperationCode int, resData []byte, errCode int) []byte {
	var dataBuf []byte
	iOpCode := int(iOperationCode)
	if isConvertToGbk(iOpCode) {
		//tmpBuf := make([]byte, 0)
		// if errCode == util.ErrNull {
		// 	var err error
		// 	if tmpBuf, err = util.UTF8ToGBK(resData); err != nil {
		// 		logger.Error("UTF8ToGBK failed,err:%s\r\n", err.Error())
		// 		return respOptResWithErrMsg(util.ErrUtf8ToGbkFailed, err.Error())
		// 	}
		// }
		start := time.Now()
		defer func() {
			logger.Debug("[UTF8ToGBK,operationCode:%d] [SqlElapsed: %v]", iOperationCode, time.Since(start))
		}()
		var tmpBuf []byte
		var err error
		if tmpBuf, err = util.UTF8ToGBK(resData); err != nil {
			logger.Error("UTF8ToGBK failed,err:%s\r\n", err.Error())
			return respOptResWithErrMsg(util.ErrUtf8ToGbkFailed, err.Error())
		}
		iSize := int32(4 + len(tmpBuf))
		dataBuf = make([]byte, iSize)
		binary.LittleEndian.PutUint32(dataBuf[0:4], uint32(errCode))
		if errCode == util.ErrNull {
			copy(dataBuf[4:], tmpBuf)
		}
	} else {
		// if errCode != util.ErrNull {
		// 	resData = resData[0:0]
		// }
		iSize := int32(4 + len(resData))
		dataBuf = make([]byte, iSize)
		binary.LittleEndian.PutUint32(dataBuf[0:4], uint32(errCode))
		if errCode == util.ErrNull {
			copy(dataBuf[4:], resData)
		}
	}
	return dataBuf
}

// errCode + error message
func respOptResWithErrMsg(errCode int, errMsg string) []byte {
	dataBuf := make([]byte, 4+len(errMsg))
	binary.LittleEndian.PutUint32(dataBuf[0:4], uint32(errCode))
	if errCode == util.ErrNull {
		return dataBuf
	}
	if tmpBuf, err := util.UTF8ToGBK([]byte(errMsg)); err != nil {
		errCode = util.ErrUtf8ToGbkFailed
		logger.Error("UTF8ToGBK failed,err:%s\r\n", err.Error())
		copy(dataBuf[4:], errMsg)
	} else {
		copy(dataBuf[4:], tmpBuf)
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
	bRet := true
	switch iOperationCode {
	case util.Login, util.LoginInfoList:

	case util.OperatorCreate, util.OperatorList, util.OperatorUpdate:

	case util.CompanyCreate, util.CompanyList, util.CompanyUpdate:

	case util.AccSubCreate, util.AccSubList, util.CopyAccSubTemplate, util.AccSubUpdate:

	case util.VoucherCreate, util.VouInfoList, util.VouInfoListWithAuxCond, util.BatchAuditVouchers:
	case util.VoucherUpdate, util.VouInfoUpdate, util.VouRecordList:

	case util.VouTemplateCreate, util.VouTemplateList:

	case util.MenuInfoList:

	default:
		bRet = false
	}
	return bRet
}

// user status
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
	//return
}

// http://127.0.0.1:8888/setLogLevel?level=Info
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
	//return
}
