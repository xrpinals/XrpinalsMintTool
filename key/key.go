package key

import (
	"fmt"
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/conf"
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/tx_builder"
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/utils"
	"github.com/syndtr/goleveldb/leveldb"
)

func ImportPrivateKey(key string) (string, error) {
	addr, err := tx_builder.WifKeyToAddr(key)
	if err != nil {
		return "", err
	}

	result, err := utils.GetBindingAccount(conf.GetConfig().WalletRpcUrl, addr, "BTC")
	if err != nil {
		return "", err
	}

	if len(result) == 0 {
		// bind address
		refBlockNum, refBlockPrefix, err := utils.GetRefBlockInfo(conf.GetConfig().WalletRpcUrl)
		if err != nil {
			return "", err
		}

		chainId, err := utils.GetChainId(conf.GetConfig().WalletRpcUrl)
		if err != nil {
			return "", err
		}

		fee := uint64(0)

		_, _, tx, err := tx_builder.BuildTxAccountBind(refBlockNum, refBlockPrefix, key, fee)
		if err != nil {
			return "", err
		}

		_, txSigned, err := tx_builder.SignTx(chainId, tx, key)
		if err != nil {
			return "", err
		}

		_, err = utils.BroadcastTx(conf.GetConfig().WalletRpcUrl, txSigned)
		if err != nil {
			return "", err
		}
	}

	// import to leveldb
	keyDbDir := fmt.Sprintf("%s/keystore", conf.GetConfig().Data.DataPath)
	db, err := leveldb.OpenFile(keyDbDir, nil)
	if err != nil {
		return "", err
	}
	defer db.Close()

	err = db.Put([]byte(addr), []byte(key), nil)
	if err != nil {
		return "", err
	}

	return addr, nil
}

func IsAddressExisted(addr string) (bool, error) {
	// read from leveldb
	keyDbDir := fmt.Sprintf("%s/keystore", conf.GetConfig().Data.DataPath)
	db, err := leveldb.OpenFile(keyDbDir, nil)
	if err != nil {
		return false, err
	}
	defer db.Close()

	ok, err := db.Has([]byte(addr), nil)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func GetAddressKey(addr string) (string, error) {
	// read from leveldb
	keyDbDir := fmt.Sprintf("%s/keystore", conf.GetConfig().Data.DataPath)
	db, err := leveldb.OpenFile(keyDbDir, nil)
	if err != nil {
		return "", err
	}
	defer db.Close()

	key, err := db.Get([]byte(addr), nil)
	if err != nil {
		return "", err
	}

	return string(key), nil
}
