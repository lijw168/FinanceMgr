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

	cmd.AddCommand(newVoucherRecordCreateCmd())
	cmd.AddCommand(newVoucherRecordDeleteCmd())
	cmd.AddCommand(newVoucherRecordUpdateCmd())
	cmd.AddCommand(newVoucherRecordListCmd())

	cmd.AddCommand(newVoucherInfoShowCmd())
	cmd.AddCommand(newVoucherInfoListCmd())
}

func newVoucherCreateCmd() *cobra.Command {
	var opts options.VoucherOptions
	var createRecOpt options.CreateVoucherRecordOptions
	//params := model.VoucherParams{}
	//vouInfoParam := model.VoucherInfoParams{CompanyID: &(opts.InfoOptions.CompanyID), VoucherMonth: &(opts.InfoOptions.VoucherMonth)}
	cmd := &cobra.Command{
		Use:   "voucher-create [OPTIONS] companyID voucherMonth",
		Short: "Create a voucher",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
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
			opts.RecordsOptions = append(opts.RecordsOptions, createRecOpt)
			if hv, err := Sdk.CreateVoucher(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(hv)
			}
		},
	}
	cmd.Flags().StringVar(&createRecOpt.SubjectName, "subject name", "test", "subject name")
	cmd.Flags().StringVar(&createRecOpt.Summary, "summary", "test", "summary")
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
	return deleteCmd(resource_type_voucher_info, Sdk.DeleteVoucher)
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
		"SubID1", "SubID2", "SubID3", "SubID4", "BillCount"}
	cmd := &cobra.Command{
		Use:   "vouRecord-list ",
		Short: "List voucher records Support Filter",
	}
	columns := cmd.Flags().StringArrayP("column", "c", defCs, "Columns to display")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		var opts options.ListOptions
		opts.Limit = -1
		opts.Offset = 0
		//for test
		//opts.Filter = make(map[string]interface{})
		//opts.Filter["status"] = "creating|available|in-use"
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
		Use:   "vouRecord-update [OPTIONS] vouRecordId summary",
		Short: "update a voucher record",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			if id, err := strconv.Atoi(args[0]); err != nil {
				fmt.Println("change to int fail", args[0])
			} else {
				opts.VouRecordID = id
			}
			opts.Summary = args[1]
			if err := Sdk.UpdateVoucherRecord(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	// cmd.Flags().StringVar(&opts.SubjectName, "subName", "test_update", "subjectName")
	// cmd.Flags().StringVar(&opts.Summary, "summary", "test_update", "summary")
	// var dm, cm int
	// cmd.Flags().IntVar(&dm, "dm", 1, "debit money")
	// cmd.Flags().IntVar(&cm, "cm", 1, "credit money")
	// opts.DebitMoney = float64(dm)
	// opts.CreditMoney = float64(cm)
	// cmd.Flags().IntVar(&opts.SubID1, "sub1", 1, "SubID1")
	// cmd.Flags().IntVar(&opts.SubID2, "sub2", 2, "SubID2")
	return cmd
}

func newVoucherInfoShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vouInfo-show [OPTIONS] voucherID",
		Short: "Show a voucher",
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

func newVoucherInfoListCmd() *cobra.Command {
	defCs := []string{"VoucherID", "CompanyID", "VoucherMonth", "NumOfMonth", "VoucherDate"}
	cmd := &cobra.Command{
		Use:   "vouInfo-list ",
		Short: "List voucherInfo Support Filter",
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
