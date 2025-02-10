package handler

// Handlers 包含所有HTTP处理器
type Handlers struct {
	WordHandler     *WordHandler
	AuthHandler     *AuthHandler
	LearningHandler *LearningHandler
	HealthHandler   *HealthHandler
}

// NewHandlers 创建新的Handlers实例
func NewHandlers(
	wordHandler *WordHandler,
	authHandler *AuthHandler,
	learningHandler *LearningHandler,
	healthHandler *HealthHandler,
) *Handlers {
	return &Handlers{
		WordHandler:     wordHandler,
		AuthHandler:     authHandler,
		LearningHandler: learningHandler,
		HealthHandler:   healthHandler,
	}
}
