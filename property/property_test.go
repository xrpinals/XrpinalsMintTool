package property

import (
	"fmt"
	"github.com/hedzr/assert"
	"testing"
)

func TestRFC3339ToUTC(t *testing.T) {
	timeStr := "2019-09-21T07:44:30"
	ts, err := RFC3339ToUTC(timeStr)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, uint64(1569051870), ts)
	fmt.Println("RFC3339ToUTC:", ts)
}

func TestUTCToRFC3339(t *testing.T) {
	timeStr := UTCToRFC3339(uint64(1569051870))
	assert.Equal(t, "2019-09-21T07:44:30", timeStr)
	fmt.Println("UTCToRFC3339:", timeStr)
}
