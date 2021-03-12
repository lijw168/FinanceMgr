package command

import (
	"github.com/spf13/cobra"
	"analysis-server/cli/util"
	"analysis-server/sdk/options"
)

func NewRgCommand(cmd *cobra.Command) {
	cmd.AddCommand(newRgListCmd())
	cmd.AddCommand(newRgShowCmd())
}

func newRgListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rg-list",
		Short: "List Replication Group Support Filter",
	}
	defCs := []string{"Id", "PoolId", "epoch", "version", "size"}
	columns := cmd.Flags().StringArrayP("column", "c", defCs, "Columns to display")
	poolId := cmd.Flags().String("pool_id", "", "PoolId")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		var opts options.ListOptions
		opts.Limit = -1
		opts.Offset = 0
		opts.Filter = make(map[string]interface{})
		if *poolId != "" {
			opts.Filter["poolId"] = *poolId
		}

		if _, views, err := Sdk.DescribeRgs(&opts); err != nil {
			util.FormatErrorOutput(err)
		} else {
			util.FormatListOutput(*columns, views)
		}
	}
	return cmd
}

func newRgShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rg-show Replication Group",
		Short: "Show given Replication Group",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			var opts options.BaseOptions
			opts.Id = args[0]
			view, err := Sdk.DescribeRg(&opts)
			if err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(view)
			}
		},
	}
	return cmd
}
