package service

import (
	"context"
	"errors"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
	domainErrors "github.com/lazyjean/sla2/internal/domain/errors" // Alias for domain errors
	"github.com/lazyjean/sla2/internal/domain/repository"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MemoryService 记忆服务接口
type MemoryService interface {
	// ReviewWord 复习单词
	ReviewWord(ctx context.Context, wordID entity.WordID, result bool, reviewDuration uint32) error
	// ReviewHanChar 复习汉字
	ReviewHanChar(ctx context.Context, hanCharID uint32, result bool, reviewDuration uint32) error
	// GetNextReviewWords 获取下一批需要复习的单词
	GetNextReviewWords(ctx context.Context, limit int) ([]*entity.Word, error)
	// GetWordStats 获取单词的学习统计信息
	GetWordStats(ctx context.Context, wordID entity.WordID) (*WordStats, error)
	// GetLearningProgress 获取学习进度
	GetLearningProgress(ctx context.Context) (*LearningProgress, error)
	// UpdateMemoryStatus 更新记忆单元状态
	UpdateMemoryStatus(ctx context.Context, memoryUnitID uint32, masteryLevel entity.MasteryLevel, studyDuration uint32) error
	// GetMemoryUnit 获取记忆单元
	GetMemoryUnit(ctx context.Context, id uint32) (*entity.MemoryUnit, error)
	// ListMemoriesForReview 获取需要复习的记忆单元列表
	ListMemoriesForReview(ctx context.Context, page, pageSize uint32, types []entity.MemoryUnitType, masteryLevels []entity.MasteryLevel) ([]*entity.MemoryUnit, int, error)
	// GetMemoryStats 获取记忆统计信息
	GetMemoryStats(ctx context.Context, unitType *entity.MemoryUnitType) (*repository.MemoryUnitStats, error)
	// CreateMemoryUnit 创建记忆单元
	CreateMemoryUnit(ctx context.Context, unit *entity.MemoryUnit) error
	// CreateMemoryUnits 批量创建记忆单元
	CreateMemoryUnits(ctx context.Context, units []*entity.MemoryUnit) ([]uint32, error)
	// UpdateMemoryUnit 更新记忆单元
	UpdateMemoryUnit(ctx context.Context, unit *entity.MemoryUnit) error
	// ListMemoriesForStrengthening 获取需要加强的记忆单元列表
	ListMemoriesForStrengthening(ctx context.Context, page, pageSize uint32, types []entity.MemoryUnitType, levels []entity.MasteryLevel) ([]*entity.MemoryUnit, int, error)
	// CalculateNextReviewInterval 计算下次复习间隔
	CalculateNextReviewInterval(unit *entity.MemoryUnit) *ReviewInterval
}

// WordStats 单词学习统计信息
type WordStats struct {
	WordID           entity.WordID       `json:"word_id"`
	Text             string              `json:"text"`
	MasteryLevel     entity.MasteryLevel `json:"mastery_level"`
	ReviewCount      uint32              `json:"review_count"`
	LastReviewAt     time.Time           `json:"last_review_at"`
	NextReviewAt     time.Time           `json:"next_review_at"`
	RetentionRate    float32             `json:"retention_rate"`
	ConsecutiveRight uint32              `json:"consecutive_right"`
	ConsecutiveWrong uint32              `json:"consecutive_wrong"`
}

// LearningProgress 学习进度
type LearningProgress struct {
	TotalWords      int     `json:"total_words"`
	MasteredWords   int     `json:"mastered_words"`
	LearningWords   int     `json:"learning_words"`
	NewWords        int     `json:"new_words"`
	MasteryRate     float64 `json:"mastery_rate"`
	RetentionRate   float64 `json:"retention_rate"`
	DailyReviewGoal int     `json:"daily_review_goal"`
	DailyNewGoal    int     `json:"daily_new_goal"`
}

// ReviewInterval 复习间隔
type ReviewInterval struct {
	Days    uint32
	Hours   uint32
	Minutes uint32
}

// MemoryServiceImpl 记忆服务实现
type MemoryServiceImpl struct {
	wordRepo    repository.WordRepository
	memoryRepo  repository.MemoryUnitRepository
	hanCharRepo repository.HanCharRepository // Add HanCharRepository
}

// NewMemoryService 创建记忆服务实例
func NewMemoryService(wordRepo repository.WordRepository, memoryRepo repository.MemoryUnitRepository, hanCharRepo repository.HanCharRepository) MemoryService { // Add hanCharRepo param
	return &MemoryServiceImpl{
		wordRepo:    wordRepo,
		memoryRepo:  memoryRepo,
		hanCharRepo: hanCharRepo, // Store hanCharRepo
	}
}

// ReviewWord 复习单词
func (s *MemoryServiceImpl) ReviewWord(ctx context.Context, wordID entity.WordID, result bool, reviewDuration uint32) error {
	log := logger.GetLogger(ctx)

	// Check if Word exists
	_, err := s.wordRepo.GetByID(ctx, wordID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("Attempted to review non-existent word", zap.Uint32("wordID", uint32(wordID)))
			return &domainErrors.Error{
				Code:    domainErrors.CodeWordNotFound,
				Message: "单词不存在",
			}
		}
		log.Error("Failed to check word existence", zap.Error(err), zap.Uint32("wordID", uint32(wordID)))
		return err // Return other DB errors
	}

	// Get UserID from context - Placeholder, needs interceptor fix
	userID, err := GetUserID(ctx)
	if err != nil {
		log.Error("Failed to get UserID in ReviewWord", zap.Error(err))
		// Returning error here because creating a unit without a user doesn't make sense
		return err
	}

	// 获取或创建记忆单元
	memoryUnit, err := s.memoryRepo.GetByTypeAndContentID(ctx, userID, entity.MemoryUnitTypeWord, uint32(wordID))
	if err != nil {
		log.Error("Failed to get memory unit by type and content ID", zap.Error(err), zap.Uint32("wordID", uint32(wordID)))
		return err
	}

	if memoryUnit == nil {
		memoryUnit = entity.NewMemoryUnit(entity.UID(userID), entity.MemoryUnitTypeWord, uint32(wordID))
		if err := s.memoryRepo.Create(ctx, memoryUnit); err != nil {
			log.Error("Failed to create new memory unit for word", zap.Error(err), zap.Uint32("wordID", uint32(wordID)), zap.Uint32("userID", uint32(userID)))
			return err
		}
		log.Info("Created new memory unit for word", zap.Uint32("unitID", uint32(memoryUnit.ID)), zap.Uint32("wordID", uint32(wordID)))
	}

	// 更新复习状态
	now := time.Now()
	memoryUnit.LastReviewAt = now
	memoryUnit.ReviewCount++
	if result {
		memoryUnit.ConsecutiveCorrect++
		memoryUnit.ConsecutiveWrong = 0
	} else {
		memoryUnit.ConsecutiveCorrect = 0
		memoryUnit.ConsecutiveWrong++
	}

	// 计算下次复习时间
	interval := s.calculateNextReviewInterval(memoryUnit)
	nextReviewAt := now.Add(time.Duration(interval.Days)*24*time.Hour +
		time.Duration(interval.Hours)*time.Hour +
		time.Duration(interval.Minutes)*time.Minute)
	memoryUnit.NextReviewAt = nextReviewAt

	// 添加日志
	log.Info("[Service ReviewWord] Calculated review interval",
		zap.Uint32("UnitID", uint32(memoryUnit.ID)),
		zap.Any("Interval", interval),
		zap.Time("LastReviewAt", memoryUnit.LastReviewAt),
		zap.Time("CalculatedNextReviewAt", nextReviewAt),
	)

	// 保存更新
	if err := s.memoryRepo.Update(ctx, memoryUnit); err != nil {
		log.Error("Failed to update memory unit after word review", zap.Error(err), zap.Uint32("unitID", uint32(memoryUnit.ID)))
		return err
	}
	return nil
}

