package handler

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`           // 业务状态码
	Message string      `json:"message"`        // 响应消息
	Data    interface{} `json:"data,omitempty"` // 响应数据
}

// ListResponse 列表响应结构
type ListResponse struct {
	Items    interface{} `json:"items"`     // 列表数据
	Page     int         `json:"page"`      // 当前页码
	PageSize int         `json:"page_size"` // 每页数量
	Total    int64       `json:"total"`     // 总记录数
}

// NewResponse 创建新的响应
func NewResponse(code int, message string, data interface{}) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// NewListResponse 创建新的列表响应
func NewListResponse(items interface{}, page, pageSize int, total int64) *Response {
	return NewResponse(0, "success", ListResponse{
		Items:    items,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	})
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(code int, message string) *Response {
	return NewResponse(code, message, nil)
}
