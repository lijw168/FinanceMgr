package command

import (
	"analysis-server/cli/util"
	"analysis-server/sdk/options"
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

func NewYearBalanceCommand(cmd *cobra.Command) {
	cmd.AddCommand(newYearBalanceCreateCmd())
	cmd.AddCommand(newYearBalanceDeleteCmd())
	cmd.AddCommand(newYearBalanceShowCmd())
}

func newYearBalanceCreateCmd() *cobra.Command {
	var opts options.YearBalanceOption
	cmd := &cobra.Command{
		Use:   "yearBal-create [OPTIONS] subjectId summary subjectDirection balance",
		Short: "Create the begin of year balance",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 4 {
				cmd.Help()
				return
			}
			if subjectId, err := strconv.Atoi(args[0]); err != nil {
				fmt.Println("change to int fail", args[0])
			} else {
				opts.SubjectID = subjectId
			}
			opts.Summary = args[1]
			if subDir, err := strconv.Atoi(args[2]); err != nil {
				fmt.Println("change to int fail", args[2])
			} else {
				opts.SubjectDirection = subDir
			}
			if balance, err := strconv.ParseFloat(args[3], 32); err != nil {
				fmt.Println("change to int fail", args[3])
			} else {
				opts.Balance = balance
			}
			if err := Sdk.CreateYearBalance(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	return cmd
}

func newYearBalanceDeleteCmd() *cobra.Command {
	return deleteCmd(resource_type_year_balance, Sdk.DeleteYearBalanceByID)
}

func newYearBalanceShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "yearBal-show [OPTIONS] subjectId",
		Short: "Show the begin of year balance by subjectId ",
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
			view, err := Sdk.GetYearBalanceById(&opts)
			if err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(view)
			}
		},
	}
	return cmd
}
