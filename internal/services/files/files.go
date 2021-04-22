package files

import (
	"context"

	"github.com/inconshreveable/log15"

	pb "github.com/Confialink/wallet-files/rpc/files"

	"github.com/Confialink/wallet-currencies/internal/config/connections"
)

type Service struct {
	filesServer pb.ServiceFiles
	logger      log15.Logger
}

func NewService(logger log15.Logger) *Service {
	return &Service{logger: logger.New("Service", "Files")}
}

func (s *Service) DownloadFileById(id uint64) (*File, error) {
	logger := s.logger.New("method", "DownloadFileById")
	request := pb.FileReq{Id: id}
	resp, err := s.processor().DownloadFile(context.Background(), &request)
	if err != nil {
		logger.Error("Failed to download file", "error", err)
		return nil, err
	}

	return &File{
		Data:        resp.Data,
		Size:        resp.Size,
		ContentType: resp.ContentType,
	}, nil
}

func (s *Service) processor() pb.ServiceFiles {
	if s.filesServer == nil {
		connection, err := connections.GetFilesClient()
		if err != nil {
			s.logger.Error("Failed to connect to accounts", "error", err)
			return nil
		}

		s.filesServer = connection
	}

	return s.filesServer
}
