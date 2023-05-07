package utils

import (
	"bytes"
	"financeMgr/src/common/constant"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

const (
	MaxUint64 = ^uint64(0)
	MaxInt64  = int64(MaxUint64 >> 1)
)

func Assert(expr bool) {
	if !expr {
		panic("OOPS:")
	}
}

func DumpStack() string {
	buf := make([]byte, 4096)
	runtime.Stack(buf, true)
	return string(buf)
}

func Shell(bash string) (string, error) {
	var buf bytes.Buffer
	cmd := exec.Command("sh", "-c", bash)
	cmd.Stderr = &buf
	cmd.Stdout = &buf
	err := cmd.Run()
	out := buf.String()
	return out, err
}

func BgShell(bash string) (*exec.Cmd, error) {
	cmd := exec.Command("sh", "-c", bash)
	err := cmd.Start()
	return cmd, err
}

type ConditionFunc func() bool

func WaitUntil(condFunc ConditionFunc) {
	for {
		if condFunc() {
			break
		}
		time.Sleep(time.Millisecond * 50)
	}
}

func WaitUntilTimeout(condFunc ConditionFunc, timeout time.Duration) bool {
	start := time.Now()
	for {
		if condFunc() {
			return true
		} else if time.Now().Sub(start) > timeout {
			return false
		}

		time.Sleep(time.Millisecond * 50)
	}
}

func NewConn(ipaddr string) (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ipaddr)
	if err != nil {
		return nil, err
	}

	c, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}

	c.SetKeepAlive(true)
	c.SetLinger(10)
	c.SetNoDelay(true)

	return c, nil
}

func LBUrl(ipaddr, action string) string {
	url := "http://" + ipaddr + "/analysis-server?Action=" + action
	return url
}

func ServerUrl(apiAddr, action string) string {
	return fmt.Sprintf("http://%s/analysis-server?Action=%s", apiAddr, action)
}

func NewAlignedByteSlice(bufSize, alignSize int) []byte {
	data := make([]byte, bufSize+alignSize)
	alignedPos := calculateAlignedPos(data, alignSize)
	data = data[alignedPos : alignedPos+bufSize]
	return data
}

func calculateAlignedPos(tmpBuffer []byte, alignSize int) int {
	bufPointer := uintptr(unsafe.Pointer(&tmpBuffer[0]))
	alignMask := uintptr(alignSize - 1)
	return alignSize - int(bufPointer&alignMask)
}

func MIN_INT64(n1, n2 int64) int64 {
	if n1 >= n2 {
		return n2
	}
	return n1
}

func MIN_INT(n1, n2 int) int {
	if n1 >= n2 {
		return n2
	}
	return n1
}

func FillZero(buf []byte) {
	if buf == nil {
		return
	}

	for i := 0; i < len(buf); i++ {
		buf[i] = 0
	}
}

func IsFullZero(buf []byte) bool {
	for i := 0; i < len(buf); i++ {
		if buf[i] != 0 {
			return false
		}
	}
	return true
}

var r *rand.Rand = nil

