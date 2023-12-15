package tx_builder

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/property"
)

func PackUint8(v uint8) []byte {
	return []byte{v}
}

func PackUint16(v uint16) []byte {
	res := make([]byte, 2)
	binary.LittleEndian.PutUint16(res, v)
	return res
}

func PackUint32(v uint32) []byte {
	res := make([]byte, 4)
	binary.LittleEndian.PutUint32(res, v)
	return res
}

func PackUint64(v uint64) []byte {
	res := make([]byte, 8)
	binary.LittleEndian.PutUint64(res, v)
	return res
}

func PackVarInt(v uint64) []byte {
	res := make([]byte, 0)
	var dest uint8 = 0
	for {
		if v < 0x80 {
			break
		} else {
			dest = uint8((v & 0x7f) | 0x80)
			v = v >> 7
		}
		res = append(res, dest)
	}
	res = append(res, uint8(v))
	return res
}

type OperationType interface {
	Pack() []byte
}

type Asset struct {
	Amount         int64  `json:"amount"`
	AssetId        string `json:"asset_id"`
	AssetIdNum     int64  `json:"-"`
	AssetPrecision int64  `json:"-"`
}

func (a *Asset) SetDefault() {
	a.Amount = 0
	a.AssetId = property.BaseAssetId
	a.AssetIdNum = property.BaseAssetIdNum
	a.AssetPrecision = property.BaseAssetPrecision
}

func (a Asset) Pack() []byte {
	bytesRet := make([]byte, 0)
	bytesAmount := PackUint64(uint64(a.Amount))
	bytesAssetId := PackUint8(uint8(a.AssetIdNum))
	bytesRet = append(bytesRet, bytesAmount...)
	bytesRet = append(bytesRet, bytesAssetId...)
	return bytesRet
}

type TransferOperation struct {
	Fee         Asset         `json:"fee"`
	GuaranteeId string        `json:"guarantee_id,omitempty"`
	From        string        `json:"from"`
	To          string        `json:"to"`
	FromAddr    Address       `json:"from_addr"`
	ToAddr      Address       `json:"to_addr"`
	Amount      Asset         `json:"amount"`
	Memo        *interface{}  `json:"memo,omitempty"`
	Extensions  []interface{} `json:"extensions"`
}

func (to *TransferOperation) SetValue(fromAddr string, toAddr string,
	amount uint64, fee uint64) error {

	to.Fee.SetDefault()
	to.Fee.Amount = int64(fee)

	to.Amount.SetDefault()
	to.Amount.Amount = int64(amount)

	to.From = "1.2.0"
	to.To = "1.2.0"

	to.Extensions = make([]interface{}, 0)

	fromAddrHex, err := AddrToHexAddr(fromAddr)
	if err != nil {
		return err
	}
	fromAddrBytes, _ := hex.DecodeString(fromAddrHex)
	to.FromAddr.SetBytes(fromAddrBytes)

	toAddrHex, err := AddrToHexAddr(toAddr)
	if err != nil {
		return err
	}
	toAddrBytes, _ := hex.DecodeString(toAddrHex)
	to.ToAddr.SetBytes(toAddrBytes)

	to.Memo = nil

	return nil
}

func (to *TransferOperation) Pack() []byte {
	bytesRet := make([]byte, 0)

	bytesFee := to.Fee.Pack()
	bytesAmount := to.Amount.Pack()

	bytesRet = append(bytesRet, bytesFee...)
	//guarantee_id
	bytesRet = append(bytesRet, byte(0))
	//from
	bytesRet = append(bytesRet, byte(0))
	//to
	bytesRet = append(bytesRet, byte(0))

	bytesRet = append(bytesRet, byte(UseAddressPrefix))
	bytesRet = append(bytesRet, to.FromAddr[:]...)
	bytesRet = append(bytesRet, byte(UseAddressPrefix))
	bytesRet = append(bytesRet, to.ToAddr[:]...)
	bytesRet = append(bytesRet, bytesAmount...)

	// pack empty memo
	bytesRet = append(bytesRet, byte(0))

	// Extensions
	bytesRet = append(bytesRet, PackVarInt(uint64(len(to.Extensions)))...)

	return bytesRet
}

