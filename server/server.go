package server

import (
	"github.com/xrpinals/XrpinalsMintTool/conf"
	"github.com/xrpinals/XrpinalsMintTool/logger"
	"net"
	"sync"
)

type XrpMintServer struct {
	minersLock  sync.RWMutex
	miners      map[*Miner]struct{}
	nonceReport chan uint64
	reBuildTrx  chan bool
	mineInfo    MineInfo
	maxMapCount int64
}

func (xrp *XrpMintServer) AddMiner(miner *Miner) {
	xrp.minersLock.Lock()
	defer xrp.minersLock.Unlock()
	if xrp.mineInfo.TrxId != "" {
		miner.UpdateMineInfo(xrp.mineInfo.TrxId, xrp.mineInfo.Difficult, false)
	}
	xrp.miners[miner] = struct{}{}
}
func (xrp *XrpMintServer) removeMiner(miner *Miner) {
	xrp.minersLock.Lock()
	defer xrp.minersLock.Unlock()
	if _, found := xrp.miners[miner]; found {
		delete(xrp.miners, miner)
	}
}
func (xrp *XrpMintServer) NotifyMiner(trxid string, diff uint32, nonceReport chan uint64, rebuildTrx chan bool) {
	xrp.nonceReport = nonceReport
	xrp.reBuildTrx = rebuildTrx
	xrp.mineInfo.TrxId = trxid
	xrp.mineInfo.Difficult = diff
	xrp.minersLock.RLock()
	defer xrp.minersLock.RUnlock()
	if len(xrp.miners) != 0 {
		for miner, _ := range xrp.miners {
			go func(miner *Miner) {
				miner.UpdateMineInfo(trxid, diff, true)
			}(miner)
		}

	}
}

func (xrp *XrpMintServer) ListenTcp() {
	config := conf.GetConfig()
	xrp.miners = make(map[*Miner]struct{}, 1)
	addr, err := net.ResolveTCPAddr("tcp", config.Server.IP+":"+config.Server.Port)
	if err != nil {
		logger.Logger.Fatalf("ResolveTCPAddr error: %v", err)
	}
	server, err := net.ListenTCP("tcp", addr)
	if err != nil {
		logger.Logger.Fatalf("ListenTCP error: %v", err)
	}
	defer server.Close()
	var accept = make(chan int, 1)
	n := 0
	ret := uint64(0)
	resultNonce := make(chan uint64)
	quit := make(chan struct{})
	go func() {
		ret = <-resultNonce
		close(quit)
	}()
	for {
		select {
		case <-quit:
			goto END
		default:
			conn, err := server.AcceptTCP()
			if err != nil {
				continue
			}
			conn.SetKeepAlive(true)
			ip, port, _ := net.SplitHostPort(conn.RemoteAddr().String())
			accept <- n
			miner := &Miner{conn: conn, IP: ip, Port: port, resultNonce: resultNonce}
			xrp.AddMiner(miner)
			go func(miner *Miner) {
				err := miner.handleTCPClient()
				xrp.removeMiner(miner)
				logger.Logger.Warn(err)
				conn.Close()
				<-accept
			}(miner)
		}

	}
END:
	xrp.nonceReport <- ret
}
