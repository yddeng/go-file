package account

import (
	"context"
	"github.com/pkg/errors"
	v1 "go.luxshare-ict.com/app/service/account/api/v1"
	"go.luxshare-ict.com/app/service/account/internal/dao"
	"go.luxshare-ict.com/pkg/app/buildInfo"
	mdb "go.luxshare-ict.com/pkg/db/mongo"
)

type Service struct {
	dao *dao.Dao
}

func New(d *dao.Dao) (s *Service) {
	return &Service{
		dao: d,
	}
}

func (s *Service) Close() {
	s.dao.Close()
}

func (s *Service) Login(ctx context.Context, req *v1.LoginReq) (resp *v1.ReplyResp, err error) {
	err = s.dao.Mongo(ctx, func(conn *mdb.Conn) (err error) {
		resp, err = s.dao.Login(conn, req)
		bi := buildInfo.GetInfo()
		msg := bi.Config.Protocols["grpc"].Map.String("message")
		if err == nil {
			resp.Message = msg
		} else {
			err = errors.Wrap(err, msg)
		}
		return
	})
	return
}

func (s *Service) Logout(ctx context.Context, req *v1.LogoutReq) (resp *v1.ReplyResp, err error) {
	err = s.dao.Mongo(ctx, func(conn *mdb.Conn) (err error) {
		resp, err = s.dao.Logout(conn, req)
		return
	})
	return
}
func (s *Service) CheckToken(ctx context.Context, req *v1.CheckTokenReq) (resp *v1.ReplyResp, err error) {
	err = s.dao.Mongo(ctx, func(conn *mdb.Conn) (err error) {
		resp, err = s.dao.CheckToken(conn, req)
		return
	})
	return
}
func (s *Service) ForkToken(ctx context.Context, req *v1.ForkTokenReq) (resp *v1.ReplyResp, err error) {
	err = s.dao.Mongo(ctx, func(conn *mdb.Conn) (err error) {
		resp, err = s.dao.ForkToken(conn, req)
		return
	})
	return
}
