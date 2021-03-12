package main

import (
	"flag"
	"os"

	//_ "common/tag"
	"common/tag"
	"analysis-server/cli/command"
	"analysis-server/sdk"
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
	if command.Tenant == "" {
		command.Tenant = os.Getenv("CC_TENANT_ID")
	}
	// Check required environments
	if command.Domain == "" {
		println("ERROR (CommandError): You must provide server_url via --server-url  or  env[CC_SERVER_URL]")
		return
	}
	if !command.Admin && command.Tenant == "" {
		println("ERROR (CommandError): You must provide tenant_id via --tenant-id or  env[CC_TENANT_ID]")
		return
	}
	command.Sdk.Domain = command.Domain
	command.Sdk.Tenant = command.Tenant
	command.Sdk.Verbose = command.Verbose
	command.Sdk.Admin = command.Admin
	command.Sdk.Timeout = command.Timeout
	command.Sdk.Setup()
	cmd.Execute()
}