// ReviewHanChar 复习汉字
func (s *MemoryServiceImpl) ReviewHanChar(ctx context.Context, hanCharID uint32, result bool, reviewDuration uint32) error {
	log := logger.GetLogger(ctx)

	// 1. Check if HanChar exists
	_, err := s.hanCharRepo.GetByID(ctx, entity.HanCharID(hanCharID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("Attempted to review non-existent HanChar", zap.Uint32("hanCharID", hanCharID))
			// Use a standard domain error or return gRPC status directly from handler
			// For now, returning a domain error that the gRPC handler should convert.
			return domainErrors.ErrNotFound // Return specific domain error
		}
		log.Error("Failed to check HanChar existence", zap.Error(err), zap.Uint32("hanCharID", hanCharID))
		return err // Return other DB errors
	}

	// 2. Get UserID from context - Placeholder, needs interceptor fix
	userID, err := GetUserID(ctx)
	if err != nil {
		log.Error("Failed to get UserID in ReviewHanChar", zap.Error(err))
		// Returning error here because creating a unit without a user doesn't make sense
		return err
	}

	// 3. Get or Create MemoryUnit (only if HanChar exists)
	memoryUnit, err := s.memoryRepo.GetByTypeAndContentID(ctx, userID, entity.MemoryUnitTypeHanChar, hanCharID)
	if err != nil {
		log.Error("Failed to get memory unit by type and content ID", zap.Error(err), zap.Uint32("hanCharID", hanCharID))
		return err
	}

	if memoryUnit == nil {
		memoryUnit = entity.NewMemoryUnit(userID, entity.MemoryUnitTypeHanChar, hanCharID)
		if err := s.memoryRepo.Create(ctx, memoryUnit); err != nil {
			log.Error("Failed to create new memory unit for han char", zap.Error(err), zap.Uint32("hanCharID", hanCharID), zap.Uint32("userID", uint32(userID)))
			return err
		}
		log.Info("Created new memory unit for han char", zap.Uint32("unitID", uint32(memoryUnit.ID)), zap.Uint32("hanCharID", hanCharID))
	}

	// 4. Update memory stats
	memoryUnit.LastReviewAt = time.Now()
	memoryUnit.ReviewCount++
	if result {
		memoryUnit.ConsecutiveCorrect++
		memoryUnit.ConsecutiveWrong = 0
	} else {
		memoryUnit.ConsecutiveCorrect = 0
		memoryUnit.ConsecutiveWrong++
	}

	// 5. Calculate next review time
	interval := s.calculateNextReviewInterval(memoryUnit)
	nextReviewAt := time.Now().Add(time.Duration(interval.Days)*24*time.Hour +
		time.Duration(interval.Hours)*time.Hour +
		time.Duration(interval.Minutes)*time.Minute)
	memoryUnit.NextReviewAt = nextReviewAt

	log.Info("[Service ReviewHanChar] Calculated review interval",
		zap.Uint32("UnitID", uint32(memoryUnit.ID)),
		zap.Any("Interval", interval),
		zap.Time("LastReviewAt", memoryUnit.LastReviewAt),
		zap.Time("CalculatedNextReviewAt", nextReviewAt),
	)

	// 6. Save updated memory unit
	if err := s.memoryRepo.Update(ctx, memoryUnit); err != nil {
		log.Error("Failed to update memory unit after han char review", zap.Error(err), zap.Uint32("unitID", uint32(memoryUnit.ID)))
		return err
	}
	return nil
}

