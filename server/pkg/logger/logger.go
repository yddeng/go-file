package logger

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.luxshare-ict.com/pkg/config"
	"go.luxshare-ict.com/pkg/logger/color"
	"runtime"
	"strings"
	"sync"
	"time"
)

var filePath string
var serverState sync.Map

// EsOption represents ElasticSearch Options
type EsOptions struct {
	URL   string `yaml:"url"`
	Host  string `yaml:"host"` // 本机
	Index string `yaml:"index"`
}

func PutState(key, val string) {
	serverState.Store(key, val)
}
func GetState() sync.Map {
	return serverState
}
func PrintState() {
	css := color.BlueFont.Echo
	fmt.Println(css("server info:"))
	fmt.Println(css("===================================================================================="))
	for _, key := range []string{"proto", "evn", "ver", "asset", "logs", "config", "", "gin", "http", "grpc", "fans", "socket", "tcp", "udp"} {
		if key == "" {
			fmt.Println("")
		} else if value, ok := serverState.Load(key); ok {
			fmt.Println(css(fmt.Sprintf(">> %v: \t%v", key, value)))
		}
	}
	fmt.Println(css("===================================================================================="))
	fmt.Println("")
}

func print(fn func(s string, v ...interface{}), args []interface{}) {
	l := len(args)
	f := make([]string, l)
	for i := 0; i < l; i++ {
		f[i] = "%+v"
	}
	fn(strings.Join(f, ", "), args...)
}

func callerPrettyfier(frame *runtime.Frame) (function string, file string) {
	fns := strings.Split(frame.Function, ".")
	//function = fns[len(fns)-1]
	fn := fns[len(fns)-1]
	f := strings.SplitAfter(frame.File, "LuxGoMES/server/")
	if len(f) > 1 {
		file = f[1]
	} else {
		file = f[0]
	}
	file = strings.TrimRight(file, ".go") + "." + fn
	return
}
func Init(c *config.Conf) error {
	logrus.SetReportCaller(true)
	/*if c.App.Evn == "prod" {
		logrus.SetFormatter(&logrus.JSONFormatter{
			CallerPrettyfier: callerPrettyfier,
		})
	} else {*/
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
		//TimestampFormat:  "01-02 15:04:05",
		//FullTimestamp:    true,
		QuoteEmptyFields: true,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) { file = ":"; return },
	})
	//}

	level, err := logrus.ParseLevel(c.Log.Level)
	if err != nil {
		return errors.Wrap(err, "invalid errLog level")
	}
	logrus.SetLevel(level)
	if c.App.Evn == "prod" && len(c.Log.Path) > 0 {
		logrus.SetFormatter(&logrus.JSONFormatter{
			CallerPrettyfier: callerPrettyfier,
		})
		filePath = c.Log.Path
		filePattern := fmt.Sprintf("%s/app_%%Y%%m%%d.log", c.Log.Path)
		linkName := fmt.Sprintf("%s/app.log", c.Log.Path)
		if c.Log.MaxDay <= 0 {
			c.Log.MaxDay = 7
		}
		if c.Log.MaxSize <= 0 {
			c.Log.MaxSize = 100
		}
		if c.Log.MaxCount <= 1 {
			c.Log.MaxCount = uint(c.Log.MaxDay) * 3
		}

		opts := []rotatelogs.Option{
			rotatelogs.WithLinkName(linkName),
			rotatelogs.WithRotationTime(24 * time.Hour),
			rotatelogs.WithRotationSize(c.Log.MaxSize * 1024),
		}
		if c.Log.MaxCount > 0 {
			opts = append(opts, rotatelogs.WithRotationCount(c.Log.MaxCount))
		} else if c.Log.MaxDay > 0 {
			opts = append(opts, rotatelogs.WithMaxAge(time.Duration(c.Log.MaxDay)*24*time.Hour))
		}
		logf, err := rotatelogs.New(filePattern, opts...)
		if err != nil {
			return errors.Wrap(err, "create errLog rotation file failed")
		}
		logrus.SetOutput(logf)
	}

	/*
		if c.Es != nil {
			if len(c.Es.URL) == 0 {
				return errors.New("es url is missing")
			}
			if len(c.Es.Host) == 0 {
				return errors.New("es host is missing")
			}
			if len(c.Es.Index) == 0 {
				return errors.New("es index is missing")
			}

			client, err := elastic.NewClient(elastic.SetURL(c.Es.URL))
			if err != nil {
				return errors.Wrap(err, "new es client failed")
			}
			esHook, err := elogrus.NewAsyncElasticHook(client, c.Es.Host, level, c.Es.Index)
			if err != nil {
				return errors.Wrap(err, "new async es hook failed")
			}
			logrus.AddHook(esHook)
		}
	*/
	PutState("asset", c.App.AssetPath)
	PutState("logs", c.Log.Path)
	PutState("evn", c.App.Evn)
	PutState("proto", c.App.PB)
	PutState("ver", c.App.Version)
	PutState("config", c.File)

	return nil
}

func WithError(err error) *logrus.Entry {
	return logrus.WithError(err)
}

func Mark(err error, tag string) bool {
	if err != nil {
		Error(tag + " 失败")
		return true
	}
	Info(tag + " 成功")
	return false
}

func Info(args ...interface{}) {
	print(Infof, args)
}
func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}
func Error(args ...interface{}) {
	print(Errorf, args)
}
func Errorf(format string, args ...interface{}) {
	logrus.Errorf(color.RedFont.Echo(format), args...)
}
func Warning(args ...interface{}) {
	print(Warningf, args)
}

func Warningf(format string, args ...interface{}) {
	logrus.Warningf(color.YellowFont.Echo(format), args...)
}

func Print(args ...interface{}) {
	print(Printf, args)
}
func Printf(format string, args ...interface{}) {
	logrus.Printf(color.GreenFont.Echo(format), args...)
}
