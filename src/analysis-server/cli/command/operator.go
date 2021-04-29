package command

import (
	//"errors"
	"fmt"
	// "strings"
	// "time"

	"analysis-server/cli/util"
	//"analysis-server/model"
	"analysis-server/sdk/options"
	//cons "common/constant"
	"github.com/spf13/cobra"
	"strconv"
)

func NewOperatorCommand(cmd *cobra.Command) {
	cmd.AddCommand(newOperatorCreateCmd())
	cmd.AddCommand(newOperatorDeleteCmd())
	cmd.AddCommand(newOperatorListCmd())
	cmd.AddCommand(newOperatorShowCmd())
	cmd.AddCommand(newOperatorUpdateCmd())
}

func newOperatorCreateCmd() *cobra.Command {
	var opts options.CreateOptInfoOptions
	cmd := &cobra.Command{
		Use:   "operator-create [OPTIONS] companyID name password",
		Short: "Create a operator",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			if id, err := strconv.Atoi(args[0]); err != nil {
				fmt.Println("change to int fail", args[0])
			} else {
				opts.CompanyID = id
			}
			opts.Name = args[1]
			opts.Password = args[2]

			if view, err := Sdk.CreateOperator(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(view)
			}
		},
	}
	cmd.Flags().StringVar(&opts.Department, "department", "test", "department")
	cmd.Flags().StringVar(&opts.Job, "job", "test", "job")
	cmd.Flags().IntVar(&opts.Role, "role", 1, "role")
	cmd.Flags().IntVar(&opts.Status, "status", 1, "status")
	return cmd
}

func newOperatorDeleteCmd() *cobra.Command {
	var opts options.NameOptions
	cmd := &cobra.Command{
		Use:   "operator-delete [OPTIONS] name",
		Short: "Delete a operator",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			opts.Name = args[0]
			if err := Sdk.DeleteOperator(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	return cmd
}

func newOperatorListCmd() *cobra.Command {
	defCs := []string{"CompanyID", "Name", "Password", "Job", "Department", "Status", "Role"}
	cmd := &cobra.Command{
		Use:   "operator-list ",
		Short: "List operators Support Filter",
	}
	columns := cmd.Flags().StringArrayP("column", "c", defCs, "Columns to display")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		var opts options.ListOptions
		opts.Limit = -1
		opts.Offset = 0
		//for test
		//opts.Filter = make(map[string]interface{})
		//opts.Filter["status"] = "creating|available|in-use"
		if _, views, err := Sdk.ListOperatorInfo(&opts); err != nil {
			util.FormatErrorOutput(err)
		} else {
			util.FormatListOutput(*columns, views)
		}
	}
	return cmd
}

func newOperatorShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "operator-show [OPTIONS] operatorName",
		Short: "Show a operator information",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			var opts options.NameOptions
			opts.Name = args[0]
			view, err := Sdk.GetOperatorInfo(&opts)
			if err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(view)
			}
		},
	}
	return cmd
}

func newOperatorUpdateCmd() *cobra.Command {
	var opts options.ModifyOptInfoOptions
	cmd := &cobra.Command{
		Use:   "operator-update [OPTIONS] name job",
		Short: "update a operator",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			opts.Name = args[0]
			//for test
			opts.Job = args[1]
			if err := Sdk.UpdateOperator(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	return cmd
}
