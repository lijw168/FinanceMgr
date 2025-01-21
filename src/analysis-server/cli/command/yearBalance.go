package command

import (
	"financeMgr/src/analysis-server/cli/util"
	"financeMgr/src/analysis-server/sdk/options"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func NewYearBalanceCommand(cmd *cobra.Command) {
	cmd.AddCommand(newYearBalanceCreateCmd())
	cmd.AddCommand(newYearBalanceDeleteCmd())
	cmd.AddCommand(newYearBalanceShowCmd())
	cmd.AddCommand(newAccSubYearBalValueShowCmd())
	cmd.AddCommand(newYearBalanceUpdateCmd())
	cmd.AddCommand(newYearBalanceListCmd())
}

func newYearBalanceCreateCmd() *cobra.Command {
	var opts options.YearBalanceOption
	cmd := &cobra.Command{
		Use:   "yearBal-create [OPTIONS] companyID year subjectID [FLAG] balance",
		Short: "Create the record of year balance",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 3 {
				cmd.Help()
				return
			}
			if id, err := strconv.Atoi(args[0]); err != nil {
				fmt.Println("change to int fail", args[0])
			} else {
				opts.CompanyID = id
			}
			if iYear, err := strconv.Atoi(args[1]); err != nil {
				fmt.Println("change to int fail", args[1])
			} else {
				opts.Year = iYear
			}
			if id, err := strconv.Atoi(args[2]); err != nil {
				fmt.Println("change to int fail", args[2])
			} else {
				opts.SubjectID = id
			}
			if err := Sdk.CreateYearBalance(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	cmd.Flags().Float64Var(&opts.Balance, "balance", 0, "annual closing status")
	//testBal := *(cmd.Flags().Float64("balance", 0, "year balance"))
	//cmd.Flags().IntVar(&opts.Status, "status", 0, "annual closing status")
	//fmt.Printf("testBal:%f,balance:%f\r\n", testBal, opts.Balance)
	return cmd
}

func newYearBalanceDeleteCmd() *cobra.Command {
	//return deleteCmd(resource_type_year_balance, Sdk.DeleteYearBalance)
	var opts options.BasicYearBalance
	cmd := &cobra.Command{
		Use:   "yearBal-delete [OPTIONS] companyId year subjectId",
		Short: "delete the record of year balance",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 3 {
				cmd.Help()
				return
			}
			if id, err := strconv.Atoi(args[0]); err != nil {
				fmt.Println("change to int fail", args[0])
			} else {
				opts.CompanyID = id
			}

			if iYear, err := strconv.Atoi(args[1]); err != nil {
				fmt.Println("change to int fail", args[1])
			} else {
				opts.Year = iYear
			}

			if subjectId, err := strconv.Atoi(args[2]); err != nil {
				fmt.Println("change to int fail", args[2])
			} else {
				opts.SubjectID = subjectId
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
		Use:   "yearBal-show [OPTIONS] companyId year subjectId",
		Short: "Show the record of year balance ",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 3 {
				cmd.Help()
				return
			}
			var opts options.BasicYearBalance
			if companyId, err := strconv.Atoi(args[0]); err != nil {
				fmt.Println("change companyId to int fail", args[0])
			} else {
				opts.CompanyID = companyId
			}
			if iYear, err := strconv.Atoi(args[1]); err != nil {
				fmt.Println("change iYear to int fail", args[1])
			} else {
				opts.Year = iYear
			}
			if subjectId, err := strconv.Atoi(args[2]); err != nil {
				fmt.Println("change subjectId to int fail", args[2])
			} else {
				opts.SubjectID = subjectId
			}
			yearBalView, err := Sdk.GetYearBalance(&opts)
			if err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(yearBalView)
			}
		},
	}
	return cmd
}

func newYearBalanceUpdateCmd() *cobra.Command {
	var opts options.YearBalanceOption
	cmd := &cobra.Command{
		Use:   "yearBal-update [OPTIONS] companyID  year subjectID [flag] balance/status",
		Short: "update a record  of year balance",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 3 {
				cmd.Help()
				return
			}
			if id, err := strconv.Atoi(args[0]); err != nil {
				fmt.Println("change to int fail", args[0])
			} else {
				opts.CompanyID = id
			}
			if iYear, err := strconv.Atoi(args[1]); err != nil {
				fmt.Println("change to int fail", args[1])
			} else {
				opts.Year = iYear
			}
			if id, err := strconv.Atoi(args[2]); err != nil {
				fmt.Println("change to int fail", args[2])
			} else {
				opts.SubjectID = id
			}
			if err := Sdk.UpdateYearBalance(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	//不清楚为何Status,是可以获取数据的，难到时因为是其第一个字符大写？下一个版本验证一下。
	opts.Balance = *(cmd.Flags().Float64("balance", 0, "year balance"))
	cmd.Flags().IntVar(&opts.Status, "Status", 0, "annual closing status")
	return cmd
}
func newYearBalanceListCmd() *cobra.Command {
	//因为只返回了这四列的数据
	defCs := []string{"year", "SubjectID", "Balance", "Status"}
	cmd := &cobra.Command{
		Use:   "yearBal-list [OPTIONS] companID year subjectID",
		Short: "List account subjects year balance Support Filter",
	}
	columns := cmd.Flags().StringArrayP("column", "c", defCs, "Columns to display")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			cmd.Help()
			return
		}
		var opts options.ListOptions
		opts.Limit = -1
		opts.Offset = 0
		//for test
		opts.Filter = make(map[string]interface{})
		if companyId, err := strconv.Atoi(args[0]); err != nil {
			fmt.Println("change to int fail", args[0])
		} else {
			opts.Filter["companyId"] = companyId
		}
		if iYear, err := strconv.Atoi(args[1]); err != nil {
			fmt.Println("change to int fail", args[1])
		} else {
			opts.Filter["year"] = iYear
		}
		if id, err := strconv.Atoi(args[2]); err != nil {
			fmt.Println("change to int fail", args[2])
		} else {
			opts.Filter["subjectId"] = id
		}

		if _, yearBalViews, err := Sdk.ListYearBalance(&opts); err != nil {
			util.FormatErrorOutput(err)
		} else {
			util.FormatListOutput(*columns, yearBalViews)
		}
	}
	return cmd
}

// 仅获取但个balance
func newAccSubYearBalValueShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "yearBal-accSubBal-show [OPTIONS] companyId year subjectId",
		Short: "Show the balance value of year balance ",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 3 {
				cmd.Help()
				return
			}
			var opts options.BasicYearBalance
			if companyId, err := strconv.Atoi(args[0]); err != nil {
				fmt.Println("change companyId to int fail", args[0])
			} else {
				opts.CompanyID = companyId
			}
			if iYear, err := strconv.Atoi(args[1]); err != nil {
				fmt.Println("change iYear to int fail", args[1])
			} else {
				opts.Year = iYear
			}
			if subjectId, err := strconv.Atoi(args[2]); err != nil {
				fmt.Println("change subjectId to int fail", args[2])
			} else {
				opts.SubjectID = subjectId
			}
			dYearBal, err := Sdk.GetAccSubYearBalValue(&opts)
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
