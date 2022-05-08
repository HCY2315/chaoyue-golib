package limit_test

import (
	"git.cestong.com.cn/cecf/cecf-golib/pkg/limit"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestBuildFixedMinuteDurationResolve(t *testing.T) {
	minutes := int64(5)
	Convey("duration resolve", t, func() {
		resolve := limit.BuildFixedMinuteDurationResolve(minutes)
		Convey("same id", func() {
			t1s := "2022-02-07 13:55:51"
			t2s := "2022-02-07 13:58:51"
			t1, err := parseTime(t1s)
			So(err, ShouldBeNil)
			t2, err := parseTime(t2s)
			So(err, ShouldBeNil)
			So(t1, ShouldNotEqual, t2)
			id := "id"
			right1, durationID1 := resolve(id, t1)
			right2, durationID2 := resolve(id, t2)
			So(right1, ShouldEqual, right2)
			So(durationID1, ShouldEqual, durationID2)
			So(right1.Sub(t1), ShouldBeLessThanOrEqualTo, minutes*int64(time.Minute))
			Println(durationID1, right1)
		})
		Convey("not same id", func() {
			t1s := "2022-02-07 13:55:51"
			t2s := "2022-02-07 13:48:51"
			t1, err := parseTime(t1s)
			So(err, ShouldBeNil)
			t2, err := parseTime(t2s)
			So(err, ShouldBeNil)
			So(t1, ShouldNotEqual, t2)
			id := "id"
			right1, durationID1 := resolve(id, t1)
			right2, durationID2 := resolve(id, t2)
			So(right1, ShouldNotEqual, right2)
			So(durationID1, ShouldNotEqual, durationID2)
			So(right1.Sub(t1), ShouldBeLessThanOrEqualTo, minutes*int64(time.Minute))
			So(right2.Sub(t2), ShouldBeLessThanOrEqualTo, minutes*int64(time.Minute))
			Println(durationID1, durationID2)
			Println(right1, right2)
		})
	})
}

func parseTime(s string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", s)
}