// GetNextReviewWords 获取下一批需要复习的单词
func (s *MemoryServiceImpl) GetNextReviewWords(ctx context.Context, limit int) ([]*entity.Word, error) {
	// 从 context 获取 userID
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	// 获取需要复习的记忆单元
	units, err := s.memoryRepo.ListNeedReview(ctx, userID, entity.MemoryUnitTypeWord, time.Now(), limit)
	if err != nil {
		return nil, err
	}

	// 获取对应的单词
	var words []*entity.Word
	for _, unit := range units {
		word, err := s.wordRepo.GetByID(ctx, entity.WordID(unit.ContentID))
		if err != nil {
			continue
		}
		words = append(words, word)
	}

	return words, nil
}

// GetWordStats 获取单词的学习统计信息
func (s *MemoryServiceImpl) GetWordStats(ctx context.Context, wordID entity.WordID) (*WordStats, error) {
	// 从 context 获取 userID
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	// 获取单词
	word, err := s.wordRepo.GetByID(ctx, wordID)
	if err != nil {
		return nil, err
	}

	// 获取记忆单元
	unit, err := s.memoryRepo.GetByTypeAndContentID(ctx, userID, entity.MemoryUnitTypeWord, uint32(wordID))
	if err != nil {
		return nil, err
	}
	if unit == nil {
		unit = entity.NewMemoryUnit(userID, entity.MemoryUnitTypeWord, uint32(wordID))
	}

	return &WordStats{
		WordID:           wordID,
		Text:             word.Text,
		MasteryLevel:     unit.MasteryLevel,
		ReviewCount:      unit.ReviewCount,
		LastReviewAt:     unit.LastReviewAt,
		NextReviewAt:     unit.NextReviewAt,
		RetentionRate:    unit.RetentionRate,
		ConsecutiveRight: unit.ConsecutiveCorrect,
		ConsecutiveWrong: unit.ConsecutiveWrong,
	}, nil
}

// GetLearningProgress 获取学习进度
func (s *MemoryServiceImpl) GetLearningProgress(ctx context.Context) (*LearningProgress, error) {
	// TODO: 实现学习进度统计
	return &LearningProgress{}, nil
}

