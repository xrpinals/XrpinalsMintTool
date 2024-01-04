package conf

import (
	"encoding/json"
	"fmt"
	"github.com/xrpinals/XrpinalsMintTool/utils"
	"os"
	"sync"
)

var (
	config     Config
	configOnce sync.Once
)

type Config struct {
	WalletRpcUrl string       `json:"walletRpcUrl"`
	Logs         LogsConfig   `json:"logs"`
	Data         DataConfig   `json:"data"`
	Server       ServerConfig `json:"server"`
	Gpu          bool         `json:"gpu"`
}

type LogsConfig struct {
	LogPath     string `json:"logPath"`
	Level       string `json:"level"`
	MaxSize     int64  `json:"maxSize"`
	BackupCount int64  `json:"backupCount"`
}

type DataConfig struct {
	DataPath string `json:"dataPath"`
}
type ServerConfig struct {
	IP   string `json:"IP"`
	Port string `json:"Port"`
}

func GetConfig() *Config {
	configOnce.Do(func() {
		bytes, err := os.ReadFile("conf.json")
		if err != nil {
			fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
			return
		}

		err = json.Unmarshal(bytes, &config)
		if err != nil {
			fmt.Println(utils.BoldRed("[Error]: "), utils.FgWhiteBgRed(err.Error()))
			return
		}
	})
	return &config
}
