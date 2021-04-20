package command

import (
	"analysis-server/cli/util"
	"analysis-server/sdk/options"
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

const (
	resource_type_account_sub = iota
	resource_type_company
	resource_type_voucher
	resource_type_voucher_record
)

func deleteCmd(rsc int, handler func(*options.BaseOptions) error) *cobra.Command {
	rscName := getResourceCmdName(rsc)
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("%s-delete [OPTION] ID", getResourceCmdName(rsc)),
		Short: "Delete a " + rscName,
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
			err := handler(&opts)
			if err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatMessageOutput("Delete " + rscName + " successfully.")
			}
		},
	}
	return cmd
}

func getResourceCmdName(rsc int) string {
	switch rsc {
	case resource_type_account_sub:
		return "accSub"
	case resource_type_company:
		return "company"
	case resource_type_voucher:
		return "voucher"
	case resource_type_voucher_record:
		return "vouRecord"
	default:
		panic("Unsupport resource type")
	}
}
