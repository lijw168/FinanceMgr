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
		Use:   "operator-create [OPTIONS] companyID name password,job,department,role",
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
			opts.Job = args[3]
			opts.Department = args[4]
			if role, err := strconv.Atoi(args[5]); err != nil {
				fmt.Println("change to int fail", args[5])
			} else {
				opts.Role = role
			}

			if view, err := Sdk.CreateOperator(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(view)
			}
		},
	}
	cmd.Flags().IntVar(&opts.Status, "status", 1, "status")
	return cmd
}

func newOperatorDeleteCmd() *cobra.Command {
	return deleteCmd(resource_type_operator, Sdk.DeleteOperator)
}

func newOperatorListCmd() *cobra.Command {
	defCs := []string{"OperatorID", "CompanyID", "Name", "Password", "Job", "Department", "Status", "Role"}
	cmd := &cobra.Command{
		Use:   "operator-list companyId",
		Short: "List operators Support Filter",
	}
	columns := cmd.Flags().StringArrayP("column", "c", defCs, "Columns to display")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		var opts options.ListOptions
		opts.Limit = -1
		opts.Offset = 0
		opts.Filter = make(map[string]interface{})
		if id, err := strconv.Atoi(args[0]); err != nil {
			fmt.Println("change to int fail", args[0])
		} else {
			opts.Filter["company_id"] = id
		}
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
		Use:   "operator-show [OPTIONS] operatorID",
		Short: "Show a operator information",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			var opts options.BaseOptions
			if id, err := strconv.Atoi(args[0]); err != nil {
				fmt.Println("change to int fail", args[0])
			} else {
				opts.ID = id
			}
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
		Use:   "operator-update [OPTIONS] operatorID job",
		Short: "update a operator",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			//opts.Name = args[0]
			if id, err := strconv.Atoi(args[0]); err != nil {
				fmt.Println("change to int fail", args[0])
			} else {
				opts.OperatorID = id
			}
			//for test
			opts.Job = args[1]
			if err := Sdk.UpdateOperator(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	return cmd
}
