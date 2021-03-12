package command

import (
	"strconv"

	"github.com/spf13/cobra"
	"analysis-server/cli/util"
	"analysis-server/sdk/options"
	sdk_util "analysis-server/sdk/util"
)

func newSnapshotCreateCmd() *cobra.Command {
	var opts options.CreateSnapshotOptions
	cmd := &cobra.Command{
		Use:   "snapshot-create [OPTIONS] TenantId SnapshotName volumeId",
		Short: "Create a snapshot",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 3 {
				cmd.Help()
				return
			}
			opts.TenantId = args[0]
			opts.SnapshotName = args[1]
			opts.VolumeId = args[2]
			if snapshotView, err := Sdk.CreateSnapshot(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(snapshotView)
			}
		},
	}
	cmd.Flags().StringVar(&opts.Description, "description", "", "the descrition of volume")
	return cmd
}
func newConvertImageToSnapshotCmd() *cobra.Command {
	var opts options.ConvertImageToSnapshotOptions
	cmd := &cobra.Command{
		Use:   "image-convert [OPTIONS] TenantId SnapshotName ImageId ImageFormat ImageLocation ImageType ImageHash PublicImage",
		Short: "Convert image file to snapshot",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 8 {
				cmd.Help()
				return
			}
			opts.TenantId = args[0]
			opts.SnapshotName = args[1]
			opts.ImageId = args[2]
			opts.ImageFormat = args[3]
			opts.ImageLocation = args[4]
			opts.ImageType = args[5]
			opts.ImageHash = args[6]
			if pbImage, err := strconv.ParseBool(args[7]); err == nil {
				opts.PublicImage = pbImage
				if snapshotView, err := Sdk.ConvertImageToSnapshot(&opts); err != nil {
					util.FormatErrorOutput(err)
				} else {
					util.FormatViewOutput(snapshotView)
				}
			} else {
				util.FormatErrorOutput(err)
			}
		},
	}
	opts.AzName = new(string)
	cmd.Flags().StringVar(opts.AzName, "az_name", "", "the image is to process in the zone")
	return cmd
}
func newVolumeRestoreCmd() *cobra.Command {
	var opts options.RestoreSnapshotOptions
	cmd := &cobra.Command{
		Use:   "volume-restore [OPTIONS] TenantId SnapshotId",
		Short: "Restore a volume",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}
			opts.TenantId = args[0]
			opts.SnapshotId = args[1]
			if err := Sdk.RestoreVolume(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	return cmd
}

func newSnapshotDeleteCmd() *cobra.Command {
	var opts options.DeleteSnapshotOptions
	cmd := &cobra.Command{
		Use:   "snapshot-delete [OPTIONS] tenantId snapshotId ",
		Short: "delete a snapshot",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}
			opts.TenantId = args[0]
			opts.SnapshotId = args[1]
			if err := Sdk.DeleteSnapshot(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	return cmd
}
func newDescribeSnapshotCmd() *cobra.Command {
	var opts options.DescribeSnapshotOptions
	cmd := &cobra.Command{
		Use:   "snapshot-describe [OPTIONS] tenantId snapshotId ",
		Short: "detail a snapshot",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}
			opts.TenantId = args[0]
			opts.SnapshotId = args[1]
			if snapshotView, err := Sdk.DescribeSnapshot(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(snapshotView)
			}
		},
	}
	return cmd
}

func newListSnapshotCmd() *cobra.Command {
	defCs := []string{"Id", "SnapshotName", "Size", "TenantId", "VolumeId", "PoolId", "AzName",
		"CreatedAt", "UpdatedAt", "DeletedAt", "Status", "Description", "Share"}
	cmd := &cobra.Command{
		Use:   "snapshot-list",
		Short: "List Snapshot Support Filter",
	}
	columns := cmd.Flags().StringArrayP("column", "c", defCs, "Columns to display")
	//var volumeType string
	//cmd.Flags().StringVar(&volumeType, "name", "", "the type name of volume")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		//		if len(args) < 1 {
		//			cmd.Help()
		//			return
		//		}
		var opts options.ListOptions
		opts.Limit = -1
		opts.Offset = 0
		//opts.Filter["name"] = volumeType
		if _, snapshotViews, err := Sdk.ListSnapshots(&opts); err != nil {
			util.FormatErrorOutput(err)
		} else {
			util.FormatListOutput(*columns, snapshotViews)
		}
	}
	return cmd
}
func newChangeSnapshotCmd() *cobra.Command {
	var opts options.ChangeSnapshotOptions
	cmd := &cobra.Command{
		Use:   "snapshot-change [OPTIONS] snapshotId ",
		Short: "change snapshot",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			opts.TenantId = sdk_util.Tenant
			opts.SnapshotId = args[0]
			if err := Sdk.ChangeSnapshot(&opts); err != nil {
				util.FormatErrorOutput(err)
			}
		},
	}
	cmd.Flags().StringVar(&opts.Description, "description", "", "the descrition of snapshot")
	cmd.Flags().StringVar(&opts.SnapshotName, "snapshot_name", "", "the name of snapshot")
	cmd.Flags().Int8Var(&opts.Share, "share", 0, "0(exclusive) 1(share)")
	return cmd
}

func NewSnapshotCommand(cmd *cobra.Command) {
	cmd.AddCommand(newSnapshotCreateCmd())
	cmd.AddCommand(newConvertImageToSnapshotCmd())
	cmd.AddCommand(newSnapshotDeleteCmd())
	cmd.AddCommand(newDescribeSnapshotCmd())
	cmd.AddCommand(newListSnapshotCmd())
	cmd.AddCommand(newVolumeRestoreCmd())
	cmd.AddCommand(newChangeSnapshotCmd())
}
