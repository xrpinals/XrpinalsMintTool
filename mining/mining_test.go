package mining

import (
	"crypto/sha256"
	"fmt"
	"github.com/xrpinals/XrpinalsMintTool/bitcoin"
	"math/big"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestMining(t *testing.T) {
	var nBits uint32 = 0x1e00ffff
	var wg sync.WaitGroup

	var isStop atomic.Bool
	isStop.Store(false)

	var numCPU = runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	if numCPU > 1 {
		MinerNum = numCPU - 1
	}

	start := time.Now()

	for i := 0; i < MinerNum; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, nBits uint32, i int, step int, isStop *atomic.Bool) {
			defer wg.Done()

			var nonce uint64 = 0
			target := bitcoin.NBits2Target(nBits)

			for {
				payloadBytes := make([]byte, 80)
				rand.Read(payloadBytes)

				s256 := sha256.New()
				_, _ = s256.Write(payloadBytes)
				hashBytes := s256.Sum(nil)
				result := new(big.Int).SetBytes(hashBytes)

				if result.Cmp(target) < 0 {
					if isStop.CompareAndSwap(false, true) {
						return
					} else {
						return
					}
				} else {
					if isStop.Load() {
						return
					}
				}

				// next nonce
				nonce = nonce + uint64(step)
			}

		}(&wg, nBits, i, MinerNum, &isStop)
	}
	wg.Wait()

	duration := time.Since(start)

	fmt.Println(duration)
}
