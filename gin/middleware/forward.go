package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/HCY2315/chaoyue-golib/log"
	cecfHTTP "github.com/HCY2315/chaoyue-golib/pkg/http"
	"github.com/gin-gonic/gin"
	"github.com/vulcand/oxy/utils"
)

type ReqRewriterSet struct {
	rewriterList []ReqRewriteFunc
}

type ReqRewriteFunc func(*http.Request)

func RefererHeaderRewrite(referer string) ReqRewriteFunc {
	return func(req *http.Request) {
		req.Header.Set(cecfHTTP.HeaderKeyReferer, referer)
	}
}

func SchemaHostRewrite(schemaHost string) ReqRewriteFunc {
	urlParsed, err := url.Parse(schemaHost)
	if err != nil {
		panic(fmt.Sprintf("parse [%s] to url failed:%s", schemaHost, err.Error()))
	}

	return func(req *http.Request) {
		req.URL.Scheme = urlParsed.Scheme
		req.URL.Host = urlParsed.Host
	}
}

type ForwardErrHandler utils.ErrorHandler

type LogErrHandler struct {
}

func NewLogErrHandler() *LogErrHandler {
	return &LogErrHandler{}
}

func (l LogErrHandler) ServeHTTP(w http.ResponseWriter, req *http.Request, err error) {
	log.Errorf("forward to [%s] failed:[%s]", req.URL.String(), err.Error())
}

func ForwardToSingle(backendURL *url.URL, errHandler ForwardErrHandler, ctxRewrite gin.HandlerFunc,
	rws ...ReqRewriteFunc) gin.HandlerFunc {

	proxy := httputil.NewSingleHostReverseProxy(backendURL)
	return func(c *gin.Context) {
		originURL := c.Request.URL.String()
		ctxRewrite(c)
		for _, rw := range rws {
			rw(c.Request)
		}
		finalURL := c.Request.URL.String()
		log.Debugf("forward %s => %s", originURL, finalURL)
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
