package command

import (
	"financeMgr/src/analysis-server/cli/util"
	"financeMgr/src/analysis-server/sdk/options"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func NewCompanyCommand(cmd *cobra.Command) {
	cmd.AddCommand(newCompanyCreateCmd())
	cmd.AddCommand(newCompanyDeleteCmd())
	cmd.AddCommand(newCompanyListCmd())
	cmd.AddCommand(newCompanyShowCmd())
	cmd.AddCommand(newCompanyUpdateCmd())
	cmd.AddCommand(newAssociatedCompanyGroupCmd())
	cmd.AddCommand(newInitResourceInfoCmd())
}

func newCompanyCreateCmd() *cobra.Command {
	var opts options.CreateCompanyOptions
	cmd := &cobra.Command{
		Use:   "company-create [OPTIONS] companyName abbreName,corporator,phone,e_mail,companyAdd,startAccountPeriod",
		Short: "Create a company",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 7 {
				cmd.Help()
				return
			}
			opts.CompanyName = args[0]
			opts.AbbrevName = args[1]
			opts.Corporator = args[2]
			opts.Phone = args[3]
			opts.Email = args[4]
			opts.CompanyAddr = args[5]
			if period, err := strconv.Atoi(args[6]); err != nil {
				fmt.Println("change to int fail", args[0])
			} else {
				opts.StartAccountPeriod = period
			}

			if hv, err := Sdk.CreateCompany(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(hv)
			}
		},
	}
	return cmd
}

func newCompanyDeleteCmd() *cobra.Command {
	return deleteCmd(resource_type_company, Sdk.DeleteCompany)
}

func newCompanyListCmd() *cobra.Command {
	defCs := []string{"CompanyID", "CompanyName", "AbbrevName", "Corporator", "Phone",
		"Email", "CompanyAddr", "Backup", "CreatedAt", "UpdatedAt", "CompanyGroupID",
		"StartAccountPeriod", "LatestAccountYear"}
	cmd := &cobra.Command{
		Use:   "company-list ",
		Short: "List company Support Filter",
	}
	columns := cmd.Flags().StringArrayP("column", "c", defCs, "Columns to display")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		var opts options.ListOptions
		opts.Limit = -1
		opts.Offset = 0
		//for test
		// opts.Filter = make(map[string]interface{})
		// opts.Filter["backup"] = "test"
		if _, views, err := Sdk.ListCompany(&opts); err != nil {
			util.FormatErrorOutput(err)
		} else {
			util.FormatListOutput(*columns, views)
		}
	}
	return cmd
}

func newCompanyShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "company-show [OPTIONS] companyId",
		Short: "Show company",
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
			view, err := Sdk.GetCompany(&opts)
			if err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(view)
			}
		},
	}
	return cmd
}

func newCompanyUpdateCmd() *cobra.Command {
	var opts options.ModifyCompanyOptions
	cmd := &cobra.Command{
		Use:   "company-update [OPTIONS] companyId CompanyName",
		Short: "update a accSub",
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
			opts.CompanyName = args[1]

			if err := Sdk.UpdateCompany(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	cmd.Flags().StringVar(&opts.CompanyName, "comName", "", "CompanyName")
	cmd.Flags().StringVar(&opts.AbbrevName, "abbrevName", "", "AbbrevName")
	cmd.Flags().StringVar(&opts.CompanyAddr, "comAddr", "", "CompanyAddr")
	cmd.Flags().StringVar(&opts.Corporator, "corporator", "", "Corporator")
	cmd.Flags().StringVar(&opts.Email, "email", "", "Email")
	cmd.Flags().StringVar(&opts.Backup, "bc", "", "backup")
	cmd.Flags().IntVar(&opts.LatestAccountYear, "accYear", 0, "LatestAccountYear")
	return cmd
}

func newAssociatedCompanyGroupCmd() *cobra.Command {
	var opts options.AssociatedCompanyGroupOptions
	cmd := &cobra.Command{
		Use:   "companyGroup-associated [OPTIONS] companyGroupId companyId",
		Short: "associate a company group",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}
			if id, err := strconv.Atoi(args[0]); err != nil {
				fmt.Println("change to int fail", args[0])
			} else {
				opts.CompanyGroupID = id
			}

			if id, err := strconv.Atoi(args[1]); err != nil {
				fmt.Println("change to int fail", args[1])
			} else {
				opts.CompanyID = id
			}

			if err := Sdk.AssociatedCompanyGroup(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	cmd.Flags().BoolVar(&opts.IsAttach, "isAttach", true, "is attach")
	return cmd
}

func newInitResourceInfoCmd() *cobra.Command {
	var opts options.BaseOptions
	cmd := &cobra.Command{
		Use:   "resourceInfo-init [OPTIONS] operatorId",
		Short: "init resource information",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			if id, err := strconv.Atoi(args[0]); err != nil {
				fmt.Println("change to int fail", args[0])
			} else {
				opts.ID = id
			}
			if hv, err := Sdk.InitResourceInfo(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(hv)
			}
		},
	}
	return cmd
}
