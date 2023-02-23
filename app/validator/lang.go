package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	chTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var trans ut.Translator

// LoadValidatorLocal 初始化语言包
func LoadValidatorLocal() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zhT := zh.New() //chinese
		enT := en.New() //english
		uni := ut.New(enT, zhT, enT)

		var o bool
		local := "zh"
		trans, o = uni.GetTranslator(local)
		if !o {
			return errors.New("uni.GetTranslator failed")
		}
		// register translate
		// 注册翻译器
		var err error
		switch local {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = chTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = chTranslations.RegisterDefaultTranslations(v, trans)
		}

		if err != nil {
			//variable.ZapLog.Error("初始化验证语言失败", zap.Error(err))
			log.Println("初始化验证语言失败")
			return err
		}
		return nil
	}
	return nil
}

func Translate(err error) (errMsg string) {
	switch err.(type) {
	case validator.ValidationErrors:
		errs := err.(validator.ValidationErrors)
		for _, err := range errs {
			errMsg = err.Translate(trans)
			break
		}
	case *json.UnmarshalTypeError:
		errMsg = fmt.Sprintf("字段 %s 类型匹配有误，正确类型 %v", err.(*json.UnmarshalTypeError).Field, err.(*json.UnmarshalTypeError).Type)
	default:
		errMsg = err.Error()
	}
	return
}
