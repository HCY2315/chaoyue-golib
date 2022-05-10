package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

const (
	JsonTimeFormat = "2006-01-02 15:04:05"
)

type JsonTime time.Time

func (t JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", time.Time(t).Format(JsonTimeFormat))), nil
}

func (t *JsonTime) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 {
		return errors.New("data is empty")
	}
	jt, err := time.Parse(`"`+JsonTimeFormat+`"`, string(data))
	*t = JsonTime(jt)
	return err
}

// before current time for true
func (t JsonTime) Before(u JsonTime) bool {
	return time.Time(t).Before(time.Time(u))
}

func (t JsonTime) IsZero() bool {
	return time.Time(t).IsZero()
}

type UnixTime time.Time

func (t UnixTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", time.Time(t).Unix())), nil
}

func (t *UnixTime) UnmarshalJSON(data []byte) (err error) {
	var unixSec int64
	if data[0] == '"' {
		var s string
		if err = json.Unmarshal(data, &s); err != nil {
			return
		}
		unixSec, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
	} else {
		unixSec, err = strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return err
		}
	}
	*t = UnixTime(time.Unix(unixSec, 0))
	return nil
}

// 数据库(xorm)存毫秒精度的Local时间，定义时不要用指针，指针xorm处理有问题
type DbMsLocalTime time.Time

const dbMsTimeReadFormat = "2006-01-02T15:04:05.000+08:00"
const dbMsTimeWriteFormat = "2006-01-02 15:04:05.000"

func (t DbMsLocalTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(t).Format(dbMsTimeWriteFormat) + `"`), nil
}

func (t *DbMsLocalTime) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 {
		return errors.New("data is empty")
	}
	jt, err := time.Parse(`"`+dbMsTimeWriteFormat+`"`, string(data))
	*t = DbMsLocalTime(jt)
	return err
}

func (t DbMsLocalTime) String() string {
	return time.Time(t).Format(dbMsTimeWriteFormat)
}

func (t *DbMsLocalTime) FromDB(b []byte) error {
	now, err := time.ParseInLocation(dbMsTimeReadFormat, string(b), time.Local)
	if nil == err {
		*t = DbMsLocalTime(now)
		return nil
	}
	return err
}

func (t *DbMsLocalTime) ToDB() ([]byte, error) {
	return []byte(time.Time(*t).Format(dbMsTimeWriteFormat)), nil
}
