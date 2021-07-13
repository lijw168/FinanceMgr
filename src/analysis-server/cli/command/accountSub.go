package command

import (
	"analysis-server/cli/util"
	"analysis-server/sdk/options"
	"github.com/spf13/cobra"

	//"analysis-server/model"
	"fmt"
	"strconv"
)

func NewAccSubCommand(cmd *cobra.Command) {
	cmd.AddCommand(newAccSubCreateCmd())
	cmd.AddCommand(newAccSubDeleteCmd())
	cmd.AddCommand(newAccSubListCmd())
	cmd.AddCommand(newAccSubShowCmd())
	cmd.AddCommand(newAccSubUpdateCmd())
}

func newAccSubCreateCmd() *cobra.Command {
	var opts options.CreateSubjectOptions
	cmd := &cobra.Command{
		Use:   "accSub-create [OPTIONS] common_id subject_name subject_level companyId",
		Short: "Create a accSub",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			var (
				err       error
				subLevel  int
				companyID int
			)

			opts.CommonID = args[0]
			opts.SubjectName = args[1]

			if subLevel, err = strconv.Atoi(args[2]); err != nil {
				fmt.Println("change to int fail", args[2])
			}
			opts.SubjectLevel = subLevel
			if companyID, err = strconv.Atoi(args[3]); err != nil {
				fmt.Println("change to int fail", args[3])
			}
			opts.CompanyID = companyID

			if hv, err := Sdk.CreateAccSub(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(hv)
			}
		},
	}
	return cmd
}

func newAccSubDeleteCmd() *cobra.Command {
	return deleteCmd(resource_type_account_sub, Sdk.DeleteAccSub)
}

func newAccSubListCmd() *cobra.Command {
	defCs := []string{"SubjectID", "CommonID", "SubjectName", "SubjectLevel", "CompanyID"}
	cmd := &cobra.Command{
		Use:   "accSub-list",
		Short: "List account subjects Support Filter",
	}
	columns := cmd.Flags().StringArrayP("column", "c", defCs, "Columns to display")
	//var volumeType string
	//cmd.Flags().StringVar(&volumeType, "name", "", "the type name of volume")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		var opts options.ListOptions
		opts.Limit = -1
		opts.Offset = 0
		//for test
		// opts.Filter = make(map[string]interface{})
		// opts.Filter["subject_id"] = int(102)
		if _, accSubViews, err := Sdk.ListAccSub(&opts); err != nil {
			util.FormatErrorOutput(err)
		} else {
			util.FormatListOutput(*columns, accSubViews)
		}
	}
	return cmd
}

func newAccSubShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accSub-show [OPTIONS] subjectId",
		Short: "Show accSub",
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
			view, err := Sdk.GetAccSub(&opts)
			if err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(view)
			}
		},
	}
	return cmd
}

func newAccSubUpdateCmd() *cobra.Command {
	var opts options.ModifySubjectOptions
	cmd := &cobra.Command{
		Use:   "accSub-update [OPTIONS] subjectID commonID subject_name subject_level company_id ",
		Short: "update a accSub",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			if id, err := strconv.Atoi(args[0]); err != nil {
				fmt.Println("change to int fail", args[0])
			} else {
				opts.SubjectID = id
			}
			opts.CommonID = args[1]
			opts.SubjectName = args[2]

			if subLevel, err := strconv.Atoi(args[3]); err != nil {
				fmt.Println("change to int fail", args[3])
			} else {
				opts.SubjectLevel = subLevel
			}
			if companyID, err := strconv.Atoi(args[4]); err != nil {
				fmt.Println("change to int fail", args[4])
			} else {
				opts.CompanyID = companyID
			}

			if err := Sdk.UpdateAccSub(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	return cmd
}
