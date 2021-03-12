package command

import (
	"github.com/spf13/cobra"
	"analysis-server/cli/util"
	"analysis-server/sdk/options"

	"fmt"
	"analysis-server/model"
	"strconv"
)

func NewDiskCommand(cmd *cobra.Command) {
	cmd.AddCommand(newDiskCreateCmd())
	cmd.AddCommand(newDiskDeleteCmd())
	cmd.AddCommand(newDiskListCmd())
	cmd.AddCommand(newDiskShowCmd())
	cmd.AddCommand(newDiskMarkDownCmd())
	cmd.AddCommand(newDiskMarkUpCmd())
}

func newDiskCreateCmd() *cobra.Command {
	var opts options.CreateDiskOptions
	cmd := &cobra.Command{
		Use:   "disk-create DEVICE_ID HOST_NAME VOLUME_TYPE CAPACITY",
		Short: "Create a Disk",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			var (
				err      error
				Capacity int
			)

			opts.DeviceId = args[0]
			hostName := args[1]
			opts.VolumeTypeName = args[2]

			if Capacity, err = strconv.Atoi(args[3]); err != nil {
				fmt.Println("change to int fail", args[1])
			}
			opts.Capacity = uint64(Capacity)

			var listOpts options.ListOptions
			listOpts.Limit = -1
			listOpts.Offset = 0
			listOpts.Filter = make(map[string]interface{})
			listOpts.Filter["hostname"] = hostName

			if numHosts, hostView, err := Sdk.DescribeHosts(&listOpts); err != nil {
				util.FormatErrorOutput(err)
				return
			} else if numHosts == 0 {
				util.FormatMessageOutput("Host " + hostName + " doesn't exit")
				return
			} else {
				opts.HostId = hostView[0].Id
				opts.ManageAddr = hostView[0].MgmtIp
				opts.ClientAddr = hostView[0].ClientIp
				opts.StorageAddr = hostView[0].DataIp
				opts.TraceAddr = hostView[0].TraceIp
			}

			if hv, err := Sdk.CreateDisk(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(hv)
			}
		},
	}
	cmd.Flags().IntVar(&opts.AdminStatus, "admin_status", 1, "Admin Status for the port")
	return cmd
}

func newDiskDeleteCmd() *cobra.Command {
	var opts options.DeleteOptions
	cmd := &cobra.Command{
		Use:   "disk-delete ID",
		Short: "Delete a Disk",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}

			opts.Id = args[0]

			if hv, err := Sdk.DeleteDisk(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(hv)
			}
		},
	}
	return cmd
}

type DiskDisplay struct {
	Id             string
	DeviceId       string
	ManageAddr     string
	StorageAddr    string
	ClientAddr     string
	TraceAddr      string
	HostName       string
	VolumeTypeName string
	Capacity       uint64
	Free           uint64
	Version        uint64
	AdminStatus    string
	Status         string
}

func newDiskListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disk-list [OPTIONS]",
		Short: "List Disk. Support Filter",
	}
	defCs := []string{"Id", "DeviceId", "ManageAddr", "StorageAddr", "ClientAddr", "TraceAddr", "HostName", "VolumeTypeName", "Capacity", "Free", "Version", "AdminStatus", "Status"}
	columns := cmd.Flags().StringArrayP("column", "c", defCs, "Columns to display")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		var opts options.ListOptions
		opts.Limit = -1
		opts.Offset = 0
		opts.Filter = make(map[string]interface{})

		if _, views, err := Sdk.DescribeDisks(&opts); err != nil {
			util.FormatErrorOutput(err)
		} else {
			diskDisplays, err := displayDiskList(views)
			if err != nil {
				util.FormatErrorOutput(err)
				return
			}
			util.FormatListOutput(*columns, diskDisplays)
		}
	}
	return cmd
}

func displayDiskList(dvs []*model.DiskView) ([]*DiskDisplay, error) {
	diskDisplays := make([]*DiskDisplay, 0)
	for _, dv := range dvs {
		hd, err := displayDisk(dv)
		if err != nil {
			return nil, err
		}
		diskDisplays = append(diskDisplays, hd)
	}
	return diskDisplays, nil
}

func displayDisk(dv *model.DiskView) (*DiskDisplay, error) {
	ass, err := parseHostAdminStatus(dv.AdminStatus)
	if err != nil {
		return nil, err
	}
	ss, err := parsteHostStatus(dv.Status)
	if err != nil {
		return nil, err
	}

	opt := &options.DescribeHostOptions{}
	opt.Id = dv.HostId
	hv, _ := Sdk.DescribeHost(opt)

	diskDisplay := &DiskDisplay{
		Id:             dv.Id,
		DeviceId:       dv.DeviceId,
		ManageAddr:     dv.ManageAddr,
		StorageAddr:    dv.StorageAddr,
		ClientAddr:     dv.ClientAddr,
		TraceAddr:      dv.TraceAddr,
		HostName:       hv.Name,
		VolumeTypeName: dv.VolumeTypeName,
		Capacity:       dv.Capacity / 1024 / 1024 / 1024,
		Free:           dv.Free / 1024 / 1024 / 1024,
		Version:        dv.Version,
		AdminStatus:    ass,
		Status:         ss,
	}
	return diskDisplay, nil
}

func newDiskShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disk-show Disk",
		Short: "Show given disk",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
			var opts options.BaseOptions
			opts.Id = args[0]
			view, err := Sdk.DescribeDisk(&opts)
			if err != nil {
				util.FormatErrorOutput(err)
			} else {
				diskDisplay, err := displayDisk(view)
				if err != nil {
					util.FormatErrorOutput(err)
					return
				}
				util.FormatViewOutput(diskDisplay)
			}
		},
	}
	return cmd
}

func newDiskMarkDownCmd() *cobra.Command {
	var opts options.DeleteOptions
	cmd := &cobra.Command{
		Use:   "disk-markdown ID",
		Short: "MarkDown a Disk",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}

			opts.Id = args[0]

			if hv, err := Sdk.MarkDownDisk(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(hv)
			}
		},
	}
	return cmd
}

func newDiskMarkUpCmd() *cobra.Command {
	var opts options.DeleteOptions
	cmd := &cobra.Command{
		Use:   "disk-markup ID",
		Short: "MarkUp a Disk",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}

			opts.Id = args[0]

			if hv, err := Sdk.MarkUpDisk(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(hv)
			}
		},
	}
	return cmd
}
