package property

import (
	"time"
)

const (
	BaseAssetId        = "1.3.0"
	BaseAssetIdNum     = 0
	BaseAssetPrecision = 100000
)

func RFC3339ToUTC(timeFormatStr string) (uint64, error) {
	t, err := time.Parse(
		time.RFC3339, timeFormatStr+"+00:00")
	if err != nil {
		return 0, err
	}
	return uint64(t.Unix()), nil
}

func UTCToRFC3339(t uint64) string {
	timeStr := time.Unix(int64(t), 0).UTC().String()
	timeStr = timeStr[0:10] + "T" + timeStr[11:19]
	return timeStr
}
