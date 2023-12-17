package tx_builder

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

const (
	// ExpireSeconds tx expire seconds
	ExpireSeconds = 3600
)

const (
	TxOpTypeTransfer    = 0
	TxOpTypeAccountBind = 10
	TxOpTypeMint        = 17
)

func BuildTxTransfer(refBlockNum uint16, refBlockPrefix uint32,
	fromAddr string, toAddr string, amount uint64, fee uint64) (string, []byte, *Transaction, error) {

	var tx Transaction
	tx.RefBlockNum = refBlockNum
	tx.RefBlockPrefix = refBlockPrefix

	tx.Expiration = UTCTime(time.Now().Unix() + ExpireSeconds)
	tx.Extensions = make([]interface{}, 0)
	tx.Signatures = make([]Signature, 0)

	var op TransferOperation
	err := op.SetValue(fromAddr, toAddr, amount, fee)
	if err != nil {
		return "", nil, nil, err
	}
	var opPair OperationPair
	opPair[0] = byte(TxOpTypeTransfer)
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
	opPair[0] = byte(TxOpTypeMint)
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
	opPair[0] = byte(TxOpTypeAccountBind)
	opPair[1] = &op
	tx.Operations = append(tx.Operations, opPair)

	txPacked := tx.Pack()

	s256 := sha256.New()
	_, err = s256.Write(txPacked)
	txHash := s256.Sum(nil)

	return hex.EncodeToString(txHash[0:20]), txPacked, &tx, nil
}
