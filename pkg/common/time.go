package common

import (
	"database/sql/driver"
	"encoding/json"
	"git.cestong.com.cn/cecf/cecf-golib/pkg/utils"
	"time"
)

type CETime time.Time

func (c *CETime) Value() (driver.Value, error) {
	return time.Time(*c), nil
}

func LocalNow() CETime {
	t := time.Now()
	return CETime(t)
}

func (c *CETime) UnmarshalJSON(dt []byte) error {
	t, err := time.ParseInLocation(CETimeLayout, utils.Bytes2String(dt), CETimeLoc)
	if err != nil {
		return err
	}
	*c = CETime(t)
	return nil
}

func LocalFromStr(str string) (CETime, error) {
	t, err := time.ParseInLocation(CETimeLayout, str, time.Local)
	if err != nil {
		return CETime{}, err
	}
	return CETime(t), nil
}

var CETimeLayout = "2006-01-02 15:04:05"
var CETimeLoc = time.Local

func (c CETime) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.AsResp())
}

func CETimeFromTime(t time.Time) CETime {
	return CETime(t)
}

func CETimeForNow() CETime {
	return CETimeFromTime(NowForLocal())
}

func (c CETime) AsResp() string {
	return time.Time(c).Format(time.RFC3339)
}

func (c *CETime) Before(t2 time.Time) bool {
	return time.Time(*c).Before(t2)
}

func NowMS() int64 {
	return time.Now().UnixMilli()
}

func NowUnixSec() int64 {
	return time.Now().Unix()
}

func NowForLocal() time.Time {
	return time.Now()
}
