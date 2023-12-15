package tx_builder

import (
	"fmt"
	"github.com/hedzr/assert"
	"testing"
)

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
	wifKey := "5JuANbr9xBwz6ASN2ifEpcMkFPAUThBJGcVWj3TLpt6PhAXbpva"
	addr, err := WifKeyToAddr(wifKey)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("WifKeyToAddr:", addr)
	assert.Equal(t, "mnQhe9c4dxq8nWkno41GEDxGL6N4w5AezR", addr)
}
