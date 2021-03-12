package command

import (
	"github.com/spf13/cobra"
	"analysis-server/cli/util"
	"analysis-server/sdk/options"
	"net"
	"strings"
)

func NewZbsProxyCommand(cmd *cobra.Command) {
	cmd.AddCommand(newZbsProxyAddCmd())
	cmd.AddCommand(newZbsProxyDeleteCmd())
	cmd.AddCommand(newZbsProxyListCmd())
}

func IPAddrCheck(ip string) bool {
	index := strings.Index(ip, ":")
	if index < 0 {
		return false
	}

	ipNew := net.ParseIP(ip[:index])
	if ipNew == nil {
		return false
	}

	return true
}

func newZbsProxyAddCmd() *cobra.Command {
	var opts options.ZbsProxyOptions
	cmd := &cobra.Command{
		Use:   "zbs-proxy-add [OPTIONS] Addr[IP:Port]",
		Short: "Add a zbs proxy",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 || !IPAddrCheck(args[0]) {
				cmd.Help()
				return
			}

			opts.Addr = args[0]

			if hv, err := Sdk.AddZbsProxy(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(hv)
			}
		},
	}
	return cmd
}

func newZbsProxyDeleteCmd() *cobra.Command {
	var opts options.ZbsProxyOptions
	cmd := &cobra.Command{
		Use:   "zbs-proxy-delete [OPTIONS] Addr[IP:Port]",
		Short: "Delete a proxy",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 || !IPAddrCheck(args[0]) {
				cmd.Help()
				return
			}

			opts.Addr = args[0]

			if hv, err := Sdk.DeleteZbsProxy(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(hv)
			}
		},
	}
	return cmd
}

func newZbsProxyListCmd() *cobra.Command {
	var opts options.ZbsProxyOptions
	cmd := &cobra.Command{
		Use:   "zbs-proxy-list [OPTIONS]",
		Short: "List proxy.",
		Run: func(cmd *cobra.Command, args []string) {
			if hv, err := Sdk.ListZbsProxy(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(hv)
			}
		},
	}
	return cmd
}
