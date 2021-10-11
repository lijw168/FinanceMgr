package command

import (
	"analysis-server/cli/util"
	"analysis-server/sdk/options"
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

func NewVoucherCommand(cmd *cobra.Command) {
	cmd.AddCommand(newVoucherCreateCmd())
	cmd.AddCommand(newVoucherDeleteCmd())
	cmd.AddCommand(newVoucherShowCmd())
	cmd.AddCommand(newVoucherArrangeCmd())

	cmd.AddCommand(newVoucherRecordCreateCmd())
	cmd.AddCommand(newVoucherRecordDeleteCmd())
	cmd.AddCommand(newVoucherRecordUpdateCmd())
	cmd.AddCommand(newVoucherRecordListCmd())

	cmd.AddCommand(newVoucherInfoShowCmd())
	cmd.AddCommand(newVoucherInfoListCmd())
	cmd.AddCommand(newGetLatestVouInfoCmd())
	cmd.AddCommand(newGetMaxNumOfMonthCmd())
	cmd.AddCommand(newVoucherInfoUpdateCmd())
}

func newVoucherCreateCmd() *cobra.Command {
	var opts options.VoucherOptions
	var createRecOpt options.CreateVoucherRecordOptions
	cmd := &cobra.Command{
		Use:   "voucher-create [OPTIONS] companyID voucherMonth voucherFiller",
		Short: "Create a voucher",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}
			if id, err := strconv.Atoi(args[0]); err != nil {
				fmt.Println("change to int fail", args[0])
			} else {
				opts.InfoOptions.CompanyID = id
			}

			if month, err := strconv.Atoi(args[1]); err != nil {
				fmt.Println("change to int fail", args[1])
			} else {
				opts.InfoOptions.VoucherMonth = month
			}
			opts.InfoOptions.VoucherFiller = args[2]
			opts.RecordsOptions = append(opts.RecordsOptions, createRecOpt)
			if hv, err := Sdk.CreateVoucher(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(hv)
			}
		},
	}
	cmd.Flags().StringVar(&createRecOpt.SubjectName, "subject", "test", "subject name")
	cmd.Flags().StringVar(&createRecOpt.Summary, "summary", "test", "summary")
	//如下的两个字段，只能用于测试。因为cmd的中，没有实现float64的操作。
	var dm, cm int
	cmd.Flags().IntVar(&dm, "dm", 1, "debit money")
	cmd.Flags().IntVar(&cm, "cm", 1, "credit money")
	createRecOpt.DebitMoney = float64(dm)
	createRecOpt.CreditMoney = float64(cm)
	cmd.Flags().IntVar(&createRecOpt.SubID1, "sub1", 1, "SubID1")
	cmd.Flags().IntVar(&createRecOpt.SubID2, "sub2", 2, "SubID2")
	return cmd
}

func newVoucherDeleteCmd() *cobra.Command {
	return deleteCmd(resource_type_voucher, Sdk.DeleteVoucher)
}

func newVoucherShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "voucher-show [OPTIONS] voucherId",
		Short: "Show voucher",
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
			view, err := Sdk.GetVoucher(&opts)
			if err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(view)
			}
		},
	}
	return cmd
}