// UpdateMemoryStatus 更新记忆单元状态
func (s *MemoryServiceImpl) UpdateMemoryStatus(ctx context.Context, memoryUnitID uint32, masteryLevel entity.MasteryLevel, studyDuration uint32) error {
	// 从 context 获取 userID
	userID, err := GetUserID(ctx)
	if err != nil {
		return err
	}

	// 获取记忆单元
	memoryUnit, err := s.memoryRepo.GetByID(ctx, userID, memoryUnitID)
	if err != nil {
		return err
	}
	if memoryUnit == nil {
		return domainErrors.ErrNotFound
	}

	// 更新状态
	memoryUnit.MasteryLevel = masteryLevel
	memoryUnit.StudyDuration += studyDuration
	memoryUnit.Update()

	// 保存更新
	return s.memoryRepo.Update(ctx, memoryUnit)
}

// GetMemoryUnit 获取记忆单元
func (s *MemoryServiceImpl) GetMemoryUnit(ctx context.Context, memoryUnitID uint32) (*entity.MemoryUnit, error) {
	// 从 context 获取 userID
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	return s.memoryRepo.GetByID(ctx, userID, memoryUnitID)
}

// ListMemoriesForReview 获取需要复习的记忆单元列表
func (s *MemoryServiceImpl) ListMemoriesForReview(ctx context.Context, page, pageSize uint32, types []entity.MemoryUnitType, masteryLevels []entity.MasteryLevel) ([]*entity.MemoryUnit, int, error) {
	// 从 context 获取 userID
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, 0, err
	}

	// 计算 offset
	offset := (page - 1) * pageSize
	limit := int(pageSize)
	now := time.Now()

	// 调用 repository 获取数据
	units, err := s.memoryRepo.ListNeedReviewByTypes(ctx, userID, types, now, offset, limit, masteryLevels)
	if err != nil {
		return nil, 0, err
	}

	// 调用 repository 获取总数
	total, err := s.memoryRepo.CountNeedReviewByTypes(ctx, userID, types, now, masteryLevels)
	if err != nil {
		return nil, 0, err
	}

	return units, int(total), nil
}

// calculateNextReviewInterval 计算下次复习间隔
func (s *MemoryServiceImpl) calculateNextReviewInterval(unit *entity.MemoryUnit) *ReviewInterval {
	// 基础间隔（小时）
	var baseInterval float64
	switch unit.MasteryLevel {
	case entity.MasteryLevelUnlearned:
		baseInterval = 1 // 1小时
	case entity.MasteryLevelBeginner:
		baseInterval = 4 // 4小时
	case entity.MasteryLevelFamiliar:
		baseInterval = 24 // 1天
	case entity.MasteryLevelMastered:
		baseInterval = 72 // 3天
	case entity.MasteryLevelExpert:
		baseInterval = 168 // 7天
	default:
		baseInterval = 1
	}

	// 根据连续正确次数调整间隔
	if unit.ConsecutiveCorrect > 0 {
		baseInterval *= 1.2 * float64(unit.ConsecutiveCorrect)
	}

	// 根据连续错误次数减少间隔
	if unit.ConsecutiveWrong > 0 {
		baseInterval /= 2.0 * float64(unit.ConsecutiveWrong)
	}

	// 转换为天、小时、分钟
	totalMinutes := int(baseInterval * 60)
	days := uint32(totalMinutes / (24 * 60))
	totalMinutes = totalMinutes % (24 * 60)
	hours := uint32(totalMinutes / 60)
	minutes := uint32(totalMinutes % 60)

	return &ReviewInterval{
		Days:    days,
		Hours:   hours,
		Minutes: minutes,
	}
}

