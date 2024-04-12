package main

import (
	"github.com/anoideaopen/testnet-cli/cmd"
	"github.com/anoideaopen/testnet-cli/logger"
)

var (
	version = "none"
	commit  = "none"
	date    = "none"
)

func main() {
	logger.Debug("Start Application: cli")

	if err := cmd.Execute(version, commit, date); err != nil {
		cmd.FatalError("Execute", err)
	}
}
