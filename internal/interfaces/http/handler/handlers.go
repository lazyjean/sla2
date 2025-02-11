package handler

// Handlers 包含所有HTTP处理器
type Handlers struct {
	WordHandler     *WordHandler
	UserHandler     *UserHandler
	LearningHandler *LearningHandler
	HealthHandler   *HealthHandler
}

// NewHandlers 创建新的Handlers实例
func NewHandlers(
	wordHandler *WordHandler,
	userHandler *UserHandler,
	learningHandler *LearningHandler,
	healthHandler *HealthHandler,
) *Handlers {
	return &Handlers{
		WordHandler:     wordHandler,
		UserHandler:     userHandler,
		LearningHandler: learningHandler,
		HealthHandler:   healthHandler,
	}
}
