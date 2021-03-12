package command

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	cons "common/constant"
	"analysis-server/cli/util"
	"analysis-server/model"
	"analysis-server/sdk/options"
)

type HostDisplay struct {
	Id          string
	HostName    string
	ManageIp    string
	StorageIp   string
	ClientIp    string
	TraceIp     string
	RackName    string
	PoolName    string
	CurPort     int
	Heartbeat   string
	AdminStatus string
	Status      string
}

func newHostCreateCmd() *cobra.Command {
	var opts options.CreateHostOptions
	cmd := &cobra.Command{
		Use:   "host-create NAME RACK_NAME POOL_NAME CLIENT_IP MANAGE_IP STORAGE_IP",
		Short: "Create a host",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}
			opts.Name = args[0]
			rackName := args[1]
			poolName := args[2]

			var listOpts options.ListOptions
			listOpts.Limit = -1
			listOpts.Offset = 0
			listOpts.Filter = make(map[string]interface{})
			listOpts.Filter["tag"] = rackName

			if numRacks, rackView, err := Sdk.DescribeRacks(&listOpts); err != nil {
				util.FormatErrorOutput(err)
				return
			} else if numRacks == 0 {
				util.FormatMessageOutput("Rack " + rackName + " doesn't exists")
				return
			} else {
				opts.Rack = rackView[0].Id
			}

			delete(listOpts.Filter, "tag")
			listOpts.Filter["name"] = poolName

			if numPools, poolView, err := Sdk.DescribePools(&listOpts); err != nil {
				util.FormatErrorOutput(err)
				return
			} else if numPools == 0 {
				util.FormatMessageOutput("Pool " + poolName + " doesn't exit")
				return
			} else {
				opts.PoolId = poolView[0].Id
			}

			if hv, err := Sdk.CreateHost(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				hd, err := displayHost(hv)
				if err != nil {
					util.FormatErrorOutput(err)
					return
				}
				util.FormatViewOutput(hd)
			}
		},
	}
	cmd.Flags().IntVar(&opts.AdminStatus, "admin_status", 1, "Admin Status for the port")
	cmd.Flags().StringVar(&opts.MgmtIp, "manage_ip", "", "Management ip of the host")
	cmd.Flags().StringVar(&opts.DataIp, "storage_ip", "", "Data ip of the host")
	cmd.Flags().StringVar(&opts.ClientIp, "client_ip", "", "Client ip of the host")
	return cmd
}

func newListHostCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "host-list [OPTIONS]",
		Short: "List Host Support Filter",
	}
	defCs := []string{"Id", "HostName", "ManageIp", "StorageIp", "ClientIp", "RackName", "PoolName", "CurPort", "Heartbeat", "AdminStatus", "Status"}
	name := cmd.Flags().String("name", "", "host name")
	adminstatus := cmd.Flags().String("adminstatus", "", "AdminStatus of the host. {UP | DOWN}")
	rack := cmd.Flags().String("rack", "", "rack name which host located at.")
	hg := cmd.Flags().String("hostgroup-id", "", "hostgroup name which host located at.")
	columns := cmd.Flags().StringArrayP("column", "c", defCs, "Columns to display")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		if len(args) < 0 {
			cmd.Help()
			return
		}
		var opts options.ListOptions
		opts.Limit = -1
		opts.Offset = 0
		if *name != "" || *adminstatus != "" || *rack != "" || *hg != "" {
			filters := make(map[string]interface{})
			if *name != "" {
				filters["name"] = *name
			}
			if *adminstatus != "" {
				ascode, err := encodeHostAdminStatus(*adminstatus)
				if err != nil {
					util.FormatErrorOutput(err)
					return
				}
				filters["admin_status"] = ascode
			}
			if *rack != "" {
				filters["rack"] = *rack
			}
			opts.Filter = filters
		}
		if _, views, err := Sdk.DescribeHosts(&opts); err != nil {
			util.FormatErrorOutput(err)
		} else {
			hds, err := displayHostList(views)
			if err != nil {
				util.FormatErrorOutput(err)
				return
			}
			util.FormatListOutput(*columns, hds)
		}
	}
	return cmd
}

func newModifyAdminStatusCmd() *cobra.Command {
	var opts options.ModifyAdminStatusOptions
	cmd := &cobra.Command{
		Use:   "host-adminstatus-modify [OPTIONS] NAME {DOWN | UP | GRAY}",
		Short: "modify host adminstatus.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}
			opts.Name = args[0]
			opts.AdminStatus = args[1]
			if isModify, err := Sdk.ModifyAdminStaus(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				hd, err := displayHost(isModify)
				if err != nil {
					util.FormatErrorOutput(err)
					return
				}
				util.FormatViewOutput(hd)
			}
		},
	}
	return cmd
}

