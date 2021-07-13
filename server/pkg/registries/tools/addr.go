package tools

import (
	"encoding/json"
	"go.luxshare-ict.com/pkg/apeol.com/crypt"
)

type Addr struct {
	PB      string `json:"p"`
	Weight  int    `json:"w"`
	Addr    string `json:"s"`
	Version string `json:"v"`
	Evn     string `json:"e"`
}

func DecodeAddr(val []byte) (addr *Addr, err error) {
	addr = &Addr{}
	err = json.Unmarshal(val, addr)
	return
}

func (a *Addr) Fix() {
	if a.Weight <= 0 {
		a.Weight = 5
	} else if a.Weight > 10 {
		a.Weight = 10
	}
	if a.PB == "" {
		a.PB = "v1"
	}
}
func (a *Addr) ToString() (res string) {
	b, err := json.Marshal(a)
	if err != nil {
		return ""
	}
	return string(b)
}
func (a *Addr) Hash() string {
	if a == nil {
		return ""
	}
	return crypt.Md5(a.Evn, a.PB, a.Version, a.Addr)
}