// GetMemoryStats 获取记忆统计信息
func (s *MemoryServiceImpl) GetMemoryStats(ctx context.Context, unitType *entity.MemoryUnitType) (*repository.MemoryUnitStats, error) {
	log := logger.GetLogger(ctx)

	// 从 context 获取 userID
	userID, err := GetUserID(ctx) // Use the function from the same package
	if err != nil {
		log.Error("Failed to get userID from context for GetMemoryStats", zap.Error(err))
		// 根据实际错误类型返回，例如 ErrNoUserInContext 可能映射为 Unauthenticated
		return nil, err // Or: status.Error(codes.Unauthenticated, "user not found in context")
	}

	// 如果未指定类型，我们需要决定是返回错误、所有类型的聚合，还是某种默认类型。
	// 当前 repository.GetStats 只接受一个明确的类型。
	// 暂时：如果 unitType 为 nil，我们返回错误或默认类型（例如 HanChar）。
	// TODO: 明确处理 unitType == nil 的情况。
	if unitType == nil {
		log.Warn("GetMemoryStats called without specific unit type, using default", zap.Uint32("userID", uint32(userID)))
		// 选项 B: 使用默认类型 (例如 HanChar)
		defaultType := entity.MemoryUnitTypeHanChar
		unitType = &defaultType
	}

	// 直接使用 entity.UID 类型调用 repository
	stats, err := s.memoryRepo.GetStats(ctx, userID, *unitType)
	if err != nil {
		log.Error("Failed to get stats from repository", zap.Error(err), zap.Uint32("userID", uint32(userID)), zap.Any("unitType", *unitType))
		return nil, err
	}
	return stats, nil
}

// CreateMemoryUnit 创建记忆单元
func (s *MemoryServiceImpl) CreateMemoryUnit(ctx context.Context, unit *entity.MemoryUnit) error {
	return s.memoryRepo.Create(ctx, unit)
}

// CreateMemoryUnits 批量创建记忆单元
func (s *MemoryServiceImpl) CreateMemoryUnits(ctx context.Context, units []*entity.MemoryUnit) ([]uint32, error) {
	// 使用事务批量创建
	err := s.memoryRepo.CreateBatch(ctx, units)
	if err != nil {
		return nil, err
	}

	// 收集创建的单元ID
	ids := make([]uint32, len(units))
	for i, unit := range units {
		ids[i] = uint32(unit.ID)
	}
	return ids, nil
}

// UpdateMemoryUnit 更新记忆单元
func (s *MemoryServiceImpl) UpdateMemoryUnit(ctx context.Context, unit *entity.MemoryUnit) error {
	return s.memoryRepo.Update(ctx, unit)
}

// ListMemoriesForStrengthening 获取需要加强的记忆单元列表
func (s *MemoryServiceImpl) ListMemoriesForStrengthening(ctx context.Context, page, pageSize uint32, types []entity.MemoryUnitType, levels []entity.MasteryLevel) ([]*entity.MemoryUnit, int, error) {
	// 从 context 获取 userID
	userID, err := GetUserID(ctx)
	if err != nil {
		return nil, 0, err
	}

	// 计算 offset
	offset := (page - 1) * pageSize
	limit := int(pageSize)
	now := time.Now()

	// 调用 repository 获取数据
	units, err := s.memoryRepo.ListNeedReviewByTypes(ctx, userID, types, now, offset, limit, levels)
	if err != nil {
		return nil, 0, err
	}

	// 调用 repository 获取总数
	total, err := s.memoryRepo.CountNeedReviewByTypes(ctx, userID, types, now, levels)
	if err != nil {
		return nil, 0, err
	}

	return units, int(total), nil
}

// CalculateNextReviewInterval 计算下次复习间隔
func (s *MemoryServiceImpl) CalculateNextReviewInterval(unit *entity.MemoryUnit) *ReviewInterval {
	// 基础间隔（小时）
	var baseInterval float64
	switch unit.MasteryLevel {
	case entity.MasteryLevelUnlearned:
		baseInterval = 1 // 1小时
	case entity.MasteryLevelBeginner:
		baseInterval = 4 // 4小时
	case entity.MasteryLevelFamiliar:
		baseInterval = 24 // 1天
	case entity.MasteryLevelMastered:
		baseInterval = 72 // 3天
	case entity.MasteryLevelExpert:
		baseInterval = 168 // 7天
	default:
		baseInterval = 1
	}

	// 根据连续正确次数调整间隔
	if unit.ConsecutiveCorrect > 0 {
		baseInterval *= 1.2 * float64(unit.ConsecutiveCorrect)
	}

	// 根据连续错误次数减少间隔
	if unit.ConsecutiveWrong > 0 {
		baseInterval /= 2.0 * float64(unit.ConsecutiveWrong)
	}

	// 转换为天、小时、分钟
	totalMinutes := int(baseInterval * 60)
	days := uint32(totalMinutes / (24 * 60))
	totalMinutes = totalMinutes % (24 * 60)
	hours := uint32(totalMinutes / 60)
	minutes := uint32(totalMinutes % 60)

	return &ReviewInterval{
		Days:    days,
		Hours:   hours,
		Minutes: minutes,
	}
}

var _ MemoryService = (*MemoryServiceImpl)(nil)
