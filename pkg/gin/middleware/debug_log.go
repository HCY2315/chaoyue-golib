package middleware

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"

	"github.com/HCY2315/chaoyue-golib/pkg/log"
	"github.com/gin-gonic/gin"
)

const CtxKeyDebugReqBody = "chaoyue:debug:body"

//RequestBodyToContext 读取请求body到context, 只能在debug阶段使用
func RequestBodyToContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, errReadBody := ioutil.ReadAll(c.Request.Body)
		if errReadBody != nil {
			log.Warnf("读取request body失败[%s], skip", errReadBody.Error())
			return
		}
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf)) //We have to create a new Buffer, because rdr1 will be read.

		c.Set(CtxKeyDebugReqBody, readBody(rdr1)) // Print request body

		c.Request.Body = rdr2
		c.Next()
	}
}

func readBody(reader io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	s := buf.String()
	return s
}

func DebugLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Debugf("")
		log.Debugf(strings.Repeat(">", 20))
		log.Debugf("[%s] %s", ctx.Request.Method, ctx.Request.URL.Path)
		log.Debugf("[Context] :" + ctx.GetString(CtxKeyDebugReqBody))
		log.Debugf(strings.Repeat("<", 20))
		log.Debugf("")
	}
}
