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
		Use:   "accSub-create [OPTIONS] subject_name subject_level",
		Short: "Create a accSub",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			var (
				err      error
				subLevel int
			)

			opts.SubjectName = args[0]

			if subLevel, err = strconv.Atoi(args[1]); err != nil {
				fmt.Println("change to int fail", args[1])
			}
			opts.SubjectLevel = subLevel

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
	var opts options.BaseOptions
	cmd := &cobra.Command{
		Use:   "accSub-delete [OPTIONS] ID",
		Short: "Delete a accountSubject",
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

			if err := Sdk.DeleteAccSub(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	return cmd
}

func newAccSubListCmd() *cobra.Command {
	defCs := []string{"Id", "subject_name", "subject_level"}
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
		//opts.Filter = make(map[string]interface{})
		//opts.Filter["status"] = "creating|available|in-use"
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
		Use:   "accSub update [OPTIONS] subjectID subject_name subject_level ",
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
			opts.SubjectName = args[1]

			if subLevel, err := strconv.Atoi(args[2]); err != nil {
				fmt.Println("change to int fail", args[2])
			} else {
				opts.SubjectLevel = subLevel
			}
			if err := Sdk.UpdateAccSub(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	return cmd
}
