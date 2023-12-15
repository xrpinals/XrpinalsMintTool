package tx_builder

import (
	"time"
)

const (
	// ExpireSeconds tx expire seconds
	ExpireSeconds = 600
)

const (
	TxOpTypeTransfer = 0
)

func BuildTxTransfer(refBlockNum uint16, refBlockPrefix uint32,
	fromAddr string, toAddr string, amount uint64, fee uint64) ([]byte, *Transaction, error) {

	var tx Transaction
	tx.RefBlockNum = refBlockNum
	tx.RefBlockPrefix = refBlockPrefix
	//expire 10 min
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
