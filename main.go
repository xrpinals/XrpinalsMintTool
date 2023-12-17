package main

import (
	"flag"
	"fmt"
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/key"
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

		var addr string
		var asset string
		fs.StringVar(&addr, "addr", "", "your address")
		fs.StringVar(&asset, "asset", "", "asset name you want to mint")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}
		mining.MintAssetName = asset

		ok, err := key.IsAddressExisted(addr)
		if err != nil {
			panic(err)
		}

		if ok {
			mining.PrivateKey, err = key.GetAddressKey(addr)
			if err != nil {
				panic(err)
			}
		} else {
			panic(fmt.Errorf("address %s is not in the storage", addr))
		}

		mining.StartMining()

	} else if os.Args[1] == "import_key" {
		fs := flag.NewFlagSet("import_key", flag.ExitOnError)

		var pKey string
		fs.StringVar(&pKey, "key", "", "your private key")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}

		addr, err := key.ImportPrivateKey(pKey)
		if err != nil {
			panic(err)
		}

		fmt.Printf("private key of address %s imported", addr)
	} else if os.Args[1] == "check_address" {
		fs := flag.NewFlagSet("check_address", flag.ExitOnError)

		var addr string
		fs.StringVar(&addr, "addr", "", "your address")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}

		ok, err := key.IsAddressExisted(addr)
		if err != nil {
			panic(err)
		}

		if ok {
			fmt.Printf("address %s is in the storage", addr)
		} else {
			panic(fmt.Errorf("address %s is not in the storage", addr))
		}
	}
}
