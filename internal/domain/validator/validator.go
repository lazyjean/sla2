package validator

// Validator 定义了验证器的接口
type Validator interface {
	// Validate 验证给定的结构体
	Validate(i interface{}) error
	// RegisterValidation 注册自定义验证函数
	RegisterValidation(tag string, fn interface{}) error
}
