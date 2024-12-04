package ping_service

import "context"

type pingRepo interface {
	PingDB(ctx context.Context) (string, error)
}

type Service struct {
	PingRepo pingRepo
}

func New(pingRepo pingRepo) *Service {
	return &Service{
		PingRepo: pingRepo,
	}
}