func newVoucherArrangeCmd() *cobra.Command {
	var opts options.VoucherArrangeOptions
	cmd := &cobra.Command{
		Use:   "voucher-arrange [OPTIONS] companyID voucherMonth",
		Short: "voucher arrange",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}
			if id, err := strconv.Atoi(args[0]); err != nil {
				fmt.Println("change to int fail", args[0])
			} else {
				opts.CompanyID = id
			}
			if voucherMonth, err := strconv.Atoi(args[1]); err != nil {
				fmt.Println("change to int fail", args[1])
			} else {
				opts.VoucherMonth = voucherMonth
			}
			if err := Sdk.ArrangeVoucher(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	cmd.Flags().BoolVar(&opts.ArrangeVoucherNum, "isArrangeVoucherNum", false, "arrange voucher Num")
	return cmd
}

func newVoucherRecordCreateCmd() *cobra.Command {
	var opts options.CreateVoucherRecordOptions
	cmd := &cobra.Command{
		Use:   "vouRecord-create [OPTIONS] voucherId subject-name ...",
		Short: "Create a voucher record",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			if id, err := strconv.Atoi(args[0]); err != nil {
				fmt.Println("change to int fail", args[0])
			} else {
				opts.VoucherID = id
			}
			opts.SubjectName = args[1]
			optSlice := []options.CreateVoucherRecordOptions{}
			optSlice = append(optSlice, opts)

			if hv, err := Sdk.CreateVoucherRecords(optSlice); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(hv)
			}
		},
	}
	//cmd.Flags().StringVar(&opts.SubjectName, "subject name", "test", "subject name")
	cmd.Flags().StringVar(&opts.Summary, "summary", "test", "summary")
	var dm, cm int
	cmd.Flags().IntVar(&dm, "dm", 1, "debit money")
	cmd.Flags().IntVar(&cm, "cm", 1, "credit money")
	opts.DebitMoney = float64(dm)
	opts.CreditMoney = float64(cm)
	cmd.Flags().IntVar(&opts.SubID1, "sub1", 1, "SubID1")
	cmd.Flags().IntVar(&opts.SubID2, "sub2", 2, "SubID2")
	return cmd
}

func newVoucherRecordDeleteCmd() *cobra.Command {
	return deleteCmd(resource_type_voucher_record, Sdk.DeleteVoucherRecord)
}

func newVoucherRecordListCmd() *cobra.Command {
	defCs := []string{"RecordID", "VoucherID", "SubjectName", "DebitMoney", "CreditMoney", "Summary",
		"SubID1", "SubID2", "SubID3", "SubID4", "BillCount", "Status"}
	cmd := &cobra.Command{
		Use:   "vouRecord-list voucherId",
		Short: "List voucher records Support Filter",
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
		opts.Filter = make(map[string]interface{})
		if id, err := strconv.Atoi(args[0]); err != nil {
			fmt.Println("change to int fail", args[0])
		} else {
			opts.Filter["voucherId"] = id
		}
		if _, views, err := Sdk.ListVoucherRecords(&opts); err != nil {
			util.FormatErrorOutput(err)
		} else {
			util.FormatListOutput(*columns, views)
		}
	}
	return cmd
}

func newVoucherRecordUpdateCmd() *cobra.Command {
	var opts options.ModifyVoucherRecordOptions
	cmd := &cobra.Command{
		Use:   "vouRecord-update [OPTIONS] vouRecordId summary status",
		Short: "update a voucher record",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}
			if id, err := strconv.Atoi(args[0]); err != nil {
				fmt.Println("change to int fail", args[0])
			} else {
				opts.VouRecordID = id
			}
			if status, err := strconv.Atoi(args[1]); err != nil {
				fmt.Println("change to int fail", args[1])
			} else {
				opts.Status = status
			}
			opts.Summary = args[2]
			opts.BillCount = -1
			opts.CreditMoney = -1
			opts.DebitMoney = -1
			if err := Sdk.UpdateVoucherRecordByID(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	cmd.Flags().StringVar(&opts.SubjectName, "subName", "", "subjectName")
	cmd.Flags().StringVar(&opts.Summary, "summary", "", "summary")
	//因为金额，可以允许为负数，所以不能用这样的默认参数。所以不能用cli修改金额。可以用客户端来修改。
	// var dm, cm int
	// cmd.Flags().IntVar(&dm, "dm", 0, "debit money")
	// cmd.Flags().IntVar(&cm, "cm", 0, "credit money")
	// opts.DebitMoney = float64(dm)
	// opts.CreditMoney = float64(cm)
	cmd.Flags().IntVar(&opts.SubID1, "sub1", 0, "SubID1")
	//cmd.Flags().IntVar(&opts.SubID2, "sub2", 0, "SubID2")
	return cmd
}

func newVoucherInfoShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vouInfo-show [OPTIONS] voucherID",
		Short: "Show a voucher information",
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
			view, err := Sdk.GetVoucherInfo(&opts)
			if err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(view)
			}
		},
	}
	return cmd
}

