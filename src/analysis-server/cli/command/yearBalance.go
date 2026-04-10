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
		Use:   "yearBal-create [OPTIONS] companyID year subjectID ",
		Short: "Create the record of year balance",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 3 {
				cmd.Help()
				return
			}
			errMsgs := []string{}
			companyID, err := strconv.Atoi(args[0])
			if err != nil {
				errMsgs = append(errMsgs, fmt.Sprintf("companyID转换失败: %s", args[0]))
			}
			year, err := strconv.Atoi(args[1])
			if err != nil {
				errMsgs = append(errMsgs, fmt.Sprintf("year转换失败: %s", args[1]))
			}
			subjectID, err := strconv.Atoi(args[2])
			if err != nil {
				errMsgs = append(errMsgs, fmt.Sprintf("subjectID转换失败: %s", args[2]))
			}
			if len(errMsgs) > 0 {
				for _, msg := range errMsgs {
					fmt.Println(msg)
				}
				return
			}
			opts.CompanyID = companyID
			opts.Year = year
			opts.SubjectID = subjectID
			if err := Sdk.CreateYearBalance(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	cmd.Flags().Float64Var(&opts.Balance, "balance", 0, "year balance")
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
			errMsgs := []string{}
			companyID, err := strconv.Atoi(args[0])
			if err != nil {
				errMsgs = append(errMsgs, fmt.Sprintf("companyID转换失败: %s", args[0]))
			}
			year, err := strconv.Atoi(args[1])
			if err != nil {
				errMsgs = append(errMsgs, fmt.Sprintf("year转换失败: %s", args[1]))
			}
			subjectID, err := strconv.Atoi(args[2])
			if err != nil {
				errMsgs = append(errMsgs, fmt.Sprintf("subjectID转换失败: %s", args[2]))
			}
			if len(errMsgs) > 0 {
				for _, msg := range errMsgs {
					fmt.Println(msg)
				}
				return
			}
			opts.CompanyID = companyID
			opts.Year = year
			opts.SubjectID = subjectID
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
			errMsgs := []string{}
			var opts options.BasicYearBalance
			companyID, err := strconv.Atoi(args[0])
			if err != nil {
				errMsgs = append(errMsgs, fmt.Sprintf("companyID转换失败: %s", args[0]))
			}
			year, err := strconv.Atoi(args[1])
			if err != nil {
				errMsgs = append(errMsgs, fmt.Sprintf("year转换失败: %s", args[1]))
			}
			subjectID, err := strconv.Atoi(args[2])
			if err != nil {
				errMsgs = append(errMsgs, fmt.Sprintf("subjectID转换失败: %s", args[2]))
			}
			if len(errMsgs) > 0 {
				for _, msg := range errMsgs {
					fmt.Println(msg)
				}
				return
			}
			opts.CompanyID = companyID
			opts.Year = year
			opts.SubjectID = subjectID
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
		Use:   "yearBal-update [OPTIONS] companyID  year subjectID",
		Short: "update a record  of year balance",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 3 {
				cmd.Help()
				return
			}
			errMsgs := []string{}
			companyID, err := strconv.Atoi(args[0])
			if err != nil {
				errMsgs = append(errMsgs, fmt.Sprintf("companyID转换失败: %s", args[0]))
			}
			year, err := strconv.Atoi(args[1])
			if err != nil {
				errMsgs = append(errMsgs, fmt.Sprintf("year转换失败: %s", args[1]))
			}
			subjectID, err := strconv.Atoi(args[2])
			if err != nil {
				errMsgs = append(errMsgs, fmt.Sprintf("subjectID转换失败: %s", args[2]))
			}
			if len(errMsgs) > 0 {
				for _, msg := range errMsgs {
					fmt.Println(msg)
				}
				return
			}
			opts.CompanyID = companyID
			opts.Year = year
			opts.SubjectID = subjectID
			if err := Sdk.UpdateYearBalance(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	cmd.Flags().Float64Var(&opts.Balance, "balance", 0, "year balance")
	cmd.Flags().IntVar(&opts.Status, "status", 0, "annual closing status")
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
		errMsgs := []string{}
		var opts options.ListOptions
		opts.Limit = -1
		opts.Offset = 0
		//for test
		opts.Filter = make(map[string]interface{})
		companyID, err := strconv.Atoi(args[0])
		if err != nil {
			errMsgs = append(errMsgs, fmt.Sprintf("companyID转换失败: %s", args[0]))
		}
		year, err := strconv.Atoi(args[1])
		if err != nil {
			errMsgs = append(errMsgs, fmt.Sprintf("year转换失败: %s", args[1]))
		}
		subjectID, err := strconv.Atoi(args[2])
		if err != nil {
			errMsgs = append(errMsgs, fmt.Sprintf("subjectID转换失败: %s", args[2]))
		}
		if len(errMsgs) > 0 {
			for _, msg := range errMsgs {
				fmt.Println(msg)
			}
			return
		}
		opts.Filter["companyId"] = companyID
		opts.Filter["year"] = year
		opts.Filter["subjectId"] = subjectID

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
		Short: "Show the balance value of year balance record ",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 3 {
				cmd.Help()
				return
			}
			errMsgs := []string{}
			var opts options.BasicYearBalance
			companyID, err := strconv.Atoi(args[0])
			if err != nil {
				errMsgs = append(errMsgs, fmt.Sprintf("companyID转换失败: %s", args[0]))
			}
			year, err := strconv.Atoi(args[1])
			if err != nil {
				errMsgs = append(errMsgs, fmt.Sprintf("year转换失败: %s", args[1]))
			}
			subjectID, err := strconv.Atoi(args[2])
			if err != nil {
				errMsgs = append(errMsgs, fmt.Sprintf("subjectID转换失败: %s", args[2]))
			}
			if len(errMsgs) > 0 {
				for _, msg := range errMsgs {
					fmt.Println(msg)
				}
				return
			}
			opts.CompanyID = companyID
			opts.Year = year
			opts.SubjectID = subjectID
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
