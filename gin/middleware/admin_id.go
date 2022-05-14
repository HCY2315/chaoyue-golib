package middleware

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type AdminIDFinder func(c *gin.Context) (uint, bool)

const (
	CtxAdminIdKey = "hx:adminID"
)

const (
	debugAdminIDKey = "c2NvcGUiOiIiLCJl"
)

func BuildConstAdminIDFinderFunc(adminID uint) AdminIDFinder {
	return func(c *gin.Context) (uint, bool) {
		return adminID, true
	}
}

func DebugAdminIDFindFunc(c *gin.Context) (uint, bool) {
	adminIDHeader := c.GetHeader(debugAdminIDKey)
	if adminIDHeader == "" {
		return 0, false
	}
	adminID, err := strconv.ParseUint(adminIDHeader, 10, 64)
	return uint(adminID), err == nil
}

func BuildSetAdminIDMiddleware(finders ...AdminIDFinder) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, finder := range finders {
			adminID, find := finder(c)
			if find {
				c.Set(CtxAdminIdKey, adminID)
				break
			}
		}
		c.Next()
	}
}
