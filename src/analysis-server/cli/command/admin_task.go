package command

import (
	"github.com/spf13/cobra"
	"analysis-server/cli/util"
	"analysis-server/sdk/options"
)

func NewAtCommand(cmd *cobra.Command) {
	cmd.AddCommand(newAtShowCmd())
}
func newAtShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "at-show At",
		Short: "Show given At",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			var opts options.BaseOptions
			opts.Id = args[0]
			view, err := Sdk.DescribeAt(&opts)
			if err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(view)
			}
		},
	}
	return cmd
}