func newModifyStatusCmd() *cobra.Command {
	var opts options.ModifyStatusOptions
	cmd := &cobra.Command{
		Use:   "host-status-modify [OPTIONS] NAME {DOWN | UP}",
		Short: "modify host status.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}
			opts.Name = args[0]
			_status, err := encodeHostStatus(args[1])
			if err != nil {
				util.FormatErrorOutput(err)
				return
			}
			opts.Status = _status
			if err := Sdk.ModifyStatus(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatMessageOutput("Modify host status successfully.")
			}
		},
	}
	cmd.Flags().BoolVarP(&opts.Force, "force", "f", false, "Modify host's status forcedly.")
	return cmd
}

func newDescribeHostCmd() *cobra.Command {
	var opts options.DescribeHostOptions
	cmd := &cobra.Command{
		Use:   "host-show [OPTIONS] ID",
		Short: "show host details.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.Help()
				return
			}
			opts.Id = args[0]
			if hv, err := Sdk.DescribeHost(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				hd, err := displayHost(hv)
				if err != nil {
					util.FormatErrorOutput(err)
					return
				}
				util.FormatViewOutput(hd)
			}
		},
	}
	return cmd
}

func newDeleteHostCmd() *cobra.Command {
	var opts options.DeleteHostOptions
	cmd := &cobra.Command{
		Use:   "host-delete ID",
		Short: "delete host.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.Help()
				return
			}
			opts.Id = args[0]
			if hv, err := Sdk.DeleteHost(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(hv)
			}
		},
	}
	return cmd
}
func NewHostCommand(cmd *cobra.Command) {
	cmd.AddCommand(newHostCreateCmd())
	cmd.AddCommand(newListHostCmd())
	cmd.AddCommand(newModifyAdminStatusCmd())
	cmd.AddCommand(newModifyStatusCmd())
	cmd.AddCommand(newDescribeHostCmd())
	cmd.AddCommand(newDeleteHostCmd())
}

func displayHostList(hvs []*model.HostView) ([]*HostDisplay, error) {
	hds := make([]*HostDisplay, 0)
	for _, hv := range hvs {
		hd, err := displayHost(hv)
		if err != nil {
			return nil, err
		}
		hds = append(hds, hd)
	}
	return hds, nil
}

func displayHost(hv *model.HostView) (*HostDisplay, error) {
	ass, err := parseHostAdminStatus(hv.AdminStatus)
	if err != nil {
		return nil, err
	}
	ss, err := parsteHostStatus(hv.Status)
	if err != nil {
		return nil, err
	}

	opt := &options.BaseOptions{}
	opt.Id = hv.Rack
	rv, _ := Sdk.DescribeRack(opt)

	opt.Id = hv.PoolId
	pv, _ := Sdk.DescribePool(opt)
	hd := &HostDisplay{
		Id:          hv.Id,
		HostName:    hv.Name,
		ManageIp:    hv.MgmtIp,
		StorageIp:   hv.DataIp,
		ClientIp:    hv.ClientIp,
		TraceIp:     hv.TraceIp,
		RackName:    rv.Tag,
		PoolName:    pv.Name,
		CurPort:     hv.CurPort,
		Heartbeat:   hv.Heartbeat.Format(time.RFC3339),
		AdminStatus: ass,
		Status:      ss,
	}
	return hd, nil
}

func parseHostAdminStatus(status int) (string, error) {
	switch status {
	case cons.HOST_ADMINSTATUS_DOWN:
		return "DOWN", nil
	case cons.HOST_ADMINSTATUS_UP:
		return "UP", nil
	case cons.HOST_ADMINSTATUS_GRAY:
		return "GRAY", nil
	default:
		return fmt.Sprintf("Unknown(%d)", status), nil
	}
}

func encodeHostAdminStatus(statusstr string) (int, error) {
	switch strings.ToUpper(statusstr) {
	case "DOWN":
		return cons.HOST_ADMINSTATUS_DOWN, nil
	case "UP":
		return cons.HOST_ADMINSTATUS_UP, nil
	case "GRAY":
		return cons.HOST_ADMINSTATUS_GRAY, nil
	default:
		return -1, errors.New("Bad Host AdminStatus")
	}
}

func parsteHostStatus(status int) (string, error) {
	switch status {
	case cons.HOST_STATUS_DOWN:
		return "DOWN", nil
	case cons.HOST_STATUS_UP:
		return "UP", nil
	default:
		return fmt.Sprintf("Unknown(%d)", status), nil
	}
}

func encodeHostStatus(statusstr string) (int, error) {
	switch strings.ToUpper(statusstr) {
	case "DOWN":
		return cons.HOST_STATUS_DOWN, nil
	case "UP":
		return cons.HOST_STATUS_UP, nil
	default:
		return -1, errors.New("Bad Host Status")
	}
}
