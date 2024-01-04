package tx_builder

import (
	"github.com/xrpinals/XrpinalsMintTool/conf"
	"github.com/xrpinals/XrpinalsMintTool/utils"
	"strconv"
)

func Transfer(from, to, asset, amount, keyWif string) (txHash string, err error) {
	amountUint, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		return "", err
	}

	refBlockNum, refBlockPrefix, err := utils.GetRefBlockInfo(conf.GetConfig().WalletRpcUrl)
	if err != nil {
		return "", err
	}

	chainId, err := utils.GetChainId(conf.GetConfig().WalletRpcUrl)
	if err != nil {
		return "", err
	}

	fee := uint64(100)

	_, _, tx, err := BuildTxTransfer(refBlockNum, refBlockPrefix, from, to, asset, amountUint, fee)
	if err != nil {
		return "", err
	}

	_, txSigned, err := SignTx(chainId, tx, keyWif)
	if err != nil {
		return "", err
	}

	txHash, err = utils.BroadcastTx(conf.GetConfig().WalletRpcUrl, txSigned)
	if err != nil {
		return "", err
	}

	return txHash, nil
}

func Withdraw(from, to, amount, memo, keyWif string) (txHash string, err error) {

	refBlockNum, refBlockPrefix, err := utils.GetRefBlockInfo(conf.GetConfig().WalletRpcUrl)
	if err != nil {
		return "", err
	}

	chainId, err := utils.GetChainId(conf.GetConfig().WalletRpcUrl)
	if err != nil {
		return "", err
	}

	_, _, tx, err := BuildTxWithdraw(refBlockNum, refBlockPrefix, from, amount, to, memo)
	if err != nil {
		return "", err
	}

	_, txSigned, err := SignTx(chainId, tx, keyWif)
	if err != nil {
		return "", err
	}

	txHash, err = utils.BroadcastTx(conf.GetConfig().WalletRpcUrl, txSigned)
	if err != nil {
		return "", err
	}

	return txHash, nil
}
