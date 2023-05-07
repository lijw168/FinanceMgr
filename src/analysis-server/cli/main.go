package main

import (
	"financeMgr/src/analysis-server/cli/command"
	"financeMgr/src/analysis-server/sdk"
	"financeMgr/src/common/tag"
	"flag"
	"os"
)

func main() {
	flag.Parse()
	if tag.CheckAndShowVersion() {
		return
	}

	command.Sdk = &sdk.CcSdk{}
	cmd := command.NewCcCli()
	if command.Help {
		cmd.Execute()
		return
	}
	// Read environments
	if command.Domain == "" {
		command.Domain = os.Getenv("CC_SERVER_URL")
	}
	// if command.Tenant == "" {
	// 	command.Tenant = os.Getenv("CC_TENANT_ID")
	// }
	// Check required environments
	if command.Domain == "" {
		println("ERROR (CommandError): You must provide server_url via --server-url  or  env[CC_SERVER_URL]")
		return
	}

	command.Sdk.Domain = command.Domain
	//command.Sdk.Tenant = command.Tenant
	command.Sdk.Verbose = command.Verbose
	command.Sdk.Admin = command.Admin
	command.Sdk.Timeout = command.Timeout
	command.Sdk.Setup()
	command.Sdk.SetAccessToken(command.AccessToken)
	cmd.Execute()
}
