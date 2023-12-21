package mining

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
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
	"time"
)

var (
	MinerNum  = 1
	isStop    atomic.Bool
	Difficult uint32
)

var (
	PrivateKey    = ""
	MintAssetName = ""
)

type Miner struct{}

func init() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	if numCPU > 1 {
		MinerNum = numCPU - 1
	}
}

func preCheck(assetInfo *utils.AssetInfoRsp) error {
	addr, err := tx_builder.WifKeyToAddr(PrivateKey)
	if err != nil {
		return err
	}

	maxMintCountLimit, err := utils.Uint64Supply(assetInfo.Result.Options.MaxMintCountLimit)
	if err != nil {
		return err
	}

	maxSupply, err := utils.Uint64Supply(assetInfo.Result.Options.MaxSupply)
	if err != nil {
		return err
	}

	currentSupply, err := utils.Uint64Supply(assetInfo.Result.DynamicData.CurrentSupply)
	if err != nil {
		return err
	}

	mintInfo, err := utils.GetAddressMintInfo(conf.GetConfig().WalletRpcUrl, addr, MintAssetName)
	if err != nil {
		return err
	}

	lastMintTime, err := utils.DataTimeToTimestamp(mintInfo.Result.Time)
	if err != nil {
		return err
	}

	if time.Now().Unix()-lastMintTime < assetInfo.Result.Options.MintInterval {
		return errors.New("less than the mint interval")
	}

	if mintInfo.Result.MintCount >= maxMintCountLimit {
		return errors.New("address had mint max count")
	}

	if mintInfo.Result.Amount+currentSupply > maxSupply {
		return errors.New("beyond max mint amount")
	}

	return nil
}

func StartMining() {
	isStop.Store(false)
	resp, err := utils.GetAssetInfo(conf.GetConfig().WalletRpcUrl, MintAssetName)
	if err != nil {
		fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
		return
	}
	if !resp.Result.Options.Brc20Token {
		fmt.Println(utils.FgWhiteBgRed("not brc20 token, can not mint"))
		return
	}

	err = preCheck(resp)
	if err != nil {
		fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
		return
	}

	Difficult = resp.Result.DynamicData.CurrentNBits

	var wg sync.WaitGroup
	for i := 0; i < MinerNum; i++ {
		wg.Add(1)
		miner := Miner{}
		go miner.mining(&wg, uint64(i))
	}
	wg.Wait()

}

func (m *Miner) buildMintTx() (string, *tx_builder.Transaction, error) {
	// build mint tx
	txHash, tx, err := m.getMintTx()
	if err != nil {
		Logger.Errorf("buildMintTx: getMintTx err: %v", err)
		return "", nil, err
	}
	return txHash, tx, nil
}

func (m *Miner) signMintTx(tx *tx_builder.Transaction) (*tx_builder.Transaction, error) {
	chainId, err := utils.GetChainId(conf.GetConfig().WalletRpcUrl)
	if err != nil {
		Logger.Errorf("signMintTx: GetChainId err: %v", err)
		return nil, err
	}

	_, txSigned, err := tx_builder.SignTx(chainId, tx, PrivateKey)
	if err != nil {
		Logger.Errorf("signMintTx: SignTx err: %v", err)
		return nil, err
	}

	return txSigned, nil
}

func (m *Miner) mining(wg *sync.WaitGroup, nonce uint64) {
	defer wg.Done()

	statHash := false
	origNonce := nonce
	if nonce == 0 {
		statHash = true
	}

ReBuildTx:
	nonce = origNonce
	statIdx := int64(0)
	txHash, unSignedTx, err := m.buildMintTx()
	if err != nil {
		Logger.Errorf("mining: buildMintTx err: %v", err)
		return
	}

	target := bitcoin.NBits2Target(Difficult)
	statStart := time.Now().UnixMicro()
	txBuildTime := time.Now().Unix()

	for {
		if statHash {
			if statIdx > 10000000 {
				hashRate := float64(statIdx) / float64(time.Now().UnixMicro()-statStart)
				hashRateStr := fmt.Sprintf("%.03f", hashRate)
				fmt.Println(utils.BoldYellow("[Mining]: "),
					utils.Bold("Pow Hash Speed ------------------------- "),
					utils.FgWhiteBgGreen(hashRateStr), utils.Bold("MHash/s"))

				statIdx = 0
				statStart = time.Now().UnixMicro()
			}
			statIdx = statIdx + int64(MinerNum)
		}

		if (time.Now().Unix() - txBuildTime) > tx_builder.ExpireSeconds/2 {
			fmt.Println(utils.BoldYellow("[Mining]: "),
				utils.BoldGreen("Re-build mint transaction in case of expiring"))
			goto ReBuildTx
		}

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
		result := new(big.Int).SetBytes(hashBytes)

		if result.Cmp(target) < 0 {
			if isStop.CompareAndSwap(false, true) {
				break
			} else {
				return
			}
		} else {
			if isStop.Load() {
				return
			}
		}

		// next nonce
		nonce = nonce + uint64(MinerNum)
	}

	// broadcast tx
	unSignedTx.NoncePow = nonce
	signedTx, err := m.signMintTx(unSignedTx)
	if err != nil {
		Logger.Errorf("mining: signMintTx err: %v", err)
		return
	}

	_, err = utils.BroadcastTx(conf.GetConfig().WalletRpcUrl, signedTx)
	if err != nil {
		fmt.Printf("mining failed: err: %v\n", err)
		Logger.Errorf("mining: utils.BroadcastTx err: %v", err)
		return
	}

	fmt.Println(utils.BoldYellow("[Info]: "), utils.Bold("mining success, txHash: "), utils.FgWhiteBgBlue(txHash))
	Logger.Infof("mining success, txHash:%v", txHash)
}

func (m *Miner) getMintTx() (string, *tx_builder.Transaction, error) {
	refBlockNum, refBlockPrefix, err := utils.GetRefBlockInfo(conf.GetConfig().WalletRpcUrl)
	if err != nil {
		return "", nil, err
	}

	resp, err := utils.GetAssetInfo(conf.GetConfig().WalletRpcUrl, MintAssetName)
	if err != nil {
		return "", nil, err
	}

	if !resp.Result.Options.Brc20Token {
		return "", nil, fmt.Errorf("not brc20 token, can not mint")
	}

	issueAddr, err := tx_builder.WifKeyToAddr(PrivateKey)
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
	fee := uint64(100)

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
