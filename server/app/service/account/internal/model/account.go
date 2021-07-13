package model

import (
	v1 "go.luxshare-ict.com/app/service/account/api/v1"
	"go.luxshare-ict.com/pkg/apeol.com/slice"
	"go.luxshare-ict.com/pkg/time"
)

type Account struct {
	Account      string    `json:"account" xorm:"id" bson:"_id"`
	Name         string    `json:"name" xorm:"name" bson:"name"`
	Activate     int32     `json:"activate" xorm:"activate" bson:"activate"`
	Expire       time.Time `json:"expire" xorm:"expire" bson:"expire"`
	Group        int32     `json:"group" xorm:"group" bson:"group"`
	Role         string    `json:"role" xorm:"role" bson:"role"`
	Tel          string    `json:"tel" xorm:"tel" bson:"tel"`
	Email        string    `json:"email" xorm:"email" bson:"email"`
	Position     string    `json:"position" xorm:"position" bson:"position"`
	PositionId   string    `json:"positionId" xorm:"position_id" bson:"position_id"`
	Department   string    `json:"department" xorm:"department" bson:"department"`
	DepartmentId string    `json:"departmentId" xorm:"department_id" bson:"department_id"`
	Hiredate     time.Time `json:"hiredate" xorm:"hiredate" bson:"hiredate"`
	Education    string    `json:"education" xorm:"education" bson:"education"`
	Gender       string    `json:"gender" xorm:"gender" bson:"gender"`
	Party        string    `json:"party" xorm:"party" bson:"party"`
	IdCard       string    `json:"idCard" xorm:"id_card" bson:"id_card"`
	Hometown     string    `json:"hometown" xorm:"hometown" bson:"hometown"`
	Disable      int32     `json:"disable" xorm:"disable" bson:"disable"`
	Salt         string    `json:"salt" xorm:"salt" bson:"salt"`
	Auth         []string  `json:"auth" xorm:"auth" bson:"auth"`
	UnAuth       []string  `json:"unAuth" xorm:"un_auth" bson:"un_auth"`
	Bu           string    `json:"bu" xorm:"bu" bson:"bu"`
	Birthday     time.Time `json:"birthday" xorm:"birthday" bson:"birthday"`
	Oaid         string    `json:"oaid" xorm:"oaid" bson:"oaid"`
	Version      int       `json:"version" xorm:"version" bson:"version"`
}
type SNS struct {
	Id              string `json:"id" xorm:"id" bson:"_id"`
	Type            string `json:"type" xorm:"type" bson:"type"`
	Account         string `json:"account" xorm:"account" bson:"account"`
	OpenId          string `json:"openId" xorm:"open_id" bson:"open_id"`
	AccessToken     string `json:"accessToken" xorm:"access_token" bson:"access_token"`
	AccessTokenExp  int64  `json:"accessTokenExp" xorm:"access_token_exp" bson:"access_token_exp"`
	RefreshToken    string `json:"refreshToken" xorm:"refresh_token" bson:"refresh_token"`
	RefreshTokenExp int64  `json:"refreshTokenExp" xorm:"refresh_token_exp" bson:"refresh_token_exp"`
}

func (Account) FromProto(from *v1.Account) *Account {
	if from == nil {
		return nil
	}
	return &Account{
		Account:      from.Account,
		Name:         from.Name,
		Activate:     from.Activate,
		Expire:       from.Expire,
		Group:        from.Group,
		Role:         from.Role,
		Tel:          from.Tel,
		Email:        from.Email,
		Position:     from.Position,
		PositionId:   from.PositionId,
		Department:   from.Department,
		DepartmentId: from.DepartmentId,
		Hiredate:     from.Hiredate,
		Education:    from.Education,
		Gender:       from.Gender,
		Party:        from.Party,
		IdCard:       from.IdCard,
		Hometown:     from.Hometown,
		Disable:      from.Disable,
		Bu:           from.Bu,
		Oaid:         from.Oaid,
	}
}
func (from *Account) ToProto(withall ...bool) *v1.Account {
	if from == nil {
		return nil
	}
	if len(withall) > 0 && withall[0] {
		return &v1.Account{
			Account:      from.Account,
			Name:         from.Name,
			Activate:     from.Activate,
			Expire:       from.Expire,
			Group:        from.Group,
			Role:         from.Role,
			Tel:          from.Tel,
			Email:        from.Email,
			Position:     from.Position,
			PositionId:   from.PositionId,
			Department:   from.Department,
			DepartmentId: from.DepartmentId,
			Hiredate:     from.Hiredate,
			Education:    from.Education,
			Gender:       from.Gender,
			Party:        from.Party,
			IdCard:       from.IdCard,
			Hometown:     from.Hometown,
			Disable:      from.Disable,
			Bu:           from.Bu,
			Oaid:         from.Oaid,
		}
	}
	return &v1.Account{
		Account:      from.Account,
		Name:         from.Name,
		Activate:     from.Activate,
		Group:        from.Group,
		Role:         from.Role,
		Tel:          from.Tel,
		Email:        from.Email,
		Position:     from.Position,
		PositionId:   from.PositionId,
		Department:   from.Department,
		DepartmentId: from.DepartmentId,
		Gender:       from.Gender,
		Disable:      from.Disable,
		Bu:           from.Bu,
		Oaid:         from.Oaid,
	}
}
func (from *Account) CheckPermission(auth ...string) (has bool) {
	if from == nil {
		return false
	}
	if len(auth) <= 0 {
		return true
	}
	if from.Account == "system" {
		return true
	}
	perm := slice.String(from.Auth)
	if perm.Contains("all") {
		// todo 这里可以，也需要改一下，可以改成字典，走字典，只允许核心人员管理这个权限
		if from.Bu != "江苏机器人" {
			for _, a := range auth {
				if a == "systemSetup" {
					return false
				}
			}
		}
		return true
	}
	disable := slice.String(from.UnAuth)
	if disable.Contains(auth...) {
		return false
	}
	if perm.Contains(auth...) {
		if from.Bu != "江苏机器人" {
			for _, a := range auth {
				if a == "systemSetup" {
					return false
				}
			}
		}
		return true
	}
	return false
}
