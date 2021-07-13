package dao

import (
	"context"
	"fmt"
	"go.luxshare-ict.com/pkg/config"
	"go.luxshare-ict.com/pkg/db/mongo"
	"go.luxshare-ict.com/pkg/db/redis"
	"go.luxshare-ict.com/pkg/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strings"
)

type Dao struct {
	conf     *config.Conf
	mdb      *mongo.MDB
	rdb      *redis.RDB
	isClosed bool
}

type Log struct{}

func (Log) Info(msg string, keysAndValues ...interface{}) {
	logger.Infof("cron info")
	logger.Infof("msg: %s, kv:%+v", msg, keysAndValues)
}

func (Log) Error(err error, msg string, keysAndValues ...interface{}) {
	logger.Infof("cron error")
	logger.Errorf("msg: %s, kv:%+v, err:%+v", msg, keysAndValues, err)
}
func New(cfg *config.Conf) *Dao {
	return &Dao{
		conf: cfg,
		mdb:  mongo.New(cfg.DB["mongodb"]),
		rdb:  redis.New(cfg.DB["redis"], nil),
	}
}

func serviceName(region, name string) string {
	return fmt.Sprintf("%s.%s", region, name)
}

func (d *Dao) Close() {

}

func (d *Dao) Mongo(ctx interface{}, h func(conn *mongo.Conn) (err error)) (err error) {
	return d.mdb.Conn(ctx, h)
}
func (d *Dao) Mongodb(ctx interface{}, h func(conn *mongo.Conn) (resp interface{}, err error)) (resp interface{}, err error) {
	err = d.mdb.Conn(ctx, func(c *mongo.Conn) (err error) {
		resp, err = h(c)
		return
	})
	return
}
func (d *Dao) MDB() *mongo.MDB {
	return d.mdb
}

func getIp(ctx context.Context) (ip string) {
	if t, ok := ctx.Value("ClientIP").(string); ok && t != "" {
		return t
	}
	return "127.0.0.1"
}

// 生成新的mongo id
func genID() string {
	return primitive.NewObjectID().Hex()
}
func getID(str ...string) string {
	if len(str) == 1 && !empty(str[0]) {
		return str[0]
	}
	return genID()
}

// not a status
func nas(i interface{}) bool {
	s := reflect.ValueOf(i).Int()
	if s <= 0 {
		return true
	}
	return s != 1<<(s>>1)
}
func str(val interface{}) string {
	return fmt.Sprint(val)
}
func ternaryStr(tf bool, tv, fv string) string {
	if tf {
		return tv
	}
	return fv
}
func empty(str string) bool {
	return len(str) <= 0 || len(strings.TrimSpace(str)) <= 0
}
func exist(str string) bool {
	return !empty(str)
}
func subString(str string, start, end int, join ...string) (res string) {
	res = ""
	s := ""
	if len(join) > 0 {
		s = join[0]
	}
	d := []rune(str)
	if start < 0 {
		start = 0
	}
	if end <= start {
		return
	}
	l := len(d)

	if l < start {
		return
	}
	if l < end {
		end = l
	}
	if end <= start {
		return
	}
	res = string(d[start:end])
	if l > end {
		res += s
	}
	return
}
func (d *Dao) SubString(str string, start, end int, join ...string) (res string) {
	return subString(str, start, end, join...)
}
func (d *Dao) GenID(str ...string) string {
	return getID(str...)
}

/**
如果 a 为空则为b,否则为a
*/
func ab(a string, b string) string {
	return ternaryStr(empty(a), b, a)
}
