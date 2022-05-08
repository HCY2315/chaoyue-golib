package middleware

import (
	"fmt"
	"git.cestong.com.cn/cecf/cecf-golib/pkg/log"
	"git.cestong.com.cn/cecf/cecf-golib/pkg/thirdparty"
	"github.com/gin-gonic/gin"
	"strings"
)

type statusFilterFunc func(int) bool

func RangeStatusFilter(lower, upper int) statusFilterFunc {
	if lower >= upper {
		panic(fmt.Sprintf("lower(%d)>=upper(%d)", lower, upper))
	}
	return func(i int) bool {
		return i >= lower && i < upper
	}
}

type statusListenerFunc func(int, *gin.Context)

func SendDingMsg(dingService thirdparty.DingService, id string, tags ...string) statusListenerFunc {
	return func(status int, c *gin.Context) {
		var sb strings.Builder
		for i, tag := range tags {
			sb.WriteString("`")
			sb.WriteString(tag)
			sb.WriteString("`")
			if i < len(tags)-1 {
				sb.WriteRune('\t')
			}
		}
		title := "服务器处理出错"
		markdownText := fmt.Sprintf(" ## %s\n #### id: **%s**\n #### tags: %s\n #### status: **%d**\n #### url: %s %s \n ",
			title, id, sb.String(), status, c.Request.Method, c.Request.URL.String())
		if err := dingService.SendMarkdown(title, markdownText); err != nil {
			log.Errorf("发送钉钉告警失败:%s", err.Error())
		}
	}
}

func BuildStatusMonitorMiddleware(statusFilter statusFilterFunc, listeners ...statusListenerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		status := c.Writer.Status()
		if statusFilter(status) {
			for _, lis := range listeners {
				lis(status, c)
			}
		}
	}
}
