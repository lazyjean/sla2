package docs

// Swagger type definitions for protobuf types
type swaggerTimestamp struct {
	Seconds int64 `json:"seconds"`
	Nanos   int32 `json:"nanos"`
}

// Swagger type definitions for protobuf messages
type swaggerSessionResponse struct {
	SessionID    string           `json:"session_id"`
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	CreatedAt    swaggerTimestamp `json:"created_at"`
	UpdatedAt    swaggerTimestamp `json:"updated_at"`
	MessageCount uint64           `json:"message_count"`
}

// Swagger type definitions for protobuf requests
type swaggerCreateSessionRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Swagger type definitions for protobuf responses
type swaggerListSessionsResponse struct {
	Sessions   []swaggerSessionResponse `json:"sessions"`
	TotalCount uint32                   `json:"total_count"`
	Page       uint32                   `json:"page"`
	PageSize   uint32                   `json:"page_size"`
}

// Swagger type definitions for protobuf empty response
type swaggerEmpty struct{}
