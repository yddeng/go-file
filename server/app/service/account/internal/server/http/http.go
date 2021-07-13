package http

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	log "github.com/sirupsen/logrus"
	v1 "go.luxshare-ict.com/app/service/account/api/v1"
	"go.luxshare-ict.com/app/service/account/internal/service"
	"go.luxshare-ict.com/pkg/bizerror"
	response "go.luxshare-ict.com/pkg/response/gin"
	"go.luxshare-ict.com/pkg/validator"
	xvalidator "gopkg.in/go-playground/validator.v9"
	"net/http"
)

var (
	svc      *service.Service
	validate *xvalidator.Validate
	trans    ut.Translator
)

func Init(s *service.Service) {
	svc = s
}

type ApiConfig struct {
	Ifb     string `yaml:"ifb"'`
	Pms     string `yaml:"pms"'`
	Account string `yaml:"account"'`
}
type ApiSrvConfig struct {
	Http *ApiConfig `yaml:"http"`
	Rpc  *ApiConfig `yaml:"rpc"`
}

func init() {
	validate, trans = validator.New(validator.LocalZH)
	bizerror.SetTranslator(trans)
}

func CheckToken(c *gin.Context) {
	param := &v1.CheckTokenReq{}
	if checkErr(c, c.BindJSON(param), "Bind JSON failed", &bizerror.InvalidParam) {
		return
	}
	if checkErr(c, validate.Struct(*param), "param validate failed", &bizerror.InvalidParam) {
		return
	}
	resp, err := svc.Account.CheckToken(c, param)
	if checkErr(c, err, "checkToken failed", nil) {
		return
	}
	c.JSON(http.StatusOK, resp)
}

func asset(h gin.HandlerFunc) gin.HandlerFunc {
	return gin.HandlerFunc(
		func(c *gin.Context) {
			h(c)
		},
	)
}

func checkErr(c *gin.Context, err error, msg string, berr *bizerror.Error) bool {
	if err != nil {
		if berr != nil {
			if _, ok := err.(xvalidator.ValidationErrors); ok {
				err = berr.TransErr(err)
			} else {
				err = berr.FromErr(err)
			}
		}
		log.WithError(err).Infof(msg)
		response.JSON(c, err, nil)
		return true
	}
	return false
}
