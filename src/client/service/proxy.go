package service

import (
	//"bytes"
	"client/business"
	"client/util"
	"common/log"
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

type Proxy struct {
	iListenPort  int
	listenCon    net.Listener
	quitCheckCh  chan bool
	processResWg sync.WaitGroup
	logger       *log.Logger
	profPort     string
	auth         *business.Authen
	//tr           *http.Transport
}

//Init ...
func (proxy *Proxy) Init(iListenPort, iServerPort int, strServerHost string, uTimeout uint64, logger *log.Logger) {
	proxy.iListenPort = iListenPort
	domain := fmt.Sprintf("http://%s:%d/analysis_server", strServerHost, iServerPort)
	verbose := true
	business.InitBusiness(logger, verbose, domain, uTimeout)
	proxy.logger = logger
	proxy.quitCheckCh = make(chan bool, 1)
	proxy.listenCon = nil
	proxy.profPort = "20000"
	proxy.initDebugPort()
	//proxy.tr = proxy.generateTr()
	proxy.auth = new(business.Authen)
}

func (proxy *Proxy) initDebugPort() {
	http.HandleFunc("/setLogLevel", proxy.setLogLevel)
	//http.HandleFunc("/", proxy.response)
	//init profile server
	go func() {
		proxy.logger.LogInfo("profPort:", proxy.profPort)
		fmt.Println(http.ListenAndServe(":"+proxy.profPort, nil))
		os.Exit(0)
	}()
}

func (proxy *Proxy) StartTcpService() {
	address := fmt.Sprintf("127.0.0.1:%d", proxy.iListenPort)
	var err error
	proxy.listenCon, err = net.Listen("tcp", address)
	if err != nil {
		proxy.logger.LogError("listen error:", err.Error())
		return
	}

	for {
		conn, err := proxy.listenCon.Accept()
		if err != nil {
			proxy.logger.LogError("accept error:", err.Error())
			break
		}
		// start a new goroutine to handle
		// the new connection.
		go proxy.handleConn(conn)
	}
	if proxy.listenCon != nil {
		proxy.listenCon.Close()
	}
}

// func (proxy *Proxy) stopService() {
// 	proxy.processResWg.Wait()
// }

func (proxy *Proxy) handleConn(conn net.Conn) {
	//defer proxy.processResWg.Done()
	defer func() {
		if err := recover(); err != nil {
			proxy.logger.LogError("in function handleConn [panic: %v, stack %s]", err, string(debug.Stack()))
			panic(err)
		}
	}()
	defer conn.Close()
	var err error
	errCode := util.ErrNull
	for {
		//fmt.Printf("begin,receive data\r\n")
		proxy.logger.LogDebug("begin,receive data")
		pk := NewPacket()
		err = pk.ReadFromConn(conn)
		if err != nil {
			proxy.logger.LogError("ReadFromConn,failed,err:", err.Error())
			break
		}
		proxy.logger.LogDebug("receive data :", pk)
		switch pk.OpCode {
		case util.QuitApp:
			errCode = proxy.quitApp(pk)
			proxy.respOptResWithoutData(conn, pk, errCode)
			//proxy.quitCheckCh <- true
			close(proxy.quitCheckCh)
			proxy.listenCon.Close()
			proxy.listenCon = nil
			time.Sleep(3 * time.Second)
			os.Exit(0)
		case util.Heartbeat:
			proxy.respHeartbeatInfo(conn, pk)
		default:
			if pk.OpCode != util.UserLogin {
				if proxy.auth.GetUserStatus() != util.Online {
					proxy.respOptResWithoutData(conn, pk, util.ErrOffline)
					break
				}
			}
			switch {
			case pk.OpCode >= util.UserLogin && pk.OpCode <= util.OperatorUpdate:
				proxy.processOperator(conn, pk)
			case pk.OpCode >= util.CompanyCreate && pk.OpCode <= util.CompanyUpdate:
				proxy.processCompany(conn, pk)
			case pk.OpCode >= util.AccSubCreate && pk.OpCode <= util.AccSubUpdate:
				proxy.processAccSub(conn, pk)
			case pk.OpCode >= util.VoucherCreate && pk.OpCode <= util.VouTemplateList:
				proxy.processVoucher(conn, pk)
			// case pk.OpCode >= util.YearBalanceCreate && pk.OpCode <= util.YearBalanceUpdate:
			// 	proxy.processYearBalance(conn, pk)
			case pk.OpCode == util.MenuInfoList:
				proxy.processMenu(conn, pk)
			default:
				proxy.logger.LogError("opcode is mistake,the mistake operation code is: \r\n", pk.OpCode)
			}
		}
		proxy.logger.LogInfo("response data :", pk)
	}
}

func (proxy *Proxy) processOperator(conn net.Conn, reqPk *Packet) {
	//var err error
	var optGate business.OperatorGateway
	switch reqPk.OpCode {
	case util.UserLogin:
		errCode := proxy.auth.UserLogin(reqPk.Buf)
		proxy.respAuthResInfo(conn, reqPk, errCode)
		go proxy.onLineLoopCheck()
	case util.LoginInfoList:
		resData, errCode := proxy.auth.ListLoginInfo(reqPk.Buf)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	case util.UserLogout:
		errCode := proxy.auth.Logout()
		proxy.respAuthResInfo(conn, reqPk, errCode)
		proxy.quitCheckCh <- true
	case util.OperatorCreate:
		resData, errCode := optGate.CreateOperator(reqPk.Buf)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	case util.OperatorList:
		resData, errCode := optGate.ListOperatorInfo(reqPk.Buf)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	case util.OperatorShow:
		resData, errCode := optGate.GetOperatorInfo(reqPk.Buf)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	case util.OperatorDel:
		errCode := optGate.DeleteOperator(reqPk.Buf)
		proxy.respOptResWithoutData(conn, reqPk, errCode)
	case util.OperatorUpdate:
		errCode := optGate.UpdateOperator(reqPk.Buf)
		proxy.respOptResWithoutData(conn, reqPk, errCode)
	default:
		proxy.logger.LogError("opcode is mistake,the mistake operation code is: \r\n", reqPk.OpCode)
		panic("bug")
		//break
	}
	return
}

func (proxy *Proxy) processCompany(conn net.Conn, reqPk *Packet) {
	var comGate business.CompanyGateway
	switch reqPk.OpCode {
	case util.CompanyCreate:
		resData, errCode := comGate.CreateCompany(reqPk.Buf)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	case util.CompanyList:
		resData, errCode := comGate.ListCompany(reqPk.Buf)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	case util.CompanyDel:
		errCode := comGate.DeleteCompany(reqPk.Buf)
		proxy.respOptResWithoutData(conn, reqPk, errCode)
	case util.CompanyShow:
		resData, errCode := comGate.GetCompany(reqPk.Buf)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	case util.CompanyUpdate:
		errCode := comGate.UpdateCompany(reqPk.Buf)
		proxy.respOptResWithoutData(conn, reqPk, errCode)
	case util.InitResourceInfo:
		resData, errCode := comGate.InitResourceInfo(proxy.auth.OperatorID)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	default:
		proxy.logger.LogError("opcode is mistake,the mistake operation code is: \r\n", reqPk.OpCode)
		panic("bug")
	}
	return
}

func (proxy *Proxy) processAccSub(conn net.Conn, reqPk *Packet) {
	var accSubGate business.AccSubGateway
	switch reqPk.OpCode {
	case util.AccSubCreate:
		resData, errCode := accSubGate.CreateAccSub(reqPk.Buf)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	case util.AccSubReferenceQuery:
		resData, errCode := accSubGate.QueryAccSubReference(reqPk.Buf)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	case util.AccSubList:
		resData, errCode := accSubGate.ListAccSub(reqPk.Buf)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	case util.AccSubShow:
		resData, errCode := accSubGate.GetAccSub(reqPk.Buf)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	case util.AccSubDel:
		errCode := accSubGate.DeleteAccSub(reqPk.Buf)
		proxy.respOptResWithoutData(conn, reqPk, errCode)
	case util.YearBalanceShow:
		resData, errCode := accSubGate.GetYearBalance(reqPk.Buf)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	case util.YearBalanceUpdate:
		errCode := accSubGate.UpdateYearBalance(reqPk.Buf)
		proxy.respOptResWithoutData(conn, reqPk, errCode)
	case util.YearBalanceList:
		resData, errCode := accSubGate.ListYearBalance(reqPk.Buf)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	case util.AccSubUpdate:
		errCode := accSubGate.UpdateAccSub(reqPk.Buf)
		proxy.respOptResWithoutData(conn, reqPk, errCode)
	default:
		proxy.logger.LogError("opcode is mistake,the mistake operation code is: \r\n", reqPk.OpCode)
		panic("bug")
	}
	return
}

func (proxy *Proxy) processVoucher(conn net.Conn, reqPk *Packet) {
	var voucherGate business.VoucherGateway
	switch reqPk.OpCode {
	case util.VoucherCreate:
		resData, errCode := voucherGate.CreateVoucher(reqPk.Buf)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	case util.VoucherUpdate:
		resData, errCode := voucherGate.UpdateVoucher(reqPk.Buf)
		if resData == nil {
			proxy.respOptResWithoutData(conn, reqPk, errCode)
		} else {
			proxy.respOptResWithData(conn, reqPk, errCode, resData)
		}
	case util.VoucherDel:
		errCode := voucherGate.DeleteVoucher(reqPk.Buf)
		proxy.respOptResWithoutData(conn, reqPk, errCode)
	case util.VoucherShow:
		resData, errCode := voucherGate.GetVoucher(reqPk.Buf)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	case util.VoucherArrange:
		errCode := voucherGate.ArrangeVoucher(reqPk.Buf)
		proxy.respOptResWithoutData(conn, reqPk, errCode)
	case util.VouInfoShow:
		resData, errCode := voucherGate.GetVoucherInfo(reqPk.Buf)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	case util.VouInfoList:
		resData, errCode := voucherGate.ListVoucherInfo(reqPk.Buf)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	case util.VouInfoListByMulCon:
		resData, errCode := voucherGate.ListVoucherInfoByMulCondition(reqPk.Buf)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	case util.VouInfoListLatest:
		resData, errCode := voucherGate.GetLatestVoucherInfo(reqPk.Buf)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	case util.VouInfoMaxNumOfMonth:
		resData, errCode := voucherGate.GetMaxNumOfMonth(reqPk.Buf)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	case util.BatchAuditVouchers:
		errCode := voucherGate.BatchAuditVouchers(reqPk.Buf)
		proxy.respOptResWithoutData(conn, reqPk, errCode)
	case util.VouInfoUpdate:
		errCode := voucherGate.UpdateVoucherInfo(reqPk.Buf)
		proxy.respOptResWithoutData(conn, reqPk, errCode)
	// case util.VouRecordCreate:
	// 	resData, errCode := voucherGate.CreateVoucherRecords(reqPk.Buf)
	// 	proxy.respOptResWithData(conn, reqPk, errCode, resData)
	// case util.VouRecordDel:
	// 	errCode := voucherGate.DeleteVoucherRecord(reqPk.Buf)
	// 	proxy.respOptResWithoutData(conn, reqPk, errCode)
	// case util.VouRecordsDel:
	// 	errCode := voucherGate.DeleteVoucherRecords(reqPk.Buf)
	// 	proxy.respOptResWithoutData(conn, reqPk, errCode)
	case util.VouRecordList:
		resData, errCode := voucherGate.ListVoucherRecords(reqPk.Buf)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	// case util.VouRecordUpdate:
	// 	errCode := voucherGate.UpdateVoucherRecordByID(reqPk.Buf)
	// 	proxy.respOptResWithoutData(conn, reqPk, errCode)
	case util.CalculateAccuMoney:
		resData, errCode := voucherGate.CalculateAccumulativeMoney(reqPk.Buf)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	case util.VouTemplateCreate:
		resData, errCode := voucherGate.CreateVoucherTemplate(reqPk.Buf)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	case util.VouTemplateDel:
		errCode := voucherGate.DeleteVoucherTemplate(reqPk.Buf)
		proxy.respOptResWithoutData(conn, reqPk, errCode)
	case util.VouTemplateShow:
		resData, errCode := voucherGate.GetVoucherTemplate(reqPk.Buf)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	case util.VouTemplateList:
		resData, errCode := voucherGate.ListVoucherTemplate(reqPk.Buf)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	default:
		proxy.logger.LogError("opcode is mistake,the mistake operation code is: \r\n", reqPk.OpCode)
		panic("bug")
	}
	return
}

func (proxy *Proxy) processMenu(conn net.Conn, reqPk *Packet) {
	var menuGate business.MenuInfoGateway
	switch reqPk.OpCode {
	case util.MenuInfoList:
		resData, errCode := menuGate.ListMenuInfo(reqPk.Buf)
		proxy.respOptResWithData(conn, reqPk, errCode, resData)
	default:
		proxy.logger.LogError("opcode is mistake,the mistake operation code is: \r\n", reqPk.OpCode)
		panic("bug")
	}
	return
}

// func (proxy *Proxy) processYearBalance(conn net.Conn, reqPk *Packet) {
// 	var yearBalGate business.YearBalGateway
// 	switch reqPk.OpCode {
// 	case util.YearBalanceCreate:
// 		errCode := yearBalGate.CreateYearBalance(reqPk.Buf)
// 		proxy.respOptResWithoutData(conn, reqPk, errCode)
// 	case util.YearBalanceDel:
// 		errCode := yearBalGate.DeleteYearBalanceByID(reqPk.Buf)
// 		proxy.respOptResWithoutData(conn, reqPk, errCode)
// 	case util.YearBalanceShow:
// 		resData, errCode := yearBalGate.GetYearBalanceById(reqPk.Buf)
// 		proxy.respOptResWithData(conn, reqPk, errCode, resData)
// 	case util.YearBalanceUpdate:
// 		errCode := yearBalGate.UpdateYearBalanceById(reqPk.Buf)
// 		proxy.respOptResWithoutData(conn, reqPk, errCode)
// 	default:
// 		proxy.logger.LogError("opcode is mistake,the mistake operation code is: \r\n", reqPk.OpCode)
// 		panic("bug")
// 	}
// 	return
// }

func (proxy *Proxy) quitApp(pk *Packet) int {
	//fmt.Println("quitApp,begin")
	proxy.logger.LogInfo("quitApp,begin")
	return util.ErrNull
}

//login/logout information;user status + errCode
func (proxy *Proxy) respAuthResInfo(conn net.Conn, reqPk *Packet, errCode int) (err error) {
	reqPk.Size = 8
	reqPk.Buf = reqPk.Buf[0:0]
	reqPk.Buf = make([]byte, 8)
	binary.LittleEndian.PutUint32(reqPk.Buf[0:4], uint32(errCode))
	binary.LittleEndian.PutUint32(reqPk.Buf[4:], uint32(proxy.auth.GetUserStatus()))
	err = reqPk.WriteToConn(conn)
	if err != nil {
		proxy.logger.LogError("respAuthResInfo,failed,err:", err.Error())
	}
	return
}

//errCode
func (proxy *Proxy) respOptResWithoutData(conn net.Conn, reqPk *Packet, errCode int) (err error) {
	reqPk.Size = 4
	reqPk.Buf = reqPk.Buf[0:0]
	reqPk.Buf = make([]byte, 4)
	binary.LittleEndian.PutUint32(reqPk.Buf[0:4], uint32(errCode))
	err = reqPk.WriteToConn(conn)
	if err != nil {
		proxy.logger.LogError("respOptResWithoutData,failed,err:", err.Error())
	}
	return
}

func (proxy *Proxy) isConvertToGbk(reqPk *Packet) bool {
	bRet := false
	switch reqPk.OpCode {
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

//errCode + data
func (proxy *Proxy) respOptResWithData(conn net.Conn, reqPk *Packet, errCode int, resData []byte) (err error) {
	if proxy.isConvertToGbk(reqPk) {
		tmpBuf := make([]byte, 0)
		if errCode == util.ErrNull {
			if tmpBuf, err = util.UTF8ToGBK(resData); err != nil {
				errCode = util.ErrUtf8ToGbkFailed
				tmpBuf = tmpBuf[0:0]
			}
		}
		reqPk.Size = int32(4 + len(tmpBuf))
		reqPk.Buf = reqPk.Buf[0:0]
		reqPk.Buf = make([]byte, reqPk.Size)
		binary.LittleEndian.PutUint32(reqPk.Buf[0:4], uint32(errCode))
		if errCode == util.ErrNull {
			copy(reqPk.Buf[4:], tmpBuf)
		}
	} else {
		if errCode != util.ErrNull {
			resData = resData[0:0]
		}
		reqPk.Size = int32(4 + len(resData))
		reqPk.Buf = reqPk.Buf[0:0]
		reqPk.Buf = make([]byte, reqPk.Size)
		binary.LittleEndian.PutUint32(reqPk.Buf[0:4], uint32(errCode))
		if errCode == util.ErrNull {
			copy(reqPk.Buf[4:], resData)
		}
	}
	err = reqPk.WriteToConn(conn)
	if err != nil {
		proxy.logger.LogError("respQuitAppResInfo,failed,err:", err.Error())
	}
	return
}

//user status
func (proxy *Proxy) respHeartbeatInfo(conn net.Conn, reqPk *Packet) (err error) {
	reqPk.Size = 4
	reqPk.Buf = reqPk.Buf[0:0]
	reqPk.Buf = make([]byte, 4)
	binary.LittleEndian.PutUint32(reqPk.Buf[0:4], uint32(proxy.auth.GetUserStatus()))
	err = reqPk.WriteToConn(conn)
	if err != nil {
		proxy.logger.LogError("respHeartbeatInfo,failed,err:", err.Error())
	}
	return
}

func (proxy *Proxy) onLineLoopCheck() {
	proxy.logger.LogInfo("onLineLoopCheck,begin")
	for {
		select {
		case <-proxy.quitCheckCh:
			goto end
		case <-time.Tick(time.Second * 30):
			if proxy.auth.GetUserStatus() == util.Online {
				proxy.auth.OnlineCheck()
			} else {
				proxy.logger.LogInfo("It's going to quit onLineLoopCheck,because the user status is not online")
				goto end
			}
		}
	}
end:
	proxy.logger.LogInfo("onLineLoopCheck,end")
	return
}

//http://127.0.0.1:8888/setLogLevel?level=Info
func (proxy *Proxy) setLogLevel(w http.ResponseWriter, r *http.Request) {
	levelVal := r.FormValue("level")
	if len(levelVal) == 0 {
		w.Write([]byte("lack the log's level value"))
		return
	}
	logLevel := ""
	for index := 0; index < len(log.LevelName); index++ {
		if strings.Contains(levelVal, log.LevelName[index]) {
			proxy.logger.SetLevel(index)
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
