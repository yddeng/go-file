package validator

import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	en_treanslations "gopkg.in/go-playground/validator.v9/translations/en"
	zh_treanslations "gopkg.in/go-playground/validator.v9/translations/zh"
)

type Local string

const (
	LocalZH Local = "zh"
	LocalEN Local = "en"
)

func New(local Local) (validate *validator.Validate, trans ut.Translator) {
	validate = validator.New()
	trans = NewTranslator(local)

	var err error
	switch local {
	case LocalZH:
		err = zh_treanslations.RegisterDefaultTranslations(validate, trans)
	case LocalEN:
		err = en_treanslations.RegisterDefaultTranslations(validate, trans)
	}
	if err != nil {
		panic(err)
	}

	return validate, trans
}

func NewTranslator(local Local) (trans ut.Translator) {
	switch local {
	case LocalZH:
		_zh := zh.New()
		uni := ut.New(_zh, _zh)
		trans, _ = uni.GetTranslator("zh")
	case LocalEN:
		_en := en.New()
		uni := ut.New(_en, _en)
		trans, _ = uni.GetTranslator("en")
	}
	return trans
}
