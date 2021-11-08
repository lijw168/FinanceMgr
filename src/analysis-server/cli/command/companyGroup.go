package command

import (
	"analysis-server/cli/util"
	"analysis-server/sdk/options"
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

func NewCompanyGroupCommand(cmd *cobra.Command) {
	cmd.AddCommand(newCompanyGroupCreateCmd())
	cmd.AddCommand(newCompanyGroupDeleteCmd())
	cmd.AddCommand(newCompanyGroupListCmd())
	cmd.AddCommand(newCompanyGroupShowCmd())
	cmd.AddCommand(newCompanyGroupUpdateCmd())
}

func newCompanyGroupCreateCmd() *cobra.Command {
	var opts options.CreateCompanyGroupOptions
	cmd := &cobra.Command{
		Use:   "companyGroup-create [OPTIONS] groupName groupStatus",
		Short: "Create a company group",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			opts.GroupName = args[0]
			if iStatus, err := strconv.Atoi(args[1]); err != nil {
				fmt.Println("change to int fail", args[1])
			} else {
				opts.GroupStatus = iStatus
			}

			if hv, err := Sdk.CreateCompanyGroup(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(hv)
			}
		},
	}
	return cmd
}

func newCompanyGroupDeleteCmd() *cobra.Command {
	return deleteCmd(resource_type_company_group, Sdk.DeleteCompanyGroup)
}

func newCompanyGroupListCmd() *cobra.Command {
	defCs := []string{"CompanyGroupID", "GroupName", "GroupStatus", "CreatedAt", "UpdatedAt"}
	cmd := &cobra.Command{
		Use:   "companyGroup-list ",
		Short: "List company group Support Filter",
	}
	columns := cmd.Flags().StringArrayP("column", "c", defCs, "Columns to display")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		var opts options.ListOptions
		opts.Limit = -1
		opts.Offset = 0
		if _, views, err := Sdk.ListCompanyGroup(&opts); err != nil {
			util.FormatErrorOutput(err)
		} else {
			util.FormatListOutput(*columns, views)
		}
	}
	return cmd
}

func newCompanyGroupShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "companyGroup-show [OPTIONS] companyGroupId",
		Short: "Show company group",
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
			view, err := Sdk.GetCompanyGroup(&opts)
			if err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(view)
			}
		},
	}
	return cmd
}

func newCompanyGroupUpdateCmd() *cobra.Command {
	var opts options.ModifyCompanyGroupOptions
	cmd := &cobra.Command{
		Use:   "company-update [OPTIONS] companyGroupId groupStatus",
		Short: "update a accSub",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			if id, err := strconv.Atoi(args[0]); err != nil {
				fmt.Println("change to int fail", args[0])
			} else {
				opts.CompanyGroupID = id
			}
			if iStatus, err := strconv.Atoi(args[1]); err != nil {
				fmt.Println("change to int fail", args[1])
			} else {
				opts.GroupStatus = iStatus
			}

			if err := Sdk.UpdateCompanyGroup(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	cmd.Flags().StringVar(&opts.GroupName, "comGroupName", "", "Company Group Name")
	return cmd
}
