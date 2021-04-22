package connections

import (
	"github.com/Confialink/wallet-currencies/internal/srvdiscovery"
	"net/http"

	pb "github.com/Confialink/wallet-accounts/rpc/accounts"
)

func GetAccountsClient() (pb.AccountsProcessor, error) {
	accountsUrl, err := srvdiscovery.ResolveRPC(srvdiscovery.ServiceNameAccounts)
	if err != nil {
		return nil, err
	}
	return pb.NewAccountsProcessorProtobufClient(accountsUrl.String(), http.DefaultClient), nil
}