func newGetLatestVouInfoCmd() *cobra.Command {
	defCs := []string{"VoucherID", "CompanyID", "VoucherMonth", "NumOfMonth", "VoucherDate",
		"VoucherFiller", "VoucherAuditor"}
	cmd := &cobra.Command{
		Use:   "vouInfo-getLatest [OPTIONS] companyID",
		Short: "get latest voucher information",
	}
	columns := cmd.Flags().StringArrayP("column", "c", defCs, "Columns to display")
	cmd.Run = func(cmd *cobra.Command, args []string) {
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
		if _, views, err := Sdk.GetLatestVoucherInfo(&opts); err != nil {
			util.FormatErrorOutput(err)
		} else {
			util.FormatListOutput(*columns, views)
		}
	}

	return cmd
}

func newGetMaxNumOfMonthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vouInfo-getMaxNumOfMan [OPTIONS] companyID voucherMonth",
		Short: "get the max numOfMonth in a month",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}
			var opts options.QueryMaxNumOfMonthOption
			if id, err := strconv.Atoi(args[0]); err != nil {
				fmt.Println("change to int fail", args[0])
			} else {
				opts.CompanyID = id
			}
			if month, err := strconv.Atoi(args[1]); err != nil {
				fmt.Println("change to int fail", args[1])
			} else {
				opts.VoucherMonth = month
			}
			iCount, err := Sdk.GetMaxNumOfMonth(&opts)
			if err != nil {
				util.FormatErrorOutput(err)
			} else {
				//util.FormatViewOutput(iCount)
				fmt.Println("the max numOfMonth is:", iCount)
			}
		},
	}
	return cmd
}

func newVoucherInfoListCmd() *cobra.Command {
	defCs := []string{"VoucherID", "CompanyID", "VoucherMonth", "NumOfMonth", "VoucherDate", "VoucherFiller", "VoucherAuditor"}
	cmd := &cobra.Command{
		Use:   "vouInfo-list companyId",
		Short: "List voucherInfo Support Filter",
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
		opts.Filter = make(map[string]interface{})
		if id, err := strconv.Atoi(args[0]); err != nil {
			fmt.Println("change to int fail", args[0])
		} else {
			opts.Filter["companyId"] = id
		}
		if _, views, err := Sdk.ListVoucherInfo(&opts); err != nil {
			util.FormatErrorOutput(err)
		} else {
			util.FormatListOutput(*columns, views)
		}
	}
	return cmd
}

func newVoucherInfoUpdateCmd() *cobra.Command {
	var opts options.ModifyVoucherInfoOptions
	cmd := &cobra.Command{
		Use:   "vouInfo-update [OPTIONS] voucherId",
		Short: "update a voucher information",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			if id, err := strconv.Atoi(args[0]); err != nil {
				fmt.Println("change to int fail", args[0])
			} else {
				opts.VoucherID = id
			}
			if err := Sdk.UpdateVoucherInfo(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	cmd.Flags().StringVar(&opts.VoucherFiller, "vouFiller", "", "voucher filler")
	cmd.Flags().StringVar(&opts.VoucherAuditor, "vouAuditor", "", "voucher auditor")
	cmd.Flags().IntVar(&opts.VoucherDate, "vouDate", 0, "voucher date")
	cmd.Flags().IntVar(&opts.BillCount, "billCount", 0, "voucher bill count")
	cmd.Flags().IntVar(&opts.Status, "status", 0, "voucher status")
	return cmd
}
