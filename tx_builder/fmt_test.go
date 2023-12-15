package tx_builder

import (
	"fmt"
	"github.com/hedzr/assert"
	"testing"
)

func TestAddrToHexAddr(t *testing.T) {
	addr := "mumtmaYKH3ttGpVaAJRCiiWZsn5zAB9hU9"
	hexAddr, _ := AddrToHexAddr(addr)
	fmt.Println("TestAddrToHexAddr:", hexAddr)
	assert.Equal(t, "9c650eb99282993f42151c9862e2d1509ac83781", hexAddr)
}

func TestHexAddrToAddr(t *testing.T) {
	hexAddr := "9c650eb99282993f42151c9862e2d1509ac83781"
	addr, _ := HexAddrToAddr(hexAddr)
	fmt.Println("TestHexAddrToAddr:", addr)
	assert.Equal(t, "mumtmaYKH3ttGpVaAJRCiiWZsn5zAB9hU9", addr)
}

func TestPubKeyToAddr(t *testing.T) {
	pubKey := "17vz1qN7yu6eDA7hZKpGyHweJs12XhtWiJRhn8wuwRpkfC17JCH"
	addr, _ := PubKeyToAddr(pubKey)
	fmt.Println("PubKeyToAddr:", addr)
	assert.Equal(t, "mumtmaYKH3ttGpVaAJRCiiWZsn5zAB9hU9", addr)
}
