package server

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xrpinals/XrpinalsMintTool/logger"
	"github.com/xrpinals/XrpinalsMintTool/tx_builder"
	"github.com/xrpinals/XrpinalsMintTool/utils"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
)

const (
	NotCalcNonce  = 0
	CalcNonceDone = 1
	CalcNonceEnd  = 2
)

type MinerRequest struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}
type Miner struct {
	sync.Mutex
	IP            string
	Port          string
	conn          *net.TCPConn
	enc           *json.Encoder
	mineinfoMutex sync.Mutex
	mineInfo      MineInfo
	resultNonce   chan uint64
}
type MineInfo struct {
	TrxId     string
	Difficult uint32
}

func (miner *Miner) UpdateMineInfo(trxid string, diff uint32, send bool) {
	miner.mineinfoMutex.Lock()
	defer miner.mineinfoMutex.Unlock()
	miner.mineInfo.TrxId = trxid
	miner.mineInfo.Difficult = diff
	if send {
		respUpdate := fmt.Sprintf("{\"method\":\"mining.notify\",\"params\":[\"%s\",\"%s\"]}\n", miner.mineInfo.TrxId, strconv.FormatUint(uint64(miner.mineInfo.Difficult), 16))
		miner.SendMessageToMiner(respUpdate)
	}
}

func (miner *Miner) handleTCPMessage(req *MinerRequest) (int, error) {
	var err error
	calcState := NotCalcNonce
	switch req.Method {
	case "mining_subscribe", "mining.subscribe", "getwork":
		miner.getJob()
	case "mining_submit", "mining.submit":
		if req.Params[1].(string) == miner.mineInfo.TrxId {
			err = miner.submitHandle(req)
			if err == nil {
				calcState = CalcNonceDone
			}
		}
	default:
		miner.SendMessageToMiner("{\"type\":\"solo\"}\n")
	}
	return calcState, err
}
func (miner *Miner) getJob() {

	if miner.mineInfo.TrxId == "" {
		miner.SendMessageToMiner("{\"method\":\"erro\"}\n")
		return
	}
	msg := fmt.Sprintf("{\"method\":\"mining.notify\",\"params\":[\"%s\",\"%s\"]}\n", miner.mineInfo.TrxId, strconv.FormatUint(uint64(miner.mineInfo.Difficult), 16))
	fmt.Println(utils.BoldYellow("[Info]: "), utils.Bold("GPU miner getJob: [\""+miner.mineInfo.TrxId+"\",\""+strconv.FormatUint(uint64(miner.mineInfo.Difficult), 16)+"\"]\n"))
	miner.SendMessageToMiner(msg)
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
func Reverse(b []byte) []byte {
	length := len(b)
	dest := make([]byte, length)
	for i := range dest {
		dest[i] = b[length-1-i]
	}
	return dest
}
func (miner *Miner) submitHandle(r *MinerRequest) error {
	strNonce := r.Params[0].(string)
	bytes, _ := hex.DecodeString(strNonce)
	reversedBytes := Reverse(bytes)
	result := hex.EncodeToString(reversedBytes)
	realNonce, err := strconv.ParseUint(result, 16, 64)
	if err != nil {
		logger.Logger.Warn(err.Error())
		return err
	}
	/*
		payload := PowPayload{
			Version:  1,
			TxHash:   miner.mineInfo.TrxId,
			Reserved: [44]byte{},
			NBits:    miner.mineInfo.Difficult,
			Nonce:    realNonce,
		}
		payloadBytes, _ := payload.pack()
		fmt.Println("payLoad is ", hex.EncodeToString(payloadBytes))
		hashBytes := utils.X17_Byte_Sum256(payloadBytes)
		hashBytes = Reverse(hashBytes)
		calcHard := new(big.Int).SetUint64(0)
		calcHard.SetBytes(hashBytes)
		target := bitcoin.NBits2Target(miner.mineInfo.Difficult)
		if calcHard.Cmp(target) >= 0 {
			return errors.New("Not Target")
		}*/
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	go func() {
		miner.resultNonce <- realNonce
	}()
	return nil
}
func (miner *Miner) handleTCPClient() error {
	miner.enc = json.NewEncoder(miner.conn)
	connbuff := bufio.NewReaderSize(miner.conn, 1024)
	go func() {
		for _ = range time.Tick(1 * time.Minute) {
			ping := fmt.Sprintf("{\"method\":\"mining.ping\",\"params\":[\"ping\"]}\n")
			miner.SendMessageToMiner(ping)
		}
	}()
	for {
		line, isPrefix, err := connbuff.ReadLine()
		if isPrefix {
			logger.Logger.Warn("Socket flood detected from %s", miner.IP)
			return err

		} else if err == io.EOF {
			return err

		} else if err != nil {
			return err
		}
		var req MinerRequest
		err = json.Unmarshal(line, &req)
		if err != nil {
			return err
		}
		calcStatue, err := miner.handleTCPMessage(&req)
		if calcStatue == CalcNonceDone {
			return errors.New("Find Nonce")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (miner *Miner) SendMessageToMiner(msg string) error {
	_, err := miner.conn.Write([]byte(msg))
	logger.Logger.Info("SendMessageToMiner:", msg)
	if err != nil {
		logger.Logger.Info("SendMessageToMiner failed")
	}
	return err
}
