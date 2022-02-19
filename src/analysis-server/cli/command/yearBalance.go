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
	cmd.AddCommand(newYearBalanceUpdateCmd())
	cmd.AddCommand(newYearBalanceListCmd())
}

func newYearBalanceCreateCmd() *cobra.Command {
	var opts options.YearBalanceOption
	cmd := &cobra.Command{
		Use:   "yearBal-create [OPTIONS] subjectId year balance",
		Short: "Create the begin of year balance",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 3 {
				cmd.Help()
				return
			}
			if subjectId, err := strconv.Atoi(args[0]); err != nil {
				fmt.Println("change to int fail", args[0])
			} else {
				opts.SubjectID = subjectId
			}
			if iYear, err := strconv.Atoi(args[1]); err != nil {
				fmt.Println("change to int fail", args[1])
			} else {
				opts.Year = iYear
			}
			if balance, err := strconv.ParseFloat(args[2], 32); err != nil {
				fmt.Println("change to int fail", args[2])
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
	//return deleteCmd(resource_type_year_balance, Sdk.DeleteYearBalance)
	var opts options.BasicYearBalance
	cmd := &cobra.Command{
		Use:   "yearBal-delete [OPTIONS] subjectId year",
		Short: "delete the begin of year balance",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}
			if subjectId, err := strconv.Atoi(args[0]); err != nil {
				fmt.Println("change to int fail", args[0])
			} else {
				opts.SubjectID = subjectId
			}
			if iYear, err := strconv.Atoi(args[1]); err != nil {
				fmt.Println("change to int fail", args[1])
			} else {
				opts.Year = iYear
			}
			if err := Sdk.DeleteYearBalance(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	return cmd
}

func newYearBalanceShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "yearBal-show [OPTIONS] subjectId year",
		Short: "Show the begin of year balance ",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}
			var opts options.BasicYearBalance
			if subjectId, err := strconv.Atoi(args[0]); err != nil {
				fmt.Println("change to int fail", args[0])
			} else {
				opts.SubjectID = subjectId
			}
			if iYear, err := strconv.Atoi(args[1]); err != nil {
				fmt.Println("change to int fail", args[1])
			} else {
				opts.Year = iYear
			}
			dYearBal, err := Sdk.GetYearBalance(&opts)
			if err != nil {
				util.FormatErrorOutput(err)
			} else {
				//util.FormatViewOutput(dYearBal)
				fmt.Printf("the dYearBal is %f\r\n", dYearBal)
			}
		},
	}
	return cmd
}

func newYearBalanceUpdateCmd() *cobra.Command {
	var opts options.YearBalanceOption
	cmd := &cobra.Command{
		Use:   "yearBal-update [OPTIONS] subjectID year balance",
		Short: "update a year balance of a accountSubject",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 5 {
				cmd.Help()
				return
			}
			if id, err := strconv.Atoi(args[0]); err != nil {
				fmt.Println("change to int fail", args[0])
			} else {
				opts.SubjectID = id
			}
			if iYear, err := strconv.Atoi(args[1]); err != nil {
				fmt.Println("change to int fail", args[1])
			} else {
				opts.Year = iYear
			}
			if dYearBal, err := strconv.ParseFloat(args[2], 64); err != nil {
				fmt.Println("change to int fail", args[2])
			} else {
				opts.Balance = dYearBal
			}

			if err := Sdk.UpdateYearBalance(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	return cmd
}
func newYearBalanceListCmd() *cobra.Command {
	defCs := []string{"SubjectID", "year", "yearBalance"}
	cmd := &cobra.Command{
		Use:   "yearBal-list year",
		Short: "List account subjects Support Filter",
	}
	columns := cmd.Flags().StringArrayP("column", "c", defCs, "Columns to display")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Help()
			return
		}
		var opts options.ListOptions
		opts.Limit = -1
		opts.Offset = 0
		//for test
		opts.Filter = make(map[string]interface{})
		if iYear, err := strconv.Atoi(args[0]); err != nil {
			fmt.Println("change to int fail", args[0])
		} else {
			opts.Filter["year"] = iYear
		}

		if _, yearBalViews, err := Sdk.ListYearBalance(&opts); err != nil {
			util.FormatErrorOutput(err)
		} else {
			util.FormatListOutput(*columns, yearBalViews)
		}
	}
	return cmd
}
