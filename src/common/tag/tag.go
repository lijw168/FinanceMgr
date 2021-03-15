package tag

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"path/filepath"
)

type GitInfo struct {
	GitVersion string
	BuildTime  string
}

const (
	HEADER  = "\033[95m"
	WARNING = "\033[93m"
	ENDC    = "\033[0m"
	FAIL    = "\033[91m"
)

const (
	DefaultVersion = "2000.01.01.release"
	DefaultTime    = "2000-01-01T00:00:00+0800"
)

var (
	flagSet        = flag.NewFlagSet("zbs", flag.ContinueOnError)
	version        = flagSet.Bool("tag", false, "git version and build time")
	GitVersion     = DefaultVersion
	BuildTime      = DefaultTime
	binaryFileName = filepath.Base(os.Args[0])
)

func (p *GitInfo) Server(dummy, reply *GitInfo) error {
	reply.GitVersion = GitVersion
	reply.BuildTime = BuildTime
	return nil
}

func domainsocket() string {
	binaryName := os.TempDir() + "/" + binaryFileName + ".sock"
	return binaryName
}

func server() {
	gitInfo := &GitInfo{}
	rpcServer := rpc.NewServer()
	err := rpcServer.Register(gitInfo)
	if err != nil {
		panic(err)
	}

	os.Remove(domainsocket())
	l, err := net.Listen("unix", domainsocket())
	if err != nil {
		fail(err.Error())
		return
	}

	go rpcServer.Accept(l)
}

func client() (*GitInfo, error) {
	client, err := rpc.Dial("unix", domainsocket())
	if err != nil {
		return nil, err
	}
	reply := &GitInfo{}
	err = client.Call("GitInfo.Server", reply, reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func info(body string) {
	fmt.Println(body + ENDC)
}

func header(header string) {
	info(HEADER + header)
}

func warning(body string) {
	info(WARNING + body)
}

func fail(body string) {
	info(FAIL + body)
}

/*
func init() {
	flag.Bool("tag", false, "git version and build time")
	flagSet.SetOutput(ioutil.Discard)
	flagSet.Parse(os.Args[1:])

	if *version {
		func() {
			header("---------------默认的git版本信息-------------------")
			info("git version:" + DefaultVersion)
			info("build time:" + DefaultTime)
			if DefaultVersion == GitVersion && DefaultTime == BuildTime {
				fail(binaryFileName + "二进制文件没有烧录git版本信息")
			}
			header("---------" + binaryFileName + "二进制文件的git版本信息--------------")
			info("git version:" + GitVersion)
			info("build time:" + BuildTime)

			header("---------" + binaryFileName + "正在执行的进程的git版本信息---------")
			runtimeGitInfo, err := client()
			if err != nil {
				fail("获取正在执行的进程的git版本信息, 错误信息:" + err.Error())
				return
			}
			info("git version:" + runtimeGitInfo.GitVersion)
			info("build time:" + runtimeGitInfo.BuildTime)
			if GitVersion != runtimeGitInfo.GitVersion ||
				BuildTime != runtimeGitInfo.BuildTime {
				warning("warning: 正在执行的进程的版本信息与二进制文件的版本信息不一致")
			}
		}()

		os.Exit(0)
	}
	server()
}
*/
