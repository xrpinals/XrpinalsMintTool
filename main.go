package main

import (
	"flag"
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/logger"
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/mining"
	"os"
)

func appInit() {
	// init log
	err := logger.InitAppLog("mint-tool.log")
	if err != nil {
		panic(err)
	}
}

func main() {
	appInit()

	if os.Args[1] == "mint" {
		fs := flag.NewFlagSet("mint", flag.ExitOnError)

		var pKey string
		var asset string
		fs.StringVar(&pKey, "key", "", "your private key")
		fs.StringVar(&asset, "asset", "", "asset name you want to mint")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}

		mining.PrivateKey = pKey
		mining.MintAssetName = asset

		for i := 0; i < 4; i++ {
			mining.StartMining()
		}
	} else if os.Args[1] == "import_key" {
		fs := flag.NewFlagSet("import_key", flag.ExitOnError)

		var pKey string
		fs.StringVar(&pKey, "key", "", "your private key")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}

		mining.PrivateKey = pKey

	}
}
