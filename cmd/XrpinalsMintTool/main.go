package main

import (
	"flag"
	"fmt"
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/conf"
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/key"
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/logger"
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/mining"
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/tx_builder"
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/utils"
	"github.com/fatih/color"
	"github.com/shopspring/decimal"
	"math"
	"os"
)

func appInit() {
	err := logger.InitAppLog("XrpinalsMintTool.log")
	if err != nil {
		fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
		return
	}
}

func main() {
	appInit()
	fmt.Println(color.GreenString(utils.AsciiPicure))

	if os.Args[1] == "mint" {
		fs := flag.NewFlagSet("mint", flag.ExitOnError)

		var addr string
		var asset string
		fs.StringVar(&addr, "addr", "", "your address")
		fs.StringVar(&asset, "asset", "", "asset name you want to mint")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
			return
		}
		mining.MintAssetName = asset

		ok, err := key.IsAddressExisted(addr)
		if err != nil {
			fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
			return
		}

		if ok {
			mining.PrivateKey, err = key.GetAddressKey(addr)
			if err != nil {
				fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
				return
			}
		} else {
			fmt.Println(utils.BoldRed("[Error]: "), utils.Bold("Address"), utils.FgWhiteBgBlue(addr), utils.Bold("is"),
				utils.BoldRed("NOT"), utils.Bold("the Storage!"))
			return
		}

		mining.StartMining()

	} else if os.Args[1] == "import_key" {
		fs := flag.NewFlagSet("import_key", flag.ExitOnError)

		var pKey string
		fs.StringVar(&pKey, "key", "", "your private key")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
			return
		}

		addr, err := key.ImportPrivateKey(pKey)
		if err != nil {
			fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
			return
		}

		fmt.Println(utils.BoldYellow("[Info]: "), utils.Bold("Private Key of Address"),
			utils.FgWhiteBgBlue(addr), utils.Bold("imported."))

	} else if os.Args[1] == "check_address" {
		fs := flag.NewFlagSet("check_address", flag.ExitOnError)

		var addr string
		fs.StringVar(&addr, "addr", "", "your address")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
			return
		}

		ok, err := key.IsAddressExisted(addr)
		if err != nil {
			fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
			return
		}

		if ok {
			fmt.Println(utils.BoldGreen("[Result]: "), utils.Bold("Address"),
				utils.FgWhiteBgBlue(addr), utils.Bold("is in the Storage."))
			return
		} else {
			fmt.Println(utils.BoldGreen("[Result]: "), utils.Bold("Address"), utils.FgWhiteBgBlue(addr),
				utils.Bold("is"), utils.BoldRed("NOT"), utils.Bold("the Storage."))
			return
		}

	} else if os.Args[1] == "get_balance" {
		fs := flag.NewFlagSet("get_balance", flag.ExitOnError)

		var addr string
		var asset string
		fs.StringVar(&addr, "addr", "", "your address")
		fs.StringVar(&asset, "asset", "", "asset name you want to query")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
			return
		}

		resp, err := utils.GetAssetInfo(conf.GetConfig().WalletRpcUrl, asset)
		if err != nil {
			fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
			return
		}

		balance, err := utils.GetAddressBalance(conf.GetConfig().WalletRpcUrl, addr, resp.Result.Id)
		if err != nil {
			fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
			return
		}

		balanceDecimal := decimal.NewFromBigInt(balance, 0)
		precisionDecimal := decimal.NewFromFloat(math.Pow(10, float64(resp.Result.Precision)))
		balanceDecimal = balanceDecimal.Div(precisionDecimal)

		fmt.Println(utils.BoldGreen("[Result]: "), utils.Bold("Balance: "),
			utils.Bold(balanceDecimal.String()), utils.BoldYellow(asset))

	} else if os.Args[1] == "get_deposit_address" {
		fs := flag.NewFlagSet("get_deposit_address", flag.ExitOnError)

		var addr string
		fs.StringVar(&addr, "addr", "", "your address")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
			return
		}

		ok, err := key.IsAddressExisted(addr)
		if err != nil {
			fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
			return
		}

		if !ok {
			fmt.Println(utils.BoldRed("[Error]: "), utils.Bold("Address"), utils.FgWhiteBgBlue(addr),
				utils.Bold("is"), utils.BoldRed("NOT"), utils.Bold("the Storage, you MUST call import_key first!"))
			return
		}

		result, err := utils.GetDepositAddress(conf.GetConfig().WalletRpcUrl, "BTC")
		if err != nil {
			fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
			return
		}

		fmt.Println(utils.BoldGreen("[Result]: "), utils.Bold("BTC deposit address: "), utils.FgWhiteBgBlue(result.BindAccountHot))

	} else if os.Args[1] == "transfer" {
		fs := flag.NewFlagSet("transfer", flag.ExitOnError)

		var from string
		var to string
		var asset string
		var amount string
		var keyWif string

		fs.StringVar(&from, "from", "", "your address")
		fs.StringVar(&to, "to", "", "receiver address")
		fs.StringVar(&asset, "asset", "", "asset name you want to transfer")
		fs.StringVar(&amount, "amount", "", "asset amount you want to transfer")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
			return
		}

		ok, err := key.IsAddressExisted(from)
		if err != nil {
			fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
			return
		}

		if ok {
			keyWif, err = key.GetAddressKey(from)
			if err != nil {
				fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
				return
			}
		} else {
			fmt.Println(utils.BoldRed("[Error]: "), utils.Bold("Address"), utils.FgWhiteBgBlue(from),
				utils.Bold("is"), utils.BoldRed("NOT"), utils.Bold("the Storage!"))
			return
		}

		resp, err := utils.GetAssetInfo(conf.GetConfig().WalletRpcUrl, asset)
		if err != nil {
			fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
			return
		}

		amountDecimal, err := decimal.NewFromString(amount)
		if err != nil {
			fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
			return
		}

		precisionDecimal := decimal.NewFromFloat(math.Pow(10, float64(resp.Result.Precision)))
		amountDecimal = amountDecimal.Mul(precisionDecimal)

		txHash, err := tx_builder.Transfer(from, to, resp.Result.Id, amountDecimal.StringFixed(0), keyWif)
		if err != nil {
			fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
			return
		}

		fmt.Println(utils.BoldYellow("[Info]: "), utils.Bold("Transfer Success, txHash:"), utils.FgWhiteBgBlue(txHash))

	} else if os.Args[1] == "get_mint_info" {
		fs := flag.NewFlagSet("get_mint_info", flag.ExitOnError)

		var addr string
		var asset string
		fs.StringVar(&addr, "addr", "", "your address")
		fs.StringVar(&asset, "asset", "", "asset name you want to query")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
			return
		}

		mintInfo, err := utils.GetAddressMintInfo(conf.GetConfig().WalletRpcUrl, addr, asset)
		if err != nil {
			fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
			return
		}

		resp, err := utils.GetAssetInfo(conf.GetConfig().WalletRpcUrl, asset)
		if err != nil {
			fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
			return
		}

		amountDecimal := decimal.NewFromInt(int64(mintInfo.Result.Amount))
		precisionDecimal := decimal.NewFromFloat(math.Pow(10, float64(resp.Result.Precision)))
		amountDecimal = amountDecimal.Div(precisionDecimal)

		fmt.Println(utils.BoldGreen("[Result]: "), utils.Bold("Mint Info:"))
		fmt.Println(utils.BoldGreen("[Result]: "), utils.Bold("Mint amount total:"), utils.Bold(amountDecimal.String()))
		fmt.Println(utils.BoldGreen("[Result]: "), utils.Bold("Mint count:"), utils.Bold(mintInfo.Result.MintCount))
		fmt.Println(utils.BoldGreen("[Result]: "), utils.Bold("Last mint time:"), utils.Bold(mintInfo.Result.Time))

	}
}
