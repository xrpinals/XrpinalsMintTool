package tx_builder

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
)

const (
	MainNetPrefix = 0x0
	TestNetPrefix = 0x6f
)

const (
	UsePrefix = TestNetPrefix
)

func AddrToHexAddr(addr string) (string, error) {
	addrBytes, err := base58.Decode(addr)
	if err != nil {
		return "", err
	}
	if len(addrBytes) != AddressLength+PrefixLength+CheckSumLength {
		return "", fmt.Errorf("invalid wif address: length")
	}
	if addrBytes[0] != UsePrefix {
		return "", fmt.Errorf("invalid wif address: prefix")
	}

	s256 := sha256.New()
	_, _ = s256.Write(addrBytes[0 : AddressLength+PrefixLength])
	checkSum := s256.Sum(nil)

	s256 = sha256.New()
	_, _ = s256.Write(checkSum)
	checkSum = s256.Sum(nil)

	if bytes.Compare(addrBytes[AddressLength+PrefixLength:], checkSum[0:CheckSumLength]) != 0 {
		return "", fmt.Errorf("invalid wif format address: invalid checksum")
	}

	return hex.EncodeToString(addrBytes[PrefixLength : AddressLength+PrefixLength]), nil
}

func HexAddrToAddr(hexAddr string) (string, error) {
	hexAddrBytes, err := hex.DecodeString(hexAddr)
	if err != nil {
		return "", err
	}
	if len(hexAddrBytes) != AddressLength {
		return "", fmt.Errorf("invalid hex address: length")
	}
	calcBytes := make([]byte, 0)
	calcBytes = append(calcBytes, UsePrefix)
	calcBytes = append(calcBytes, hexAddrBytes...)

	s256 := sha256.New()
	_, _ = s256.Write(calcBytes)
	checkSum := s256.Sum(nil)

	s256 = sha256.New()
	_, _ = s256.Write(checkSum)
	checkSum = s256.Sum(nil)

	calcBytes = append(calcBytes, checkSum[0:CheckSumLength]...)

	return base58.Encode(calcBytes), nil
}

func PubKeyToAddr(pubKey string) (string, error) {
	pubKeyBytes, err := base58.Decode(pubKey)
	if err != nil {
		return "", err
	}
	if len(pubKeyBytes) != PubKeyLength+PrefixLength+CheckSumLength {
		return "", fmt.Errorf("invalid wif pubkey: length")
	}

	s256 := sha256.New()
	_, err = s256.Write(pubKeyBytes[PrefixLength : PubKeyLength+PrefixLength])
	pubHash := s256.Sum(nil)

	r160 := ripemd160.New()
	_, err = r160.Write(pubHash)
	pubHash = r160.Sum(nil)

	return HexAddrToAddr(hex.EncodeToString(pubHash))
}
