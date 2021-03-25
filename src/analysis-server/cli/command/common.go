package command

import (
// "analysis-server/cli/util"
// "analysis-server/sdk/options"
// "errors"
// "fmt"
// "github.com/spf13/cobra"
)

// const (
// 	RESOURCE_TYPE_POOL = iota
// 	RESOURCE_TYPE_RACK
// 	RESOURCE_TYPE_HOST
// 	RESOURCE_TYPE_DISK
// 	RESOURCE_TYPE_RG
// 	RESOURCE_TYPE_PROXY
// 	RESOURCE_TYPE_AT
// 	RESOURCE_TYPE_CACHE
// )

// type cmdHandler func(cmd *cobra.Command, args []string)

// func updateCmd(rsc int, handler func(*options.ModifyAttributeOptions) error) *cobra.Command {
// 	rscName := getResourceName(rsc)
// 	cmd := &cobra.Command{
// 		Use:   fmt.Sprintf("%s-update [OPTION] ID", getResourceCmdName(rsc)),
// 		Short: "Update attributes of a " + rscName,
// 	}
// 	name := cmd.Flags().String("name", "", "Name for the "+rscName)
// 	descs := cmd.Flags().StringArray("description", nil, "Description of the "+rscName)
// 	cmd.Run = func(cmd *cobra.Command, args []string) {
// 		if len(args) < 1 {
// 			cmd.Help()
// 			return
// 		}
// 		if *name == "" && (*descs == nil || len(*descs) == 0) {
// 			util.FormatErrorOutput(errors.New("Least one of Name and Description are required"))
// 			return
// 		}
// 		var opts options.ModifyAttributeOptions
// 		opts.Id = args[0]
// 		if *name != "" {
// 			opts.Name = name
// 		}
// 		if *descs != nil && len(*descs) > 0 {
// 			opts.Description = &((*descs)[0])
// 		}
// 		err := handler(&opts)
// 		if err != nil {
// 			util.FormatErrorOutput(err)
// 		} else {
// 			util.FormatMessageOutput(rscName + " " + opts.Id + " has been updated successfully")
// 		}
// 	}
// 	return cmd
// }

// func deleteCmd(rsc int, handler func(*options.BaseOptions) error) *cobra.Command {
// 	rscName := getResourceName(rsc)
// 	cmd := &cobra.Command{
// 		Use:   fmt.Sprintf("%s-delete [OPTION] ID", getResourceCmdName(rsc)),
// 		Short: "Delete a " + rscName,
// 		Run: func(cmd *cobra.Command, args []string) {
// 			if len(args) < 1 {
// 				cmd.Help()
// 				return
// 			}
// 			var opts options.BaseOptions
// 			opts.Id = args[0]
// 			err := handler(&opts)
// 			if err != nil {
// 				util.FormatErrorOutput(err)
// 			} else {
// 				util.FormatMessageOutput("Delete " + rscName + " successfully.")
// 			}
// 		},
// 	}
// 	return cmd
// }

// func getResourceName(rsc int) string {
// 	switch rsc {
// 	case RESOURCE_TYPE_POOL:
// 		return "Pool"
// 	case RESOURCE_TYPE_RACK:
// 		return "Rack"
// 	case RESOURCE_TYPE_HOST:
// 		return "Host"
// 	case RESOURCE_TYPE_DISK:
// 		return "Disk"
// 	case RESOURCE_TYPE_RG:
// 		return "Rg"
// 	case RESOURCE_TYPE_PROXY:
// 		return "Proxy"
// 	case RESOURCE_TYPE_AT:
// 		return "At"
// 	case RESOURCE_TYPE_CACHE:
// 		return "Cache"
// 	default:
// 		panic("Unsupport resource type")
// 	}
// }

// func getResourceCmdName(rsc int) string {
// 	s := getResourceName(rsc)
// 	var tmp string
// 	for i, c := range s {
// 		if c >= 65 && c <= 90 {
// 			if i != 0 {
// 				tmp += "-" + string(c+32)
// 			} else {
// 				tmp += string(c + 32)
// 			}

// 		} else if c == 32 {
// 			continue
// 		} else {
// 			tmp += string(c)
// 		}
// 	}
// 	return tmp
// }
