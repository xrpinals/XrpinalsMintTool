package tx_builder

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/xrpinals/XrpinalsMintTool/property"
	"time"
)

const (
	// ExpireSeconds tx expire seconds
	ExpireSeconds = 43200
)

func BuildTxTransfer(refBlockNum uint16, refBlockPrefix uint32,
	fromAddr string, toAddr, assetName string, amount uint64, fee uint64) (string, []byte, *Transaction, error) {

	var tx Transaction
	tx.RefBlockNum = refBlockNum
	tx.RefBlockPrefix = refBlockPrefix

	tx.Expiration = UTCTime(time.Now().Unix() + ExpireSeconds)
	tx.Extensions = make([]interface{}, 0)
	tx.Signatures = make([]Signature, 0)

	var op TransferOperation
	err := op.SetValue(fromAddr, toAddr, assetName, amount, fee)
	if err != nil {
		return "", nil, nil, err
	}
	var opPair OperationPair
	opPair[0] = byte(property.TxOpTypeTransfer)
	opPair[1] = &op
	tx.Operations = append(tx.Operations, opPair)

	txPacked := tx.Pack()

	s256 := sha256.New()
	_, err = s256.Write(txPacked)
	txHash := s256.Sum(nil)

	return hex.EncodeToString(txHash[0:20]), txPacked, &tx, nil
}

func BuildTxMint(refBlockNum uint16, refBlockPrefix uint32,
	issueAddr string,
	issueAssetId string, issueAssetIdNum int64, issueAmount int64, fee uint64) (string, []byte, *Transaction, error) {

	var tx Transaction
	tx.RefBlockNum = refBlockNum
	tx.RefBlockPrefix = refBlockPrefix

	tx.Expiration = UTCTime(time.Now().Unix() + ExpireSeconds)
	tx.Extensions = make([]interface{}, 0)
	tx.Signatures = make([]Signature, 0)

	var op MintOperation
	err := op.SetValue(issueAddr, issueAssetId, issueAssetIdNum, issueAmount, fee)
	if err != nil {
		return "", nil, nil, err
	}
	var opPair OperationPair
	opPair[0] = byte(property.TxOpTypeMint)
	opPair[1] = &op
	tx.Operations = append(tx.Operations, opPair)

	txPacked := tx.Pack()

	s256 := sha256.New()
	_, err = s256.Write(txPacked)
	txHash := s256.Sum(nil)

	return hex.EncodeToString(txHash[0:20]), txPacked, &tx, nil
}

func BuildTxAccountBind(refBlockNum uint16, refBlockPrefix uint32,
	keyWif string, fee uint64) (string, []byte, *Transaction, error) {

	var tx Transaction
	tx.RefBlockNum = refBlockNum
	tx.RefBlockPrefix = refBlockPrefix

	tx.Expiration = UTCTime(time.Now().Unix() + ExpireSeconds)
	tx.Extensions = make([]interface{}, 0)
	tx.Signatures = make([]Signature, 0)

	var op AccountBindOperation
	err := op.SetValue(keyWif, fee)
	if err != nil {
		return "", nil, nil, err
	}
	var opPair OperationPair
	opPair[0] = byte(property.TxOpTypeAccountBind)
	opPair[1] = &op
	tx.Operations = append(tx.Operations, opPair)

	txPacked := tx.Pack()

	s256 := sha256.New()
	_, err = s256.Write(txPacked)
	txHash := s256.Sum(nil)

	return hex.EncodeToString(txHash[0:20]), txPacked, &tx, nil
}

func BuildTxWithdraw(refBlockNum uint16, refBlockPrefix uint32,
	withdrawAccount string, amount string, toAddr string, memo string) (string, []byte, *Transaction, error) {

	var tx Transaction
	tx.RefBlockNum = refBlockNum
	tx.RefBlockPrefix = refBlockPrefix

	tx.Expiration = UTCTime(time.Now().Unix() + ExpireSeconds)
	tx.Extensions = make([]interface{}, 0)
	tx.Signatures = make([]Signature, 0)

	var op CrossChainWithdrawOperation
	var defaultAsset Asset
	defaultAsset.SetDefault()
	err := op.SetValue(withdrawAccount, amount, "BTC", defaultAsset.AssetId, toAddr, memo)
	if err != nil {
		return "", nil, nil, err
	}
	var opPair OperationPair
	opPair[0] = byte(property.TxOpTypeCrossChainWithdraw)
	opPair[1] = &op
	tx.Operations = append(tx.Operations, opPair)

	txPacked := tx.Pack()

	s256 := sha256.New()
	_, err = s256.Write(txPacked)
	txHash := s256.Sum(nil)

	return hex.EncodeToString(txHash[0:20]), txPacked, &tx, nil
}
