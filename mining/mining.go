package mining

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/bitcoin"
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/conf"
	. "github.com/Xrpinals-Protocol/XrpinalsMintTool/logger"
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/tx_builder"
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/utils"
	"math/big"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	MinerNum  = 1
	isStop    atomic.Bool
	config    *conf.Config
	Difficult uint32
)

type Miner struct{}

func init() {
	numCPU := runtime.NumCPU()
	fmt.Println("Number of CPUs: ", numCPU)
	runtime.GOMAXPROCS(numCPU)
	if numCPU > 1 {
		MinerNum = numCPU - 1
	}
	isStop.Store(false)
	config = conf.GetConfig()
}

func StartMining() {
	var wg sync.WaitGroup

	resp, err := utils.GetAssetInfo(config.WalletRpcUrl, config.AssetName)
	if err != nil {
		panic(err)
	}
	Difficult = resp.Result.DynamicData.CurrentNBits

	for i := 0; i < MinerNum; i++ {
		wg.Add(1)
		miner := Miner{}
		go miner.mining(&wg, uint64(i))
	}
	wg.Wait()
}

func (m *Miner) buildTx() (string, *tx_builder.Transaction, error) {
	// build mint tx
	txHash, tx, err := m.getMintTxHash()
	if err != nil {
		Logger.Errorf("buildTx: getMintTxHash err: %v", err)
		return "", nil, err
	}
	return txHash, tx, nil
}

func (m *Miner) signTx(tx *tx_builder.Transaction) (*tx_builder.Transaction, error) {
	chainId, err := utils.GetChainId(config.WalletRpcUrl)
	if err != nil {
		Logger.Errorf("buildTx: GetChainId err: %v", err)
		return nil, err
	}

	_, txSigned, err := tx_builder.SignTx(chainId, tx, config.PrivateKey)
	if err != nil {
		Logger.Errorf("signTx: SignTx err: %v", err)
		return nil, err
	}

	return txSigned, nil
}

func (m *Miner) mining(wg *sync.WaitGroup, nonce uint64) {
	defer wg.Done()

	txHash, unSignedTx, err := m.buildTx()
	if err != nil {
		Logger.Errorf("mining: buildTx err: %v", err)
		return
	}

	target := bitcoin.NBits2Target(Difficult)

	for {
		payload := PowPayload{
			Version:  1,
			TxHash:   txHash,
			Reserved: [44]byte{},
			NBits:    Difficult,
			Nonce:    nonce,
		}

		payloadBytes, err := payload.pack()
		if err != nil {
			Logger.Errorf("mining: payload.pack err: %v", err)
			return
		}

		s256 := sha256.New()
		_, err = s256.Write(payloadBytes)
		if err != nil {
			Logger.Errorf("mining: s256.Write err: %v", err)
			return
		}
		hashBytes := s256.Sum(nil)
		utils.ReverseBytesInPlace(hashBytes)
		result := new(big.Int).SetBytes(hashBytes)

		if result.Cmp(target) < 0 {
			if isStop.CompareAndSwap(false, true) {
				fmt.Println("payload: ", hex.EncodeToString(payloadBytes))
				fmt.Println("hash: ", hex.EncodeToString(hashBytes))
				break
			} else {
				return
			}
		}

		// next nonce
		nonce = nonce + uint64(MinerNum)
	}

	// broadcast tx
	unSignedTx.NoncePow = nonce
	signedTx, err := m.signTx(unSignedTx)
	if err != nil {
		Logger.Errorf("mining: signTx err: %v", err)
		return
	}

	_, err = utils.BroadcastTx(config.WalletRpcUrl, signedTx)
	if err != nil {
		Logger.Errorf("mining: utils.BroadcastTx err: %v", err)
		return
	}

	Logger.Infof("Mining success, txHash:%v", txHash)
}

func (m *Miner) getMintTxHash() (string, *tx_builder.Transaction, error) {
	refBlockNum, refBlockPrefix, err := utils.GetRefBlockInfo(config.WalletRpcUrl)
	if err != nil {
		return "", nil, err
	}

	resp, err := utils.GetAssetInfo(config.WalletRpcUrl, config.AssetName)
	if err != nil {
		return "", nil, err
	}

	issueAddr, err := tx_builder.WifKeyToAddr(config.PrivateKey)
	if err != nil {
		return "", nil, err
	}

	issueAssetId := resp.Result.Id
	l := strings.Split(resp.Result.Id, ".")
	issueAssetIdNum, err := strconv.Atoi(l[len(l)-1])
	if err != nil {
		return "", nil, err
	}
	issueAmount, err := utils.Uint64Supply(resp.Result.Options.MaxPerMint)
	if err != nil {
		return "", nil, err
	}
	fee := uint64(100000)

	txHashCalc, _, tx, err := tx_builder.BuildTxMint(refBlockNum, refBlockPrefix, issueAddr, issueAssetId, int64(issueAssetIdNum), int64(issueAmount), fee)
	if err != nil {
		return "", nil, err
	}

	return txHashCalc, tx, nil
}

type PowPayload struct {
	Version  uint32
	TxHash   string
	Reserved [44]byte
	NBits    uint32
	Nonce    uint64
}

func (p *PowPayload) pack() ([]byte, error) {
	bytesRet := make([]byte, 0)

	bytesRet = append(bytesRet, tx_builder.PackUint32(p.Version)...)

	hashBytes, err := hex.DecodeString(p.TxHash)
	if err != nil {
		return nil, err
	}
	bytesRet = append(bytesRet, hashBytes...)
	bytesRet = append(bytesRet, p.Reserved[:]...)
	bytesRet = append(bytesRet, tx_builder.PackUint32(p.NBits)...)
	bytesRet = append(bytesRet, tx_builder.PackUint64(p.Nonce)...)

	return bytesRet, nil
}