type MintOperation struct {
	Fee            Asset         `json:"fee"`
	Issuer         string        `json:"issuer"`
	AssetToIssue   Asset         `json:"asset_to_issue"`
	IssueToAccount string        `json:"issue_to_account"`
	IssueAddress   Address       `json:"issue_address"`
	BosToken       bool          `json:"bos_token"`
	Memo           *interface{}  `json:"memo,omitempty"`
	Extensions     []interface{} `json:"extensions"`
}

func (to *MintOperation) SetValue(issueAddr string,
	issueAssetId string, issueAssetIdNum int64, issueAmount int64, fee uint64) error {

	to.Fee.SetDefault()
	to.Fee.Amount = int64(fee)

	to.Issuer = "1.2.0"

	to.AssetToIssue.Amount = issueAmount
	to.AssetToIssue.AssetId = issueAssetId
	to.AssetToIssue.AssetIdNum = issueAssetIdNum

	to.IssueToAccount = "1.2.0"

	to.BosToken = true

	to.Extensions = make([]interface{}, 0)

	issueAddrHex, err := AddrToHexAddr(issueAddr)
	if err != nil {
		return err
	}
	issueAddrBytes, _ := hex.DecodeString(issueAddrHex)
	to.IssueAddress.SetBytes(issueAddrBytes)

	to.Memo = nil

	return nil
}

func (to *MintOperation) Pack() []byte {
	bytesRet := make([]byte, 0)

	bytesFee := to.Fee.Pack()
	bytesAssetToIssue := to.AssetToIssue.Pack()

	bytesRet = append(bytesRet, bytesFee...)
	//issuer
	bytesRet = append(bytesRet, byte(0))

	bytesRet = append(bytesRet, bytesAssetToIssue...)
	//issue_to_account
	bytesRet = append(bytesRet, byte(0))

	bytesRet = append(bytesRet, byte(UseAddressPrefix))
	bytesRet = append(bytesRet, to.IssueAddress[:]...)

	//bos_token
	bytesRet = append(bytesRet, byte(1))

	// pack empty memo
	bytesRet = append(bytesRet, byte(0))

	// Extensions
	bytesRet = append(bytesRet, PackVarInt(uint64(len(to.Extensions)))...)

	return bytesRet
}

type OperationPair [2]interface{}

type Transaction struct {
	RefBlockNum    uint16          `json:"ref_block_num"`
	RefBlockPrefix uint32          `json:"ref_block_prefix"`
	Expiration     UTCTime         `json:"expiration"`
	Operations     []OperationPair `json:"operations"`
	Extensions     []interface{}   `json:"extensions"`
	Signatures     []Signature     `json:"signatures"`
}

func (tx *Transaction) Pack() []byte {
	bytesRet := make([]byte, 0)

	bytesRefBlockNum := PackUint16(tx.RefBlockNum)
	bytesRefBlockPrefix := PackUint32(tx.RefBlockPrefix)
	bytesExpiration := PackUint32(uint32(tx.Expiration))

	bytesRet = append(bytesRet, bytesRefBlockNum...)
	bytesRet = append(bytesRet, bytesRefBlockPrefix...)
	bytesRet = append(bytesRet, bytesExpiration...)

	bytesRet = append(bytesRet, PackVarInt(uint64(len(tx.Operations)))...)
	for _, opPair := range tx.Operations {
		bytesRet = append(bytesRet, PackVarInt(uint64(opPair[0].(byte)))...)
		bytesOP := opPair[1].(OperationType).Pack()
		bytesRet = append(bytesRet, bytesOP...)
	}

	//extension
	bytesRet = append(bytesRet, byte(0))

	//without sig
	return bytesRet
}

func bytesToNumber(bs []byte) uint32 {
	if len(bs) != 4 {
		return 0
	}
	return (uint32(bs[0]) << 24) + (uint32(bs[1]) << 16) + (uint32(bs[2]) << 8) + uint32(bs[3])
}
