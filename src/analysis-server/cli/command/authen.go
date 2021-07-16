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
	cmd.AddCommand(newStatusCheckoutCmd())
}

func newLoginCmd() *cobra.Command {
	var opts options.AuthenInfoOptions
	cmd := &cobra.Command{
		Use:   "login [OPTIONS] name password companyID",
		Short: "user login",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 3 {
				cmd.Help()
				return
			}
			opts.Name = args[0]
			opts.Password = args[1]
			if id, err := strconv.Atoi(args[2]); err != nil {
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
	cmd := &cobra.Command{
		Use:   "logout [OPTIONS] operatorID",
		Short: "user logout",
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
			if err := Sdk.Logout(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	return cmd
}

func newLoginListCmd() *cobra.Command {
	defCs := []string{"OperatorID", "Name", "Status", "ClientIp", "BeginedAt", "EndedAt"}
	cmd := &cobra.Command{
		Use:   "loginInfo-list companyId",
		Short: "List operators Support Filter",
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
		if id, err := strconv.Atoi(args[0]); err != nil {
			fmt.Println("change to int fail", args[0])
		} else {
			opts.Filter["companyId"] = id
		}
		if _, views, err := Sdk.ListLoginInfo(&opts); err != nil {
			util.FormatErrorOutput(err)
		} else {
			util.FormatListOutput(*columns, views)
		}
	}
	return cmd
}

//由于一个用户，可能有多条loginInfo记录，所以该函数通过list函数获取相应的数据。
func newLoginShowCmd() *cobra.Command {
	defCs := []string{"OperatorID", "Name", "Status", "ClientIp", "BeginedAt", "EndedAt"}
	cmd := &cobra.Command{
		Use:   "loginInfo-show [OPTIONS] operatorID",
		Short: "Show the information of logined user",
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
			opts.Filter["operatorId"] = id
		}

		if _, views, err := Sdk.ListLoginInfo(&opts); err != nil {
			util.FormatErrorOutput(err)
		} else {
			util.FormatListOutput(*columns, views)
		}
	}
	return cmd
}

func newStatusCheckoutCmd() *cobra.Command {
	defCs := []string{"OperatorID", "Name", "Status"}
	cmd := &cobra.Command{
		Use:   "status-checkout [OPTIONS] operatorID",
		Short: "checkout the user status",
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

		if views, err := Sdk.StatusCheckout(&opts); err != nil {
			util.FormatErrorOutput(err)
		} else {
			util.FormatListOutput(*columns, views)
		}
	}
	return cmd
}
