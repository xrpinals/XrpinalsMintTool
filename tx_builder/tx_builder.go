package tx_builder

import (
	"time"
)

const (
	// ExpireSeconds tx expire seconds
	ExpireSeconds = 3600
)

const (
	TxOpTypeTransfer = 0
	TxOpTypeMint     = 17
)

func BuildTxTransfer(refBlockNum uint16, refBlockPrefix uint32,
	fromAddr string, toAddr string, amount uint64, fee uint64) ([]byte, *Transaction, error) {

	var tx Transaction
	tx.RefBlockNum = refBlockNum
	tx.RefBlockPrefix = refBlockPrefix

	tx.Expiration = UTCTime(time.Now().Unix() + ExpireSeconds)
	tx.Extensions = make([]interface{}, 0)
	tx.Signatures = make([]Signature, 0)

	var op TransferOperation
	err := op.SetValue(fromAddr, toAddr, amount, fee)
	if err != nil {
		return nil, nil, err
	}
	var opPair OperationPair
	opPair[0] = byte(TxOpTypeTransfer)
	opPair[1] = &op
	tx.Operations = append(tx.Operations, opPair)

	return tx.Pack(), &tx, nil
}

func BuildTxMint(refBlockNum uint16, refBlockPrefix uint32,
	issueAddr string,
	issueAssetId string, issueAssetIdNum int64, issueAmount int64, fee uint64) ([]byte, *Transaction, error) {

	var tx Transaction
	tx.RefBlockNum = refBlockNum
	tx.RefBlockPrefix = refBlockPrefix

	tx.Expiration = UTCTime(time.Now().Unix() + ExpireSeconds)
	tx.Extensions = make([]interface{}, 0)
	tx.Signatures = make([]Signature, 0)

	var op MintOperation
	err := op.SetValue(issueAddr, issueAssetId, issueAssetIdNum, issueAmount, fee)
	if err != nil {
		return nil, nil, err
	}
	var opPair OperationPair
	opPair[0] = byte(TxOpTypeMint)
	opPair[1] = &op
	tx.Operations = append(tx.Operations, opPair)

	return tx.Pack(), &tx, nil
}
