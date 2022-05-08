package gin

import (
	"git.cestong.com.cn/cecf/cecf-golib/pkg/errors"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func UIntListFromQuery(c *gin.Context, key string) ([]uint, error) {
	str := c.Query(key)
	ss := strings.Split(str, ",")
	ret := make([]uint, 0, len(ss))
	for _, s := range ss {
		trimmed := strings.TrimSpace(s)
		if trimmed == "" {
			continue
		}
		i, err := strconv.ParseUint(trimmed, 10, 64)
		if err != nil {
			return nil, errors.Wrap(errors.ErrBadRequest, "%s not uint:%s", s, err.Error())
		}
		ret = append(ret, uint(i))
	}
	return ret, nil
}

func UIntPathParam(c *gin.Context, key string) (uint, error) {
	vs := c.Param(key)
	i64, errParse := strconv.ParseUint(vs, 10, 64)
	if errParse != nil {
		return 0, errors.Wrap(errors.ErrBadRequest, "invalid path param %s:%s", key, vs)
	}
	return uint(i64), nil
}
