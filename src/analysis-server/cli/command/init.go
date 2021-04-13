package command

import (
	"os"
	"strconv"

	"analysis-server/sdk"
	"github.com/spf13/cobra"
)

var (
	Verbose bool
	Admin   bool
	Domain  string
	//Tenant  string
	Sdk     *sdk.CcSdk
	Help    bool
	Timeout uint64
)

func NewCcCli() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ccs [OPTIONS] COMMAND [arg...]",
		Short: "A self-sufficient runtime for analysis system",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				cmd.Help()
				return
			}
		},
	}
	cmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	cmd.PersistentFlags().BoolVarP(&Admin, "admin", "a", false, "admin role")
	cmd.PersistentFlags().StringVar(&Domain, "server-url", "", "URL of web server api")
	//cmd.PersistentFlags().StringVar(&Tenant, "tenant-id", "", "Tenant ID used to call api")
	cmd.PersistentFlags().BoolVarP(&Help, "help", "h", false, "Type for help")
	cmd.PersistentFlags().Uint64Var(&Timeout, "timeout", 0, "network timeout (ms)")

	al := len(os.Args)
	for i := 0; i < al; i++ {
		switch os.Args[i] {
		// case "--tenant-id":
		// 	if i < al-1 {
		// 		Tenant = os.Args[i+1]
		// 		i++
		// 	}
		case "--server-url":
			if i < al-1 {
				Domain = os.Args[i+1]
				i++
			}
		case "--verbose", "-v":
			if i < al-1 {
				if os.Args[i+1] != "false" {
					Verbose = true
				}
			} else {
				Verbose = true
			}
		case "--admin", "-a":
			if i < al-1 {
				if os.Args[i+1] != "false" {
					Admin = true
				}
			} else {
				Admin = true
			}
		case "--help", "-h":
			Help = true
		case "--timeout":
			t, _ := strconv.ParseUint(os.Args[i+1], 10, 64)
			Timeout = t
		}
	}
	// add command
	NewAccSubCommand(cmd)
	NewCompanyCommand(cmd)
	NewOperatorCommand(cmd)
	NewVoucherCommand(cmd)

	return cmd
}
