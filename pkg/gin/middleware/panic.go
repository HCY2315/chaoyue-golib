package middleware

import (
	"fmt"
	"git.cestong.com.cn/cecf/cecf-golib/pkg/errors"
	"git.cestong.com.cn/cecf/cecf-golib/pkg/log"
	"git.cestong.com.cn/cecf/cecf-golib/pkg/thirdparty"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PanicHandler func(r interface{}, c *gin.Context)

func BuildDingMsgPanicHandler(service thirdparty.DingService) PanicHandler {
	return func(r interface{}, c *gin.Context) {
		panicMsg := errors.PanicMsg(r)
		service.SendMarkdown("服务Panic", fmt.Sprintf("###错误\n%s", panicMsg))
	}
}

func BuildPanicMiddleware(handlers ...PanicHandler) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		for _, handler := range handlers {
			handler(recovered, c)
		}
		log.Errorf("[CRIT]: panic recovered:%+v", recovered)
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}
