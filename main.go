package main

import (
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/logger"
)

func main() {
	appInit()
}

func appInit() {
	// init log
	err := logger.InitAppLog("mint-tool.log")
	if err != nil {
		panic(err)
	}
}
