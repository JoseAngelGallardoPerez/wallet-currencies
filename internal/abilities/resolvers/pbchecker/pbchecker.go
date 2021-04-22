package pbchecker

import (
	"github.com/Confialink/wallet-currencies/internal/srvdiscovery"
	"context"
	"net/http"

	"github.com/Confialink/wallet-permissions/rpc/permissions"
	"github.com/Confialink/wallet-users/rpc/proto/users"
	"github.com/inconshreveable/log15"

	"github.com/Confialink/wallet-currencies/internal/config/logger"
)

func HasPermission(user *users.User, permission string) bool {
	checker := getChecker()
	if checker == nil {
		return false
	}

	resp, err := checker.permissionsChecker.Check(context.Background(),
		&permissions.PermissionReq{
			UserId:    user.UID,
			ActionKey: permission,
		},
	)
	if err != nil {
		checker.logger.Error("Failed to get permission response", "error", err)
		return false
	}
	return resp.IsAllowed
}

type pbChecker struct {
	permissionsChecker permissions.PermissionChecker
	logger             log15.Logger
}

func getChecker() *pbChecker {
	newLogger := logger.Logger.New("Service", "pbChecker")
	permissionsUrl, err := srvdiscovery.ResolveRPC(srvdiscovery.ServiceNamePermissions)
	if err != nil {
		newLogger.Error("Failed to get permissions rpc url", "error", err)
		return nil
	}
	return &pbChecker{
		permissions.NewPermissionCheckerProtobufClient(permissionsUrl.String(), http.DefaultClient),
		newLogger,
	}
}
