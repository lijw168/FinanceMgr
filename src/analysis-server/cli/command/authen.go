package command

import (
	"analysis-server/cli/util"
	"analysis-server/sdk/options"
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

func NewAuthenCommand(cmd *cobra.Command) {
	cmd.AddCommand(newLoginCmd())
	cmd.AddCommand(newLogoutCmd())
	cmd.AddCommand(newLoginListCmd())
	cmd.AddCommand(newLoginShowCmd())
}

func newLoginCmd() *cobra.Command {
	var opts options.AuthenInfoOptions
	cmd := &cobra.Command{
		Use:   "login [OPTIONS] name password companyID",
		Short: "user login",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			opts.Name = args[0]
			opts.Password = args[2]
			if id, err := strconv.Atoi(args[3]); err != nil {
				fmt.Println("change to int fail", args[0])
			} else {
				opts.CompanyID = id
			}

			if view, err := Sdk.Login(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(view)
			}
		},
	}
	return cmd
}

func newLogoutCmd() *cobra.Command {
	var opts options.NameOptions
	cmd := &cobra.Command{
		Use:   "logout [OPTIONS] name ",
		Short: "user logout",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			opts.Name = args[0]
			if err := Sdk.Logout(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	return cmd
}

func newLoginListCmd() *cobra.Command {
	defCs := []string{"Name", "Role", "ClientIp", "BeginedAt", "EndedAt"}
	cmd := &cobra.Command{
		Use:   "loginInfo-list ",
		Short: "List operators Support Filter",
	}
	columns := cmd.Flags().StringArrayP("column", "c", defCs, "Columns to display")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		var opts options.ListOptions
		opts.Limit = -1
		opts.Offset = 0
		if _, views, err := Sdk.ListLoginInfo(&opts); err != nil {
			util.FormatErrorOutput(err)
		} else {
			util.FormatListOutput(*columns, views)
		}
	}
	return cmd
}

func newLoginShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "loginInfo-show [OPTIONS] username",
		Short: "Show the information of logined user",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			var opts options.NameOptions
			opts.Name = args[0]
			view, err := Sdk.GetLoginInfo(&opts)
			if err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(view)
			}
		},
	}
	return cmd
}
