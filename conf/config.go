package conf

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

var (
	config     Config
	configOnce sync.Once
)

type Config struct {
	WalletRpcUrl string     `json:"walletRpcUrl"`
	Logs         LogsConfig `json:"logs"`
}

type LogsConfig struct {
	LogPath     string `json:"logPath"`
	Level       string `json:"level"`
	MaxSize     int64  `json:"maxSize"`
	BackupCount int64  `json:"backupCount"`
}

func GetConfig() *Config {
	configOnce.Do(func() {
		bytes, err := ioutil.ReadFile("conf.json")
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal(bytes, &config)
		if err != nil {
			panic(err)
		}
	})
	return &config
}
