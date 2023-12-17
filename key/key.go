package key

import (
	"fmt"
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/conf"
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/tx_builder"
	"github.com/syndtr/goleveldb/leveldb"
)

func ImportPrivateKey(key string) (string, error) {
	// bind address
	// TODO

	// import to leveldb
	keyDbDir := fmt.Sprintf("%s/keystore", conf.GetConfig().Data.DataPath)
	db, err := leveldb.OpenFile(keyDbDir, nil)
	if err != nil {
		return "", err
	}
	defer db.Close()

	addr, err := tx_builder.WifKeyToAddr(key)
	if err != nil {
		return "", err
	}

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
