package resolvers

import (
	"github.com/Confialink/wallet-users/rpc/proto/users"

	"github.com/Confialink/wallet-currencies/internal/abilities/resolvers/pbchecker"
)

const ViewSettingsPermissionKey = "view_settings"
const EditSettingsPermissionKey = "modify_settings"

func Allow(_ *users.User) bool {
	return true
}

func CanEditSettings(user *users.User) bool {
	return pbchecker.HasPermission(user, EditSettingsPermissionKey)
}

func CanViewSettings(user *users.User) bool {
	return pbchecker.HasPermission(user, ViewSettingsPermissionKey)
}
