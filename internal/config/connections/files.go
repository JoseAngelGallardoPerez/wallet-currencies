package connections

import (
	"github.com/Confialink/wallet-currencies/internal/srvdiscovery"
	"net/http"

	pb "github.com/Confialink/wallet-files/rpc/files"
)

func GetFilesClient() (pb.ServiceFiles, error) {
	filesUrl, err := srvdiscovery.ResolveRPC(srvdiscovery.ServiceNameFiles)
	if err != nil {
		return nil, err
	}

	return pb.NewServiceFilesProtobufClient(filesUrl.String(), http.DefaultClient), nil
}
