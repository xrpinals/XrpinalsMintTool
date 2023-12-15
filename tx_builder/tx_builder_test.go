package tx_builder

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/utils"
	"strconv"
	"strings"
	"testing"
)

var (
	walletUrl = "http://192.168.1.73:8093"
)

func TestBuildTxTransfer(t *testing.T) {
	refBlockNum, refBlockPrefix, err := utils.GetRefBlockInfo(walletUrl)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(refBlockNum)
	fmt.Println(refBlockPrefix)

	fromAddr := "mfhGJnP5T7A5kYDJNxnHozxrVzC7WKHzKs"
	toAddr := "mfhGJnP5T7A5kYDJNxnHozxrVzC7WKHzKs"
	amount := uint64(100000)
	fee := uint64(100000)

	txBytes, tx, _ := BuildTxTransfer(refBlockNum, refBlockPrefix, fromAddr, toAddr, amount, fee)
	fmt.Println("BuildTxTransfer Hex:", hex.EncodeToString(txBytes))
	txJson, _ := json.Marshal(*tx)
	fmt.Println("BuildTxTransfer Tx:", string(txJson))
}

func TestBuildTxMint(t *testing.T) {
	refBlockNum, refBlockPrefix, err := utils.GetRefBlockInfo(walletUrl)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(refBlockNum)
	fmt.Println(refBlockPrefix)

	resp, err := utils.GetAssetInfo(walletUrl, "BTC")
	if err != nil {
		t.Fatal(err.Error())
	}

	issueAddr := "mfhGJnP5T7A5kYDJNxnHozxrVzC7WKHzKs"
	issueAssetId := resp.Result.Id
	l := strings.Split(resp.Result.Id, ".")
	issueAssetIdNum, err := strconv.Atoi(l[len(l)-2])
	if err != nil {
		t.Fatal(err.Error())
	}
	issueAmount := resp.Result.Options.MaxPerMint
	fee := uint64(100000)

	txBytes, tx, _ := BuildTxMint(refBlockNum, refBlockPrefix, issueAddr, issueAssetId, int64(issueAssetIdNum), int64(issueAmount), fee)
	fmt.Println("BuildTxMint Hex:", hex.EncodeToString(txBytes))
	txJson, _ := json.Marshal(*tx)
	fmt.Println("BuildTxMint Tx:", string(txJson))
}
