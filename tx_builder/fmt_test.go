package tx_builder

import (
	"fmt"
	"github.com/hedzr/assert"
	"testing"
)

func TestWifKeyToHexKey(t *testing.T) {
	wifKey := "5KcnSNrBJEdGAcmjVzzThtpncNtuZDDf74Fj81sEvYYkij7bs6u"
	hexKey, err := WifKeyToHexKey(wifKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("TestWifKeyToHexKey:", hexKey)
	assert.Equal(t, "ed4640fd09578c07bc180298b8ea4f454d0daa2fff791b5f4b3e9ae42b0e4af5", hexKey)
}

func TestHexKeyToWifKey(t *testing.T) {
	hexKey := "ed4640fd09578c07bc180298b8ea4f454d0daa2fff791b5f4b3e9ae42b0e4af5"
	wifKey, err := HexKeyToWifKey(hexKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("TestHexKeyToWifKey:", wifKey)
	assert.Equal(t, "5KcnSNrBJEdGAcmjVzzThtpncNtuZDDf74Fj81sEvYYkij7bs6u", wifKey)
}

func TestAddrToHexAddr(t *testing.T) {
	addr := "mumtmaYKH3ttGpVaAJRCiiWZsn5zAB9hU9"
	hexAddr, err := AddrToHexAddr(addr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("TestAddrToHexAddr:", hexAddr)
	assert.Equal(t, "9c650eb99282993f42151c9862e2d1509ac83781", hexAddr)
}

func TestHexAddrToAddr(t *testing.T) {
	hexAddr := "9c650eb99282993f42151c9862e2d1509ac83781"
	addr, err := HexAddrToAddr(hexAddr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("TestHexAddrToAddr:", addr)
	assert.Equal(t, "mumtmaYKH3ttGpVaAJRCiiWZsn5zAB9hU9", addr)
}

func TestPubKeyToAddr(t *testing.T) {
	pubKey := "17vz1qN7yu6eDA7hZKpGyHweJs12XhtWiJRhn8wuwRpkfC17JCH"
	addr, err := PubKeyToAddr(pubKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("PubKeyToAddr:", addr)
	assert.Equal(t, "mumtmaYKH3ttGpVaAJRCiiWZsn5zAB9hU9", addr)
}

func TestWifKeyToAddr(t *testing.T) {
	wifKey := "5JF7asAXBFzGbnLDdLyKqrkRGGKcSJByU22fvzejdU6TdLGimdf"
	addr, err := WifKeyToAddr(wifKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("WifKeyToAddr:", addr)
	assert.Equal(t, "mfhGJnP5T7A5kYDJNxnHozxrVzC7WKHzKs", addr)
}
