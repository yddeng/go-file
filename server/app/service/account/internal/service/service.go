package service

import (
	"go.luxshare-ict.com/app/service/account/internal/dao"
	"go.luxshare-ict.com/app/service/account/internal/service/account"
	"go.luxshare-ict.com/pkg/config"
)

type Service struct {
	config  *config.Conf
	dao     *dao.Dao
	Account *account.Service
}

func New(cfg *config.Conf) (s *Service) {
	d := dao.New(cfg)
	s = &Service{
		config:  cfg,
		dao:     d,
		Account: account.New(d),
	}
	return
}

func (s *Service) Close() {
	s.dao.Close()
}
