package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	domainValidator "github.com/lazyjean/sla2/internal/domain/validator"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// 注册自定义验证器
	// 例如：密码强度验证
	validate.RegisterValidation("password_strength", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		hasUpper := strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
		hasLower := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
		hasNumber := strings.ContainsAny(password, "0123456789")
		return hasUpper && hasLower && hasNumber
	})
}

// Validate 验证结构体
func Validate(s interface{}) error {
	return validate.Struct(s)
}

// TranslateError 翻译验证错误
func TranslateError(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := e.Field()
			tag := e.Tag()

			// 获取字段的 json 标签
			fieldType, _ := reflect.TypeOf(e.Value()).FieldByName(field)
			jsonTag := fieldType.Tag.Get("json")
			if jsonTag != "" {
				field = jsonTag
			}

			// 根据验证标签生成错误消息
			switch tag {
			case "required":
				errors[field] = fmt.Sprintf("%s 不能为空", field)
			case "min":
				errors[field] = fmt.Sprintf("%s 长度不能小于 %s", field, e.Param())
			case "max":
				errors[field] = fmt.Sprintf("%s 长度不能大于 %s", field, e.Param())
			case "email":
				errors[field] = fmt.Sprintf("%s 必须是有效的邮箱地址", field)
			case "password_strength":
				errors[field] = fmt.Sprintf("%s 必须包含大小写字母和数字", field)
			default:
				errors[field] = fmt.Sprintf("%s 验证失败", field)
			}
		}
	}

	return errors
}

type Validator struct {
	validator *validator.Validate
}

func NewValidator() domainValidator.Validator {
	return &Validator{
		validator: validator.New(),
	}
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func (v *Validator) RegisterValidation(tag string, fn interface{}) error {
	return v.validator.RegisterValidation(tag, fn.(validator.Func))
}
