package command

import (
	"github.com/spf13/cobra"
	"analysis-server/cli/util"
	"analysis-server/sdk/options"
)

func NewRackCommand(cmd *cobra.Command) {
	cmd.AddCommand(newRackCreateCmd())
	cmd.AddCommand(newRackDeleteCmd())
	cmd.AddCommand(newRackListCmd())
	cmd.AddCommand(newRackShowCmd())
}

func newRackCreateCmd() *cobra.Command {
	var opts options.CreateRackOptions
	cmd := &cobra.Command{
		Use:   "rack-create Tag",
		Short: "Create a Rack",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}

			opts.Tag = args[0]

			if hv, err := Sdk.CreateRack(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(hv)
			}
		},
	}
	return cmd
}

func newRackDeleteCmd() *cobra.Command {
	var opts options.DeleteOptions
	cmd := &cobra.Command{
		Use:   "rack-delete ID",
		Short: "Delete a Rack",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}

			opts.Id = args[0]

			if hv, err := Sdk.DeleteRack(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(hv)
			}
		},
	}
	return cmd
}

func newRackListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rack-list [OPTIONS]",
		Short: "List Rack. Support Filter",
	}
	defCs := []string{"Id", "Tag"}
	columns := cmd.Flags().StringArrayP("column", "c", defCs, "Columns to display")
	tag := cmd.Flags().String("tag", "", "Rack tag")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		var opts options.ListOptions
		opts.Limit = -1
		opts.Offset = 0
		opts.Filter = make(map[string]interface{})
		if *tag != "" {
			opts.Filter["tag"] = *tag
		}

		if _, views, err := Sdk.DescribeRacks(&opts); err != nil {
			util.FormatErrorOutput(err)
		} else {
			util.FormatListOutput(*columns, views)
		}
	}
	return cmd
}

func newRackShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rack-show Rack",
		Short: "Show given Rack",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			var opts options.BaseOptions
			opts.Id = args[0]
			view, err := Sdk.DescribeRack(&opts)
			if err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(view)
			}
		},
	}
	return cmd
}
