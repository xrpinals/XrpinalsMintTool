package tx_builder

import (
	//"fmt"
	"crypto/sha256"
	"encoding/hex"
	secp256k1 "github.com/bitnexty/secp256k1-go"
)

func SignTx(chainIdHex string, tx *Transaction, keyWif string) ([]byte, *Transaction, error) {
	chainIdBytes, err := hex.DecodeString(chainIdHex)
	if err != nil {
		return nil, nil, err
	}

	txData := tx.Pack()

	//fmt.Println("unsigned tx:", hex.EncodeToString(txData))
	//fmt.Println("chain id:", chainIdHex)

	s256 := sha256.New()
	_, _ = s256.Write(chainIdBytes)
	_, _ = s256.Write(txData)
	digestData := s256.Sum(nil)

	//fmt.Println("digest:", hex.EncodeToString(digestData))

	keyHex, _ := WifKeyToHexKey(keyWif)
	keyBytes, _ := hex.DecodeString(keyHex)

	var txSig []byte
	for {
		txSig, err = secp256k1.BtsSign(digestData, keyBytes, true)
		if err != nil {
			return nil, nil, err
		}

		if txSig[1] < 128 && txSig[33] < 128 {
			break
		}
	}

	// tx data with sig
	txBytesWithSig := make([]byte, 0)
	txBytesWithSig = append(txBytesWithSig, txData...)

	// sig count
	txBytesWithSig = append(txBytesWithSig, PackVarInt(1)...)

	txBytesWithSig = append(txBytesWithSig, PackVarInt(uint64(len(txSig)))...)
	txBytesWithSig = append(txBytesWithSig, txSig...)

	tx.Signatures = append(tx.Signatures, txSig)

	return txSig, tx, nil
}