func RandomDelaySeconds(maxSeconds int) {
	if r == nil {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	delayTime := time.Duration(r.Intn(maxSeconds))
	time.Sleep(time.Second * delayTime)
}

func RandomDelayMilliSeconds(maxMilliseconds int) {
	if r == nil {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	delayTime := time.Duration(r.Intn(maxMilliseconds))
	time.Sleep(time.Second * delayTime)
}

func CompareStringSlice(former, latter []string) ([]string, []string) {
	formerMap := make(map[string]bool, len(former))
	for _, s := range former {
		formerMap[s] = true
	}

	onlyInFormerSlice := make([]string, 0)
	onlyInLatterSlice := make([]string, 0)

	for _, s := range latter {
		if inFormer := formerMap[s]; inFormer {
			formerMap[s] = false
		} else {
			onlyInLatterSlice = append(onlyInLatterSlice, s)
		}
	}

	for s, onlyInFormer := range formerMap {
		if onlyInFormer { // True means s only exists in former slice.
			onlyInFormerSlice = append(onlyInFormerSlice, s)
		}
	}

	return onlyInFormerSlice, onlyInLatterSlice
}

func DialTcpByPort(startPort int32, endPort int32,
	dstAddr string, timeout time.Duration) (
	net.Conn, error) {

	if endPort <= startPort {
		return nil, constant.ERR_INVALIDARGUMENT
	}

	rand.Seed(time.Now().UnixNano())
	region := endPort - startPort + 1
	maxRetry := 4 * region
	retry := int32(0)

find:
	localPort := startPort + rand.Int31n(region)
	localAddr := ":" + strconv.Itoa(int(localPort))
	localTcpAddr, err := net.ResolveTCPAddr("tcp", localAddr)
	if err != nil {
		return nil, err
	}

	dial := &net.Dialer{
		Timeout:   timeout,
		LocalAddr: localTcpAddr,
	}

	conn, err := dial.Dial("tcp", dstAddr)
	if err != nil {
		if strings.Contains(err.Error(), syscall.EADDRINUSE.Error()) &&
			retry < maxRetry {
			time.Sleep(time.Millisecond * 2)
			retry++
			goto find
		}
	}
	return conn, err
}

// like `__FUNCTION__` of C.
func FUNC() string {
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	elems := strings.Split(fn.Name(), ".")
	return elems[len(elems)-1]
}

// like `__LINE__` of C.
func LINE() int {
	_, _, li, _ := runtime.Caller(1)
	return li
}

func getID(s []byte) int64 {
	s = s[len("goroutine "):]
	s = s[:bytes.IndexByte(s, ' ')]
	gid, _ := strconv.ParseInt(string(s), 10, 64)
	return gid
}

func GoRoutineID() int64 {
	var buf [64]byte
	return getID(buf[:runtime.Stack(buf[:], false)])
}

func PowerOfTwo(value int64) bool {
	bits := 0
	for value != 0 {
		value = value & (value - 1)
		bits++
	}
	return bits == 1
}

func GetShiftBits(value int64) int {
	if !PowerOfTwo(value) {
		Assert(false)
	}

	bits := 0
	for ; value != 1; value = value >> 1 {
		bits++
	}
	return bits
}

func IsProcessExist(pid int, subCmdline string) bool {
	cmdlinePath := fmt.Sprintf("/proc/%d/cmdline", pid)

	b, err := ioutil.ReadFile(cmdlinePath)
	if err != nil {
		return false
	}

	if strings.Contains(string(b), subCmdline) {
		return true
	}
	return false
}

// Find the first process whose cmdline contains subCmdline
func FindPID(subCmdline string) (int, error) {
	d, err := os.Open("/proc")
	if err != nil {
		return 0, err
	}
	defer d.Close()

	for {
		fis, err := d.Readdir(10)
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}

		for _, fi := range fis {
			// We only care about directories, since all pids are dirs
			if !fi.IsDir() {
				continue
			}

			// We only care if the name starts with a numeric
			name := fi.Name()
			if name[0] < '0' || name[0] > '9' {
				continue
			}

			pid, err := strconv.Atoi(name)
			if err != nil {
				continue
			}

			statPath := fmt.Sprintf("/proc/%d/cmdline", pid)
			b, err := ioutil.ReadFile(statPath)
			if err != nil {
				return 0, err
			}

			cmdLine := string(b)
			if strings.Contains(cmdLine, subCmdline) {
				return pid, nil
			}
		}
	}

	return 0, constant.ERR_NOTEXIST
}

func IsDeviceOpen(devName string) ([]string, error) {
	cmds := make([]string, 0) // process cmd
	if !strings.HasPrefix(devName, "/dev") {
		return cmds, constant.ERR_INVALIDARGUMENT
	}

	result, err := Shell("lsof " + devName)
	if err != nil {
		return cmds, nil
	}

	lines := strings.Split(result, "\n")
	for _, line := range lines {
		if strings.Contains(line, devName) {
			elems := strings.Split(line, " ")
			cmd := strings.TrimSpace(elems[0])
			cmds = append(cmds, cmd)
		}
	}

	return cmds, nil
}

func MemsetUint32(a []uint32, v uint32) {
	if len(a) == 0 {
		return
	}
	a[0] = v
	for bp := 1; bp < len(a); bp *= 2 {
		copy(a[bp:], a[:bp])
	}
}

func SetLimit() error {
	// var rLimit syscall.Rlimit
	// rLimit.Max = constant.MAX_OPEN_FD_LIMIT
	// rLimit.Cur = constant.MAX_OPEN_FD_LIMIT
	// err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	// if err != nil {
	// 	return err
	// }
	return nil
}
