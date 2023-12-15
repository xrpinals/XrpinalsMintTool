package logger

import (
	"github.com/Xrpinals-Protocol/XrpinalsMintTool/conf"
	"github.com/siddontang/go-log/log"
	"os"
	"strings"
)

var Logger *log.Logger

const (
	Mask = 0755
)

func InitAppLog(logFileName string) error {
	config := conf.GetConfig()

	fileInfo, err := os.Stat(config.Logs.LogPath)
	if err != nil {
		if os.IsNotExist(err) {
			// not exist
			err = os.Mkdir(config.Logs.LogPath, Mask)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		if !fileInfo.IsDir() {
			// exist, but is not dir
			err := os.RemoveAll(config.Logs.LogPath)
			if err != nil {
				return err
			}
			err = os.Mkdir(config.Logs.LogPath, Mask)
			if err != nil {
				return err
			}
		}
	}

	fileName := strings.Join([]string{config.Logs.LogPath, logFileName}, "/")
	fileHandle, err := log.NewRotatingFileHandler(fileName, int(config.Logs.MaxSize), int(config.Logs.BackupCount))
	if err != nil {
		return err
	}

	Logger = log.NewDefault(fileHandle)
	Logger.SetLevelByName(strings.ToLower(config.Logs.Level))

	return nil
}
