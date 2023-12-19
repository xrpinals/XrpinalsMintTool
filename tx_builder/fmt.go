package tx_builder

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
)

const (
	MainNetAddressPrefix = 0x0
	TestNetAddressPrefix = 0x6f

	MainNetPubKeyPrefix = 0x0
	TestNetPubKeyPrefix = 0x6f

	MainNetSecretPrefix = 0x80
	TestNetSecretPrefix = 0xef
)

const (
	//UseAddressPrefix = MainNetAddressPrefix
	UseAddressPrefix = TestNetAddressPrefix

	UsePubKeyPrefix = MainNetPubKeyPrefix
	UseSecretPrefix = MainNetSecretPrefix
)

func WifKeyToHexKey(wifKey string) (string, error) {
	privateKeyBytes, err := base58.Decode(wifKey)
	if err != nil {
		return "", err
	}

	if len(privateKeyBytes) != PrivateKeyLength+PrefixLength+CheckSumLength && len(privateKeyBytes) != PrivateKeyLength+PrefixLength+KeyCpsFlagLength+CheckSumLength {
		return "", fmt.Errorf("invalid wif key: length")
	}
	if privateKeyBytes[0] != UseSecretPrefix {
		return "", fmt.Errorf("invalid wif key: prefix")
	}

	return hex.EncodeToString(privateKeyBytes[PrefixLength : PrivateKeyLength+PrefixLength]), nil
}

func HexKeyToWifKey(hexKey string) (string, error) {
	hexKeyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return "", err
	}
	if len(hexKeyBytes) != PrivateKeyLength {
		return "", fmt.Errorf("invalid hex key: length")
	}
	calcBytes := make([]byte, 0)
	calcBytes = append(calcBytes, UseSecretPrefix)
	calcBytes = append(calcBytes, hexKeyBytes...)

	// double sha256
	s256 := sha256.New()
	_, err = s256.Write(calcBytes)
	if err != nil {
		return "", err
	}
	checkSum := s256.Sum(nil)

	s256 = sha256.New()
	_, err = s256.Write(checkSum)
	if err != nil {
		return "", err
	}
	checkSum = s256.Sum(nil)

	calcBytes = append(calcBytes, checkSum[0:CheckSumLength]...)

	return base58.Encode(calcBytes), nil
}

func AddrToHexAddr(addr string) (string, error) {
	addrBytes, err := base58.Decode(addr)
	if err != nil {
		return "", err
	}
	if len(addrBytes) != AddressLength+PrefixLength+CheckSumLength {
		return "", fmt.Errorf("invalid wif address: length")
	}
	if addrBytes[0] != UseAddressPrefix {
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
	calcBytes = append(calcBytes, UseAddressPrefix)
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
	if pubKeyBytes[0] != UsePubKeyPrefix {
		return "", fmt.Errorf("invalid wif pubkey: prefix")
	}

	s256 := sha256.New()
	_, err = s256.Write(pubKeyBytes[PrefixLength : PubKeyLength+PrefixLength])
	pubHash := s256.Sum(nil)

	r160 := ripemd160.New()
	_, err = r160.Write(pubHash)
	pubHash = r160.Sum(nil)

	return HexAddrToAddr(hex.EncodeToString(pubHash))
}

func GetCompressPubKey(pubKeyBytes []byte) ([]byte, error) {
	pubKeyBytes = pubKeyBytes[1:]
	if len(pubKeyBytes) != 2*(PubKeyLength-1) {
		return nil, fmt.Errorf("GetCompressPubKey: invalid pubKeyBytes size")
	}

	pubKeyCompressBytes := make([]byte, PubKeyLength)
	if pubKeyBytes[2*(PubKeyLength-1)-1]%2 == 0 {
		pubKeyCompressBytes[0] = 0x2
	} else {
		pubKeyCompressBytes[0] = 0x3
	}
	copy(pubKeyCompressBytes[1:], pubKeyBytes[0:PubKeyLength-1])

	return pubKeyCompressBytes, nil
}

func WifKeyToAddr(wifKey string) (string, error) {
	privateKeyBytes, err := base58.Decode(wifKey)
	if err != nil {
		return "", err
	}
	if len(privateKeyBytes) != PrivateKeyLength+PrefixLength+CheckSumLength && len(privateKeyBytes) != PrivateKeyLength+PrefixLength+KeyCpsFlagLength+CheckSumLength {
		return "", fmt.Errorf("invalid wif key: length")
	}
	if privateKeyBytes[0] != UseSecretPrefix {
		return "", fmt.Errorf("invalid wif key: prefix")
	}

	key, err := crypto.HexToECDSA(hex.EncodeToString(privateKeyBytes[PrefixLength : PrefixLength+PrivateKeyLength]))
	if err != nil {
		return "", err
	}
	pub, ok := key.Public().(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("invalid wif key: get public key")
	}
	compressed, err := GetCompressPubKey(crypto.FromECDSAPub(pub))

	if err != nil {
		return "", err
	}

	s256 := sha256.New()
	_, err = s256.Write(compressed)
	pubHash := s256.Sum(nil)

	r160 := ripemd160.New()
	_, err = r160.Write(pubHash)
	pubHash = r160.Sum(nil)

	return HexAddrToAddr(hex.EncodeToString(pubHash))
}
