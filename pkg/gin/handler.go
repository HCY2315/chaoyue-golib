package gin

import (
	"git.cestong.com.cn/cecf/cecf-golib/pkg/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GeneralGinHandleFunc func(*gin.Context) (interface{}, error)

func HandlerWrapper(handleFunc GeneralGinHandleFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := handleFunc(c)
		if err != nil {
			if ce, ok := err.(errors.ErrorWithCodeAndStatus); ok {
				c.JSON(ce.HTTPStatus(), GeneralResponseVO{
					Error: &EnhancedRespErr{
						Code: ce.Code(),
						Msg:  ce.Error(),
					},
					Data: nil,
				})
				c.Abort()
				return
			}
			c.Error(err)
			return
		}
		c.JSON(http.StatusOK, GeneralResponseVO{
			Error: nil,
			Data:  data,
		})
		return
	}
}
