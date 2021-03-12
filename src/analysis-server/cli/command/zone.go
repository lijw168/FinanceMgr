package command

import (
	"github.com/spf13/cobra"
	"common/constant"
	"analysis-server/cli/util"
	"analysis-server/sdk/options"

	"fmt"
	"strconv"
)

func NewPoolCommand(cmd *cobra.Command) {
	cmd.AddCommand(newPoolCreateCmd())
	cmd.AddCommand(newPoolDeleteCmd())
	cmd.AddCommand(newPoolListCmd())
	cmd.AddCommand(newPoolShowCmd())
	cmd.AddCommand(newPoolUpdateStatusCmd())
}

func newPoolCreateCmd() *cobra.Command {
	var opts options.CreatePoolOptions
	cmd := &cobra.Command{
		Use: `pool-create NAME TYPE REPLICATION_GROUP_COUNT OBJECT_SIZE MEDIA_TYPE 
		TYPE (EBS : 1, PRERELEASE : 2, GRAYSCALE : 3, PRODUCTION : 4, SPECIAL : 5)
        MEDIA_TYPE (ssd : 1, hdd : 2)`,
		Short: "Create a Pool",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 5 {
				cmd.Help()
				return
			}
			var (
				err     error
				ObjSize int
			)

			opts.Name = args[0]
			if opts.Type, err = strconv.Atoi(args[1]); err != nil {
				fmt.Println("change to int fail", args[1])
				return
			}

			if opts.Rgc, err = strconv.Atoi(args[2]); err != nil {
				fmt.Println("change to int fail", args[2])
				return
			}

			if ObjSize, err = strconv.Atoi(args[3]); err != nil {
				fmt.Println("change to int fail", args[3])
				return
			}
			opts.ObjSize = uint64(ObjSize)
			opts.MediaType = args[4]
			if args[4] != constant.MEDIA_TYPE_HDD && args[4] != constant.MEDIA_TYPE_SSD {
				fmt.Println("invalid media type")
				return
			}

			if hv, err := Sdk.CreatePool(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(hv)
			}
		},
	}
	return cmd
}

func newPoolDeleteCmd() *cobra.Command {
	//return deleteCmd(RESOURCE_TYPE_POOL, Sdk.DeletePool)
	var opts options.DeleteOptions
	cmd := &cobra.Command{
		Use:   "pool-delete ID",
		Short: "Delete a Pool",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}

			opts.Id = args[0]

			if hv, err := Sdk.DeletePool(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(hv)
			}
		},
	}
	return cmd
}

func newPoolListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool-list [OPTIONS]",
		Short: "List Pool. Support Filter",
	}
	defCs := []string{"Id", "Name", "Type", "Rgc", "ObjSize", "Status"}
	columns := cmd.Flags().StringArrayP("column", "c", defCs, "Columns to display")
	name := cmd.Flags().String("name", "", "Pool name")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		var opts options.ListOptions
		opts.Limit = -1
		opts.Offset = 0
		opts.Filter = make(map[string]interface{})
		if *name != "" {
			opts.Filter["name"] = *name
		}

		if _, views, err := Sdk.DescribePools(&opts); err != nil {
			util.FormatErrorOutput(err)
		} else {
			util.FormatListOutput(*columns, views)
		}
	}
	return cmd
}

func newPoolShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool-show POOL",
		Short: "Show given pool",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			var opts options.BaseOptions
			opts.Id = args[0]
			view, err := Sdk.DescribePool(&opts)
			if err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(view)
			}
		},
	}
	return cmd
}

func newPoolUpdateStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: `pool-update POOL STATUS 
		STATUS (1: ENABLED, 0: DISABLED)`,
		Short: "update pool status",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				cmd.Help()
				return
			}
			var err error
			var opts options.UpdateStatusOptions
			opts.Id = args[0]
			if opts.Status, err = strconv.Atoi(args[1]); err != nil {
				fmt.Println("change to int fail", args[1])
				return
			}
			view, err := Sdk.UpdatePoolStatus(&opts)
			if err != nil {
				util.FormatErrorOutput(err)
			} else if view != nil {
				util.FormatViewOutput(view)
			}
		},
	}
	return cmd
}
