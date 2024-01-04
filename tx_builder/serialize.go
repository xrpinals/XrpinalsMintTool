package tx_builder

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"github.com/bitnexty/secp256k1-go"
	"github.com/xrpinals/XrpinalsMintTool/property"
	"strconv"
	"strings"
)

func PackUint8(v uint8) []byte {
	return []byte{v}
}

func PackUint16(v uint16) []byte {
	res := make([]byte, 2)
	binary.LittleEndian.PutUint16(res, v)
	return res
}

func PackString(v string) []byte {
	res := make([]byte, 0)
	res = append(res, PackVarInt(uint64(len(v)))...)
	res = append(res, []byte(v)...)

	return res
}

func PackBytes(v []byte) []byte {
	res := make([]byte, 0)
	res = append(res, PackVarInt(uint64(len(v)))...)
	res = append(res, v...)

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

func (to *TransferOperation) SetValue(fromAddr string, toAddr string, assetName string,
	amount uint64, fee uint64) error {

	to.Fee.SetDefault()
	to.Fee.Amount = int64(fee)

	to.Amount.SetDefault()
	to.Amount.Amount = int64(amount)

	if len(assetName) > 0 {
		to.Amount.AssetId = assetName

		l := strings.Split(to.Amount.AssetId, ".")
		idNum, err := strconv.Atoi(l[len(l)-1])
		if err != nil {
			return err
		}
		to.Amount.AssetIdNum = int64(idNum)
	}

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
	Brc20Token     bool          `json:"brc20_token"`
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

	to.Brc20Token = true

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

	//brc20_token
	bytesRet = append(bytesRet, byte(1))

	// pack empty memo
	bytesRet = append(bytesRet, byte(0))

	// Extensions
	bytesRet = append(bytesRet, PackVarInt(uint64(len(to.Extensions)))...)

	return bytesRet
}

type AccountBindOperation struct {
	Fee              Asset     `json:"fee"`
	CrosschainType   string    `json:"crosschain_type"`
	Addr             Address   `json:"addr"`
	AccountSignature Signature `json:"account_signature"`
	TunnelAddress    string    `json:"tunnel_address"`
	TunnelSignature  string    `json:"tunnel_signature"`
	GuaranteeId      string    `json:"guarantee_id,omitempty"`
}

func (to *AccountBindOperation) SetValue(keyWif string, fee uint64) error {
	to.Fee.SetDefault()
	to.Fee.Amount = int64(fee)

	to.CrosschainType = "BTC"

	keyHex, err := WifKeyToHexKey(keyWif)
	if err != nil {
		return err
	}

	keyBytes, err := hex.DecodeString(keyHex)
	if err != nil {
		return err
	}

	tunnelAddr, err := WifKeyToAddr(keyWif)
	if err != nil {
		return err
	}

	to.TunnelAddress = tunnelAddr
	addrHex, err := AddrToHexAddr(tunnelAddr)
	if err != nil {
		return err
	}

	addrBytes, err := hex.DecodeString(addrHex)
	if err != nil {
		return err
	}

	to.Addr.SetBytes(addrBytes)

	var tunnelBytes []byte
	tunnelBytes = append(tunnelBytes, []byte("Bitcoin Signed Message:\n")...)
	tunnelBytes = append(tunnelBytes, PackString(tunnelAddr)...)

	// sign tunnel address
	s256 := sha256.New()
	_, _ = s256.Write(tunnelBytes)
	digestData := s256.Sum(nil)

	var tunnelSig []byte
	for {
		tunnelSig, err = secp256k1.BtsSign(digestData, keyBytes, true)
		if err != nil {
			return err
		}

		if tunnelSig[0] < 128 && tunnelSig[32] < 128 {
			break
		}
	}

	to.TunnelSignature = base64.StdEncoding.EncodeToString(tunnelSig)

	var accountBytes []byte
	accountBytes = append(accountBytes, byte(UseAddressPrefix))
	accountBytes = append(accountBytes, addrBytes...)

	s256Account := sha256.New()
	_, _ = s256.Write(accountBytes)
	digestAccountData := s256Account.Sum(nil)

	// sign account address
	var accountSig []byte
	for {
		accountSig, err = secp256k1.BtsSign(digestAccountData, keyBytes, true)
		if err != nil {
			return err
		}

		if accountSig[0] < 128 && accountSig[32] < 128 {
			break
		}
	}

	to.AccountSignature = accountSig

	return nil
}

func (to *AccountBindOperation) Pack() []byte {
	bytesRet := make([]byte, 0)

	bytesFee := to.Fee.Pack()
	bytesRet = append(bytesRet, bytesFee...)

	// crosschain_type
	bytesRet = append(bytesRet, PackString(to.CrosschainType)...)

	// addr
	bytesRet = append(bytesRet, byte(UseAddressPrefix))
	bytesRet = append(bytesRet, to.Addr[:]...)

	// account signature
	bytesRet = append(bytesRet, to.AccountSignature...)

	// tunnel address
	bytesRet = append(bytesRet, PackString(to.TunnelAddress)...)

	// tunnel signature
	bytesRet = append(bytesRet, PackString(to.TunnelSignature)...)

	//guarantee_id
	bytesRet = append(bytesRet, byte(0))

	return bytesRet
}

type OperationPair [2]interface{}

type Transaction struct {
	RefBlockNum    uint16          `json:"ref_block_num"`
	RefBlockPrefix uint32          `json:"ref_block_prefix"`
	Expiration     UTCTime         `json:"expiration"`
	Operations     []OperationPair `json:"operations"`
	Extensions     []interface{}   `json:"extensions"`
	NoncePow       uint64          `json:"nonce"`
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

	//pack nonce 0
	bytesNoncePow := PackUint64(tx.NoncePow)
	bytesRet = append(bytesRet, bytesNoncePow...)

	//without sig
	return bytesRet
}

func bytesToNumber(bs []byte) uint32 {
	if len(bs) != 4 {
		return 0
	}
	return (uint32(bs[0]) << 24) + (uint32(bs[1]) << 16) + (uint32(bs[2]) << 8) + uint32(bs[3])
}

type CrossChainWithdrawOperation struct {
	Fee               Asset   `json:"fee"`
	WithdrawAccount   Address `json:"withdraw_account"`
	Amount            string  `json:"amount"`
	AssetSymbol       string  `json:"asset_symbol"`
	AssetId           string  `json:"asset_id"`
	CrossChainAccount string  `json:"crosschain_account"`
	Memo              string  `json:"memo"`
}

func (to *CrossChainWithdrawOperation) SetValue(withdrawAddr string,
	amount string, assetSymbol string, assetId string, crossChainAccount string, memo string) error {
	to.Fee.SetDefault()

	withdrawAddrHex, err := AddrToHexAddr(withdrawAddr)
	if err != nil {
		return err
	}
	withdrawAddrBytes, _ := hex.DecodeString(withdrawAddrHex)
	to.WithdrawAccount.SetBytes(withdrawAddrBytes)

	to.Amount = amount
	to.AssetSymbol = assetSymbol
	to.AssetId = assetId
	to.CrossChainAccount = crossChainAccount
	to.Memo = memo
	return nil
}

func (to *CrossChainWithdrawOperation) Pack() []byte {
	bytesRet := make([]byte, 0)
	var assetType Asset
	assetType.SetDefault()
	bytesFee := to.Fee.Pack()
	bytesRet = append(bytesRet, bytesFee...)
	bytesRet = append(bytesRet, byte(UseAddressPrefix))
	bytesRet = append(bytesRet, to.WithdrawAccount[:]...)
	bytesRet = append(bytesRet, PackString(to.Amount)...)
	bytesRet = append(bytesRet, PackString(to.AssetSymbol)...)
	bytesRet = append(bytesRet, PackUint8(uint8(assetType.AssetIdNum))...)
	bytesRet = append(bytesRet, PackString(to.CrossChainAccount)...)
	bytesRet = append(bytesRet, PackString(to.Memo)...)

	return bytesRet
}
