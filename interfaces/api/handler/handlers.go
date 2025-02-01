package handler

// Handlers 包含所有HTTP处理器
type Handlers struct {
	WordHandler     *WordHandler
	AuthHandler     *AuthHandler
	LearningHandler *LearningHandler
}

// NewHandlers 创建新的Handlers实例
func NewHandlers(
	wordHandler *WordHandler,
	authHandler *AuthHandler,
	learningHandler *LearningHandler,
) *Handlers {
	return &Handlers{
		WordHandler:     wordHandler,
		AuthHandler:     authHandler,
		LearningHandler: learningHandler,
	}
}
