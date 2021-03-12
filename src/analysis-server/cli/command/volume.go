package command

import (
	"strconv"

	"github.com/spf13/cobra"
	"analysis-server/cli/util"
	"analysis-server/sdk/options"
	sdk_util "analysis-server/sdk/util"
)

func newVolumeCreateCmd() *cobra.Command {
	var opts options.CreateVolumeOptions
	cmd := &cobra.Command{
		Use:   "volume-create [OPTIONS] volumeName size volumeTypeName AzName",
		Short: "Create a volume",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 4 {
				cmd.Help()
				return
			}
			opts.TenantId = sdk_util.Tenant
			opts.VolumeName = args[0]
			var err error
			opts.Size, err = strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				util.FormatErrorOutput(err)
				return
			}
			opts.VolumeTypeName = args[2]
			opts.AzName = args[3]
			if volumeView, err := Sdk.CreateVolume(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(volumeView)
			}
		},
	}
	cmd.Flags().StringVar(&opts.Description, "description", "", "the descrition of volume")
	cmd.Flags().StringVar(&opts.SnapshotId, "snapshot_id", "", "create volume by the snapshot")
	cmd.Flags().Uint8Var(&opts.Bootable, "tag", 0, "0. data 1. bootable (invisible) 2. bootable ")
	return cmd
}
func newChangeVolumeCmd() *cobra.Command {
	var opts options.ChangeVolumeOptions
	cmd := &cobra.Command{
		Use:   "volume-change [OPTIONS] volumeId ",
		Short: "change volume",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			opts.TenantId = sdk_util.Tenant
			opts.VolumeId = args[0]
			if err := Sdk.ChangeVolume(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	opts.Description = new(string)
	opts.VolumeName = new(string)
	cmd.Flags().StringVar(opts.Description, "description", "", "the descrition of volume")
	cmd.Flags().StringVar(opts.VolumeName, "volume_name", "", "the name of volume")
	return cmd
}
func newVolumeResizeCmd() *cobra.Command {
	var opts options.ResizeVolumeOptions
	cmd := &cobra.Command{
		Use:   "volume-resize [OPTIONS] volumeId size",
		Short: "resize a volume",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}
			opts.TenantId = sdk_util.Tenant
			opts.VolumeId = args[0]
			var err error
			opts.Size, err = strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				util.FormatErrorOutput(err)
				return
			}
			if volumeView, err := Sdk.ResizeVolume(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(volumeView)
			}
		},
	}
	return cmd
}
func newVolumeDeleteCmd() *cobra.Command {
	var opts options.DeleteVolumeOptions
	cmd := &cobra.Command{
		Use:   "volume-delete [OPTIONS] volumeId ",
		Short: "delete a volume",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			opts.TenantId = sdk_util.Tenant
			opts.VolumeId = args[0]
			if err := Sdk.DeleteVolume(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	return cmd
}
func newDescribeVolumeCmd() *cobra.Command {
	var opts options.DescribeVolumeOptions
	cmd := &cobra.Command{
		Use:   "volume-describe [OPTIONS] volumeId ",
		Short: "detail a volume",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			opts.TenantId = sdk_util.Tenant
			opts.VolumeId = args[0]
			if volumeView, err := Sdk.DescribeVolume(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(volumeView)
			}
		},
	}
	return cmd
}

func newListVolumeCmd() *cobra.Command {
	defCs := []string{"Id", "VolumeName", "Size", "PoolId", "VolumeTypeName", "AzName", "Action", "TenantId", "CreatedAt", "UpdatedAt", "DeletedAt", "Status", "Iops", "Description"}
	cmd := &cobra.Command{
		Use:   "volume-list",
		Short: "List Volume Support Filter",
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
		if _, volumeViews, err := Sdk.ListVolumes(&opts); err != nil {
			util.FormatErrorOutput(err)
		} else {
			util.FormatListOutput(*columns, volumeViews)
		}
	}
	return cmd
}

func newVolumeAttachCmd() *cobra.Command {
	var opts options.AttachVolumeOptions
	cmd := &cobra.Command{
		Use:   "volume-attach [OPTIONS] hostIp volumeId instanceUuid [instanceType]",
		Short: "attach a volume",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 3 {
				cmd.Help()
				return
			}
			opts.TenantId = sdk_util.Tenant
			opts.HostIp = args[0]
			opts.VolumeId = args[1]
			opts.InstanceUuid = args[2]
			if len(args) > 3 {
				opts.InstanceType = args[3] // optional
			}
			if attachResultView, err := Sdk.AttachVolume(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(attachResultView)
			}
		},
	}
	cmd.Flags().StringVar(&opts.AttachMode, "attachMode", "rw", "attach mode")
	//cmd.Flags().StringVar(&opts.InstanceUuid, "instanceUuid", "", "instanceUuid")
	cmd.Flags().BoolVar(&opts.MultiAttachment, "multi_attachment", false, "multiple attachment")
	return cmd
}
func newVolumeDetachCmd() *cobra.Command {
	var opts options.DetachVolumeOptions
	cmd := &cobra.Command{
		Use:   "volume-detach [OPTIONS] volumeId attachmentId",
		Short: "detach a volume",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}
			opts.TenantId = sdk_util.Tenant
			opts.VolumeId = args[0]
			opts.AttachmentId = args[1]
			if err := Sdk.DetachVolume(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	return cmd
}
func newDescribeAttachmentCmd() *cobra.Command {
	var opts options.DescribeAttachmentOptions
	cmd := &cobra.Command{
		Use:   "attachment-describe [OPTIONS] attachmentId ",
		Short: "detail a attachment",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			opts.TenantId = sdk_util.Tenant
			opts.AttachmentId = args[0]
			if volumeAttView, err := Sdk.DescribeAttachment(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(volumeAttView)
			}
		},
	}
	return cmd
}

func newListVolumeAttCmd() *cobra.Command {
	defCs := []string{"Id", "VolumeId", "TenantId", "HostIp", "InstanceUuid", "DeviceName", "AttachMode", "AttachTime", "DetachTime", "Status"}
	cmd := &cobra.Command{
		Use:   "attachment-list [OPTIONS]",
		Short: "List attachment Support Filter",
	}
	columns := cmd.Flags().StringArrayP("column", "c", defCs, "Columns to display")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		//		if len(args) < 1 {
		//			cmd.Help()
		//			return
		//		}
		var opts options.ListOptions
		opts.Limit = -1
		opts.Offset = 0
		if _, volumeAttViews, err := Sdk.ListAttachments(&opts); err != nil {
			util.FormatErrorOutput(err)
		} else {
			util.FormatListOutput(*columns, volumeAttViews)
		}
	}
	return cmd
}
func NewVolumeCommand(cmd *cobra.Command) {
	cmd.AddCommand(newVolumeCreateCmd())
	cmd.AddCommand(newVolumeDeleteCmd())
	cmd.AddCommand(newDescribeVolumeCmd())
	cmd.AddCommand(newListVolumeCmd())
	cmd.AddCommand(newVolumeAttachCmd())
	cmd.AddCommand(newVolumeDetachCmd())
	cmd.AddCommand(newDescribeAttachmentCmd())
	cmd.AddCommand(newListVolumeAttCmd())
	cmd.AddCommand(newVolumeResizeCmd())
	cmd.AddCommand(newChangeVolumeCmd())
}
