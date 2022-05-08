package auth

import (
	"testing"
	"time"

	"github.com/HCY2315/chaoyue-golib/pkg/common"
	"github.com/HCY2315/chaoyue-golib/pkg/utils"
	. "github.com/smartystreets/goconvey/convey"
)

func TestJwtTokenBuilder(t *testing.T) {
	Convey("JWTTokenBuilder", t, func() {
		Convey("token string", func() {
			ttl := 1 * time.Minute
			signSecret := utils.RandBytes(64)
			aesPasswd := utils.RandBytes(24)
			aesCrypto, errFilter := common.NewGCMCryptTextFilterByPasswd(aesPasswd, common.NewBase64URLFilter())
			So(errFilter, ShouldBeNil)
			tb := NewJwtTokenBuilder(ttl, signSecret, WithFilter(aesCrypto))
			tokenValues := common.NewSimpleMemNamedValues()
			testValues := map[string]interface{}{
				"a":                 "b",
				"b":                 utils.RandAscii(10),
				"c":                 1,
				"d":                 nil,
				"e":                 0,
				"f":                 false,
				utils.RandAscii(20): utils.RandInt64Below(9e10),
			}
			for k, v := range testValues {
				So(tokenValues.Set(k, v), ShouldBeNil)
			}
			Convey("get token", func() {
				tokenString, errToken := tb.Token(tokenValues)
				So(errToken, ShouldBeNil)
				So(tokenString, ShouldNotBeBlank)
				Convey("parse token", func() {
					valueParsed, errParse := tb.Parse(tokenString)
					So(errParse, ShouldBeNil)
					keysParsed, errKeys := valueParsed.Keys()
					So(errKeys, ShouldBeNil)
					So(keysParsed, ShouldNotBeEmpty)
					for k, v := range testValues {
						valueGet, errGetValue := valueParsed.Get(k)
						So(errGetValue, ShouldBeNil)
						So(valueGet, ShouldEqual, v)
					}
				})
			})
		})
	})
}
