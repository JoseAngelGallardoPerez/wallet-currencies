package middleware

import (
	"net/http"

	"github.com/Confialink/wallet-pkg-acl"
	"github.com/Confialink/wallet-users/rpc/proto/users"
	"github.com/gin-gonic/gin"
)

func AdminOnly(c *gin.Context) {
	user, exist := c.Get("_user")
	if !exist {
		c.Status(http.StatusForbidden)
		c.Abort()
		return
	}

	role := acl.RolesHelper.FromName((user.(*users.User)).RoleName)
	if role < acl.Admin {
		c.Status(http.StatusForbidden)
		c.Abort()
		return
	}
}
