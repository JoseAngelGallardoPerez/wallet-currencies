package abilities

import (
	"github.com/Confialink/wallet-users/rpc/proto/users"
	"github.com/gin-gonic/gin"
)

func Can(context *gin.Context, action, resource string) bool {
	user := context.MustGet("_user").(*users.User)
	roleAbilities, ok := abilitiesList[user.RoleName]
	if !ok {
		return false
	}
	resourceAbilities, ok := roleAbilities[resource]
	if !ok {
		return false
	}
	function, ok := resourceAbilities[action]
	if !ok {
		return false
	}
	return function(user)
}

func CanNot(context *gin.Context, action, resource string) bool {
	return !Can(context, action, resource)
}
