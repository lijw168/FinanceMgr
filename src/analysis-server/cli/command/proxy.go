package command

import (
	"github.com/spf13/cobra"
	"analysis-server/cli/util"
	"analysis-server/sdk/options"
)

func NewProxyCommand(cmd *cobra.Command) {
	cmd.AddCommand(newProxyCreateCmd())
	cmd.AddCommand(newProxyDeleteCmd())
	cmd.AddCommand(newProxyListCmd())
	cmd.AddCommand(newProxyShowCmd())
}

func newProxyCreateCmd() *cobra.Command {
	var opts options.CreateProxyOptions
	cmd := &cobra.Command{
		Use:   "proxy-create [OPTIONS] Addr",
		Short: "Create a Proxy",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}

			opts.Addr = args[0]

			if hv, err := Sdk.CreateProxy(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(hv)
			}
		},
	}
	return cmd
}

func newProxyDeleteCmd() *cobra.Command {
	return deleteCmd(RESOURCE_TYPE_PROXY, Sdk.DeleteProxy)
}

func newProxyListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proxy-list [OPTIONS]",
		Short: "List Proxy. Support Filter",
	}
	defCs := []string{"Id", "Addr"}
	columns := cmd.Flags().StringArrayP("column", "c", defCs, "Columns to display")
	Addr := cmd.Flags().String("Addr", "", "Proxy Addr")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		var opts options.ListOptions
		opts.Limit = -1
		opts.Offset = 0
		opts.Filter = make(map[string]interface{})
		if *Addr != "" {
			opts.Filter["Addr"] = *Addr
		}

		if _, views, err := Sdk.DescribeProxys(&opts); err != nil {
			util.FormatErrorOutput(err)
		} else {
			util.FormatListOutput(*columns, views)
		}
	}
	return cmd
}

func newProxyShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proxy-show Proxy",
		Short: "Show given Proxy",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			var opts options.BaseOptions
			opts.Id = args[0]
			view, err := Sdk.DescribeProxy(&opts)
			if err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(view)
			}
		},
	}
	return cmd
}
