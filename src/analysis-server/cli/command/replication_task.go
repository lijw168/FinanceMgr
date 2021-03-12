package command

import (
	"github.com/spf13/cobra"
	"analysis-server/cli/util"
	"analysis-server/sdk/options"
	"strings"
)

func NewRTCommand(cmd *cobra.Command) {
	cmd.AddCommand(newReschedulerCmd())
	cmd.AddCommand(newTransferLeaderCmd())
	cmd.AddCommand(newMoveReplicaCmd())
	cmd.AddCommand(newDeleteReplicaCmd())
	cmd.AddCommand(newAddReplicaCmd())
}

func newReschedulerCmd() *cobra.Command {
	var opts options.ReschedulerOptions
	cmd := &cobra.Command{
		Use:   "Reschedule PoolId",
		Short: "Create a Reschedule Task",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}

			opts.PoolId = args[0]

			if hv, err := Sdk.Rescheduler(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(hv)
			}
		},
	}
	return cmd
}

func newTransferLeaderCmd() *cobra.Command {
	var opts options.TransferLeaderOptions
	cmd := &cobra.Command{
		Use:   "TransferLeader ReplicationGroupId ReplicaId",
		Short: "Create a TransferLeader Task",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}

			opts.RgId = args[0]
			opts.ReplicaId = args[1]

			if hv, err := Sdk.TransferLeader(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(hv)
			}
		},
	}
	return cmd
}

func newMoveReplicaCmd() *cobra.Command {
	var opts options.MoveReplicaOptions
	cmd := &cobra.Command{
		Use:   "MoveReplica ReplicaId TargetDiskId",
		Short: "Create a MoveReplica Task",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}

			opts.ReplicaId = args[0]
			opts.TargetDiskId = args[1]

			if hv, err := Sdk.MoveReplica(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(hv)
			}
		},
	}
	return cmd
}

func newDeleteReplicaCmd() *cobra.Command {
	var opts options.DeleteReplicaOptions
	cmd := &cobra.Command{
		Use:   "DeleteReplica ReplicaId [<force:|true|false>]",
		Short: "Create a DeleteReplica Task.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}

			opts.ReplicaId = args[0]
			if len(args) > 1 {
				strforce := strings.ToLower(args[1])
				if strforce == "true" {
					opts.Force = true
				} else if strforce == "false" {
					opts.Force = false
				}
			} else {
				opts.Force = false
			}

			if hv, err := Sdk.DeleteReplica(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(hv)
			}
		},
	}
	return cmd
}

func newAddReplicaCmd() *cobra.Command {
	var opts options.AddReplicaOptions
	cmd := &cobra.Command{
		Use:   "AddReplica <RgID> <DiskID>",
		Short: "Create a AddReplica Task.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				cmd.Help()
				return
			}

			opts.RgID = args[0]
			opts.DiskID = args[1]

			if hv, err := Sdk.AddReplica(&opts); err != nil {
				util.FormatErrorOutput(err)
			} else {
				util.FormatViewOutput(hv)
			}
		},
	}
	return cmd
}
