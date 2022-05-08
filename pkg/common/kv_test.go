package common

import (
	errors2 "git.cestong.com.cn/cecf/cecf-golib/pkg/errors"
	"git.cestong.com.cn/cecf/cecf-golib/pkg/utils"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestMemNamedValuesGetHook(t *testing.T) {
	Convey("override OK", t, func() {
		k, v := "key", "value"
		ov := "valueOverride"
		hk := func(values NamedValues, key string) interface{} {
			if k == key {
				return ov
			}
			return nil
		}
		mkv := NewSimpleMemNamedValuesWithHook(hk)
		So(mkv.Set(k, v), ShouldBeNil)
		valueGotObj, errGet := mkv.Get(k)
		So(errGet, ShouldBeNil)
		valueGot, ok := valueGotObj.(string)
		So(ok, ShouldBeTrue)
		So(valueGot, ShouldEqual, ov)
	})
	Convey("named values impl exp time", t, func() {
		k := "k"
		ttl := time.Minute
		hk := func(values NamedValues, key string) interface{} {
			if k == key {
				return time.Now().Add(ttl)
			}
			return nil
		}
		mkv := NewSimpleMemNamedValuesWithHook(hk)
		v1o, err1 := mkv.Get(k)
		So(err1, ShouldBeNil)
		time.Sleep(1 * time.Millisecond)
		v2o, err2 := mkv.Get(k)
		So(err2, ShouldBeNil)
		v1, ok1 := v1o.(time.Time)
		v2, ok2 := v2o.(time.Time)
		So(ok1 && ok2, ShouldBeTrue)
		So(v1.UnixMilli(), ShouldBeLessThan, v2.UnixMilli())
	})
}

func TestSortedMemNamedValues(t *testing.T) {
	Convey("sortedMemNamedValues", t, func() {
		Convey("Set OK", func() {
			skv := NewSortedMemNamedValues()
			kvs := map[string]interface{}{
				"a":                  1,
				"b":                  true,
				"c":                  nil,
				"d":                  utils.RandAscii(10),
				utils.RandAscii(100): utils.RandInt64Below(5e5),
			}
			for k, v := range kvs {
				So(skv.Set(k, v), ShouldBeNil)
			}
		})

		Convey("Get OK", func() {
			skv := NewSortedMemNamedValues()
			ks := []string{
				"a", "b", "", utils.RandAscii(10),
			}
			for _, k := range ks {
				_, err := skv.Get(k)
				So(errors.Is(err, errors2.ErrEntryNotFound), ShouldBeTrue)
			}
		})

		Convey("Get your Set", func() {
			skv := NewSortedMemNamedValues()
			kvs := map[string]interface{}{
				"a":                  1,
				"b":                  true,
				"c":                  nil,
				"d":                  utils.RandAscii(10),
				utils.RandAscii(100): utils.RandInt64Below(5e5),
			}
			for k, v := range kvs {
				assert.Nil(t, skv.Set(k, v), "%s -> %+v", k, v)
			}
			for k, v := range kvs {
				valueFound, err := skv.Get(k)
				So(err, ShouldBeNil)
				So(v, ShouldEqual, valueFound)
			}
		})

		Convey("Keys() return keys by Set Order", func() {
			skv := NewSortedMemNamedValues()
			ks := []string{
				"a", "b", "c", utils.RandAscii(10),
			}
			Convey("Set should be OK", func() {
				for _, k := range ks {
					So(skv.Set(k, utils.RandAscii(10)), ShouldBeNil)
				}
				Convey("Keys should be OK", func() {
					ksGot, err := skv.Keys()
					So(err, ShouldBeNil)
					So(len(ksGot), ShouldEqual, len(ks))
					for i, kGot := range ksGot {
						So(kGot, ShouldEqual, ks[i])
					}
					Convey("add again", func() {
						rand.Shuffle(len(ks), func(i, j int) {
							ks[i], ks[j] = ks[j], ks[i]
						})
						Convey("Set", func() {
							for _, k := range ks {
								So(skv.Set(k, utils.RandInt64Below(10)), ShouldBeNil)
							}
							Convey("Keys should be sorted", func() {
								keysGot, errGot := skv.Keys()
								So(errGot, ShouldBeNil)
								So(len(keysGot), ShouldEqual, len(ks))
								for i, keyGot := range keysGot {
									So(keyGot, ShouldEqual, ks[i])
								}
							})
						})
					})
				})
			})
		})
	})
}
