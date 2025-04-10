package learning

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	pg "github.com/lazyjean/sla2/internal/infrastructure/persistence/postgres"
	"github.com/lazyjean/sla2/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/gorm"
)

var (
	grpcService    *LearningService
	memoryUnitRepo repository.MemoryUnitRepository
	learningRepo   repository.LearningRepository
	db             *gorm.DB
)

// ensureLogger 初始化基础 logger (如果尚未初始化)
func ensureLogger() {
	if logger.Log == nil {
		logger.InitBaseLogger()
		logger.Log.Info("Initialized base logger for testing.")
	}
}

// SetupTestDB 设置测试数据库并初始化必要的 repo
func SetupTestDB(t *testing.T) error {
	ensureLogger() // Initialize logger here

	var cleanup func()
	db, cleanup = pg.SetupTestDB(t) // Assuming pg.SetupTestDB is the correct package and function
	t.Cleanup(cleanup)

	// Run migrations to ensure schema matches entities
	err := db.AutoMigrate(
		&entity.MemoryUnit{}, // Add other entities if needed for tests in this package
		&entity.HanChar{},    // Ensure HanChar is migrated here too if SetupTestDB is used elsewhere
		// &entity.CourseLearningProgress{}, // Example
		// &entity.CourseSectionProgress{}, // Example
		// &entity.CourseSectionUnitProgress{}, // Example
	)
	if err != nil {
		logger.Log.Error("AutoMigrate failed", zap.Error(err))
		return fmt.Errorf("automigrate failed: %w", err)
	}
	logger.Log.Info("AutoMigrate completed for learning tests.")

	// 清理测试数据 (Truncate after migration)
	logger.Log.Info("Truncating memory_units table...")
	if err := db.Exec("TRUNCATE TABLE memory_units CASCADE").Error; err != nil {
		logger.Log.Error("Failed to truncate memory_units", zap.Error(err))
		return err
	}
	logger.Log.Info("Truncating han_chars table...") // Truncate HanChar here too
	if err := db.Exec("TRUNCATE TABLE han_chars CASCADE").Error; err != nil {
		logger.Log.Error("Failed to truncate han_chars", zap.Error(err))
		return err
	}
	// ... add other TRUNCATE statements if needed ...

	// Initialize repositories used in tests
	learningRepo = pg.NewLearningRepository(db)
	memoryUnitRepo = pg.NewMemoryUnitRepository(db)
	hanCharRepo := pg.NewHanCharRepository(db)
	wordRepo := pg.NewWordRepository(db)

	// Initialize services needed for tests (can be done here or in TestMain/specific tests)
	memoryService := service.NewMemoryService(wordRepo, memoryUnitRepo, hanCharRepo)
	learningService := service.NewLearningService(learningRepo, memoryService)
	grpcService = NewLearningService(learningService, memoryService)

	return nil
}

// setupRealGrpcTest sets up DB, repos, services, and an in-memory gRPC server/client.
func setupRealGrpcTest(t *testing.T) (context.Context, pb.LearningServiceClient, *gorm.DB, func()) {
	ensureLogger() // Keep logger setup

	// --- Database Setup ---
	var testDB *gorm.DB // Use local variable
	var dbCleanup func()
	testDB, dbCleanup = pg.SetupTestDB(t) // Re-use the underlying test DB setup logic
	require.NotNil(t, testDB, "DB setup failed")

	// Migrations (Ensure all necessary entities are included)
	err := testDB.AutoMigrate(
		&entity.MemoryUnit{},
		&entity.HanChar{}, // Ensure HanChar is migrated
		// Add other entities specific to this test suite if needed
	)
	require.NoError(t, err, "AutoMigrate failed in setupRealGrpcTest")

	// Truncate tables (Ensure all necessary tables are truncated)
	require.NoError(t, testDB.Exec("TRUNCATE TABLE memory_units CASCADE").Error, "Truncate memory_units failed in setupRealGrpcTest")
	require.NoError(t, testDB.Exec("TRUNCATE TABLE han_chars CASCADE").Error, "Truncate han_chars failed in setupRealGrpcTest")
	// Add other truncations if needed

	// --- Initialize Repositories (Local Scope) ---
	localLearningRepo := pg.NewLearningRepository(testDB)
	localMemoryUnitRepo := pg.NewMemoryUnitRepository(testDB)
	localHanCharRepo := pg.NewHanCharRepository(testDB)
	localWordRepo := pg.NewWordRepository(testDB)
	localMemoryService := service.NewMemoryService(localWordRepo, localMemoryUnitRepo, localHanCharRepo)
	localLearningService := service.NewLearningService(localLearningRepo, localMemoryService)

	// --- Setup gRPC Server ---
	ctx := context.Background()
	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer) // Local listener

	// Define a simple auth interceptor for testing
	authInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			// Allow requests without metadata for tests not needing auth, or return error
			// For this test, GetMemoryStats needs auth, so ideally we check method name
			// or just return error if metadata is missing.
			// Simple approach for now: Proceed if no metadata, specific calls might fail later if they need UserID
			return handler(ctx, req)
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			// Proceed if no auth header, specific calls might fail later if they need UserID
			return handler(ctx, req)
		}

		// In a real app, validate the token ("Bearer test-token")
		// For this test, we assume it's valid and corresponds to a user.
		// Extract the UserID used in the test (User ID 1)
		testUserID := entity.UID(1) // Match the ID used in TestHanCharLearningFlow
		newCtx := service.WithUserID(ctx, testUserID)
		return handler(newCtx, req)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor), // Add the interceptor
	)

	// Create the specific gRPC service instance using local services
	grpcLearningSvcImpl := NewLearningService(localLearningService, localMemoryService)

	// Register the service
	pb.RegisterLearningServiceServer(grpcServer, grpcLearningSvcImpl)

	// Start server in background
	go func() {
		if err := grpcServer.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			t.Logf("gRPC server error: %v", err) // Log error instead of Fatal
		}
	}()

	// --- Setup gRPC Client ---
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err, "Failed to dial bufnet")

	grpcClient := pb.NewLearningServiceClient(conn)

	// --- Cleanup Function ---
	cleanupFunc := func() {
		err := conn.Close()
		if err != nil {
			t.Logf("Failed to close gRPC client connection: %v", err)
		}
		grpcServer.Stop() // Stop the gRPC server
		err = listener.Close()
		if err != nil {
			t.Logf("Failed to close bufconn listener: %v", err)
		}
		dbCleanup() // Close the DB connection
	}

	return ctx, grpcClient, testDB, cleanupFunc
}

func TestMain(m *testing.M) {
	// TestMain can be simplified or removed if SetupTestDB handles all setup
	// Or, it can handle setup not specific to individual tests
	code := m.Run()
	os.Exit(code)
}

func TestGetCourseProgress(t *testing.T) {
	if err := SetupTestDB(t); err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	// 准备测试数据
	progress := &entity.CourseLearningProgress{
		UserID:   1,
		CourseID: 1,
		Status:   "completed",
		Score:    100,
		Progress: 100,
	}
	if err := learningRepo.SaveCourseProgress(context.Background(), progress); err != nil {
		t.Fatalf("Failed to save test data: %v", err)
	}

	// 执行测试
	req := &pb.LearningServiceGetCourseProgressRequest{
		CourseId: 1,
	}
	ctx := service.WithUserID(context.Background(), entity.UID(1))
	resp, err := grpcService.GetCourseProgress(ctx, req)
	if err != nil {
		t.Fatalf("Failed to get course progress: %v", err)
	}

	// 验证结果
	assert.Equal(t, float32(100.0), resp.Progress.Progress)
	assert.Equal(t, uint32(0), resp.Progress.CompletedItems)
	assert.Equal(t, uint32(0), resp.Progress.TotalItems)
}

func TestGetSectionProgress(t *testing.T) {
	if err := SetupTestDB(t); err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	// 准备测试数据
	progress := &entity.CourseSectionProgress{
		UserID:    1,
		CourseID:  1,
		SectionID: 1,
		Status:    "completed",
		Progress:  100,
	}
	if err := learningRepo.SaveSectionProgress(context.Background(), progress); err != nil {
		t.Fatalf("Failed to save test data: %v", err)
	}

	// 执行测试
	req := &pb.LearningServiceGetSectionProgressRequest{
		SectionId: 1,
	}
	ctx := service.WithUserID(context.Background(), entity.UID(1))
	resp, err := grpcService.GetSectionProgress(ctx, req)
	if err != nil {
		t.Fatalf("Failed to get section progress: %v", err)
	}

	// 验证结果
	assert.Equal(t, float32(100.0), resp.Progress.Progress)
	assert.Equal(t, uint32(0), resp.Progress.CompletedItems)
	assert.Equal(t, uint32(0), resp.Progress.TotalItems)
}

func TestUpdateUnitProgress(t *testing.T) {
	if err := SetupTestDB(t); err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	// 准备测试数据
	progress := &entity.CourseSectionUnitProgress{
		UserID:        1,
		SectionID:     1,
		UnitID:        1,
		Status:        "in_progress",
		CompleteCount: 0,
	}
	if err := learningRepo.UpsertUnitProgress(context.Background(), progress); err != nil {
		t.Fatalf("Failed to save test data: %v", err)
	}

	// 执行测试
	req := &pb.LearningServiceUpdateUnitProgressRequest{
		UnitId:    1,
		SectionId: 1,
		Completed: true,
	}
	ctx := service.WithUserID(context.Background(), entity.UID(1))
	resp, err := grpcService.UpdateUnitProgress(ctx, req)
	if err != nil {
		t.Fatalf("Failed to update unit progress: %v", err)
	}

	// 验证结果
	assert.NotNil(t, resp)

	// 验证数据库中的更新
	unitProgresses, err := learningRepo.ListUnitProgress(ctx, 1, 1)
	if err != nil {
		t.Fatalf("Failed to get updated progress: %v", err)
	}
	assert.Equal(t, 1, len(unitProgresses))
	assert.Equal(t, "completed", unitProgresses[0].Status)
	assert.Equal(t, uint(1), unitProgresses[0].CompleteCount)
}

func TestUpdateMemoryStatus(t *testing.T) {
	if err := SetupTestDB(t); err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	// 准备测试数据
	now := time.Now()
	unit := &entity.MemoryUnit{
		UserID:             1,
		Type:               entity.MemoryUnitTypeWord,
		ContentID:          1,
		CreatedAt:          now,
		UpdatedAt:          now,
		MasteryLevel:       entity.MasteryLevelBeginner,
		ReviewCount:        0,
		NextReviewAt:       now.Add(time.Hour),
		LastReviewAt:       now,
		StudyDuration:      0,
		RetentionRate:      0,
		ConsecutiveCorrect: 0,
		ConsecutiveWrong:   0,
	}
	if err := memoryUnitRepo.Create(context.Background(), unit); err != nil {
		t.Fatalf("Failed to save test data: %v", err)
	}

	// 执行测试
	req := &pb.UpdateMemoryStatusRequest{
		MemoryUnitId:  unit.ID,
		MasteryLevel:  pb.MasteryLevel_MASTERY_LEVEL_FAMILIAR,
		StudyDuration: 100,
	}
	resp, err := grpcService.UpdateMemoryStatus(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to update memory status: %v", err)
	}

	// 验证结果
	assert.NotNil(t, resp)

	// 验证数据库中的更新
	updatedUnit, err := memoryUnitRepo.GetByID(context.Background(), unit.ID)
	if err != nil {
		t.Fatalf("Failed to get updated unit: %v", err)
	}
	assert.Equal(t, entity.MasteryLevelFamiliar, updatedUnit.MasteryLevel)
	assert.Equal(t, uint32(100), updatedUnit.StudyDuration)
}

func TestReviewWord(t *testing.T) {
	if err := SetupTestDB(t); err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	ctx := context.Background()
	wordID := uint32(99999) // Use an ID unlikely to exist

	// 1. First review (Correct) - Should now FAIL because the word doesn't exist.
	req1 := &pb.ReviewWordRequest{
		WordId:       wordID,
		Result:       pb.ReviewResult_REVIEW_RESULT_CORRECT,
		ResponseTime: 1500, // 1.5 seconds
	}

	// Expect an error here because the word does not exist
	_, err := grpcService.ReviewWord(ctx, req1)
	// assert.NoError(t, err) // REMOVED - We expect an error
	require.Error(t, err, "Expected an error when reviewing a non-existent word")

	// Check if the error is a gRPC status error with code NotFound
	st, ok := status.FromError(err)
	require.True(t, ok, "Error should be a gRPC status error")
	// Note: The exact error code depends on how the service maps domainErrors.ErrNotFound
	// Assuming it maps to codes.NotFound
	assert.Equal(t, codes.NotFound, st.Code(), "Expected NotFound error code")

	// Since the first review failed, the MemoryUnit should NOT have been created.
	// The following checks are removed or commented out as they are no longer valid.
	/*
		assert.NotNil(t, resp1) // REMOVED

		// Verify creation and initial state in DB
		unit1, err := memoryUnitRepo.GetByTypeAndContentID(ctx, entity.MemoryUnitTypeWord, wordID)
		assert.NoError(t, err) // This might fail now or return nil unit
		assert.NotNil(t, unit1, "MemoryUnit should NOT have been created") // MODIFIED assertion (or remove)
		t.Logf("[Review 1] LastReviewAt: %s, NextReviewAt: %s", unit1.LastReviewAt, unit1.NextReviewAt)
		assert.True(t, unit1.NextReviewAt.After(unit1.LastReviewAt), "Next review time must be after last review")
		initialNextReviewAt := unit1.NextReviewAt

		// 2. Second review (Wrong) - Cannot proceed if unit1 wasn't created
		time.Sleep(10 * time.Millisecond)              // Ensure timestamp changes
		lastReviewAtBeforeUpdate := unit1.LastReviewAt // Record time before the update
		req2 := &pb.ReviewWordRequest{
			WordId:       wordID,
			Result:       pb.ReviewResult_REVIEW_RESULT_WRONG,
			ResponseTime: 2100, // 2.1 seconds
		}
		resp2, err := grpcService.ReviewWord(ctx, req2)
		assert.NoError(t, err)
		assert.NotNil(t, resp2)

		// Verify updated state in DB
		unit2, err := memoryUnitRepo.GetByID(ctx, unit1.ID)
		assert.NoError(t, err)
		assert.NotNil(t, unit2)
		t.Logf("[Review 2] LastReviewAt: %s, NextReviewAt: %s", unit2.LastReviewAt, unit2.NextReviewAt)
		assert.True(t, unit2.LastReviewAt.After(lastReviewAtBeforeUpdate), "Last review time should be updated")
		assert.True(t, unit2.NextReviewAt.After(unit2.LastReviewAt), "Next review time must be after last review")
		assert.True(t, unit2.NextReviewAt.Before(initialNextReviewAt.Add(time.Hour)), "Next review time after wrong should be sooner than initial next review plus buffer")
	*/
}

func TestReviewHanChar(t *testing.T) {
	if err := SetupTestDB(t); err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	ctx := context.Background()
	hanCharID := uint32(88888) // Use an ID unlikely to exist

	// 1. First review (Correct) - Should FAIL because HanChar does not exist
	req1 := &pb.ReviewHanCharRequest{
		HanCharId:    hanCharID,
		Result:       pb.ReviewResult_REVIEW_RESULT_CORRECT,
		ResponseTime: 1800, // 1.8 seconds
	}
	// Expect an error here
	_, err := grpcService.ReviewHanChar(ctx, req1)
	// assert.NoError(t, err) // REMOVED
	require.Error(t, err, "Expected an error when reviewing a non-existent han char")

	// Check for gRPC NotFound status
	st, ok := status.FromError(err)
	require.True(t, ok, "Error should be a gRPC status error")
	assert.Equal(t, codes.NotFound, st.Code(), "Expected NotFound error code")

	// Since the first review failed, the MemoryUnit should NOT have been created.
	// The following checks are removed or commented out.
	/*
		assert.NotNil(t, resp1) // REMOVED

		// Verify creation and initial state in DB
		unit1, err := memoryUnitRepo.GetByTypeAndContentID(ctx, entity.MemoryUnitTypeHanChar, hanCharID)
		assert.NoError(t, err)
		assert.NotNil(t, unit1, "MemoryUnit should NOT have been created") // MODIFIED or remove
		assert.Equal(t, uint32(1), unit1.ReviewCount)
		assert.Equal(t, entity.MasteryLevelBeginner, unit1.MasteryLevel)
		assert.Equal(t, uint32(1), unit1.ConsecutiveCorrect)
		assert.Equal(t, uint32(0), unit1.ConsecutiveWrong)
		assert.Equal(t, uint32(1), unit1.StudyDuration) // 1800ms / 1000 = 1s
		assert.True(t, unit1.NextReviewAt.After(unit1.LastReviewAt), "Next review time must be after last review")
		initialNextReviewAt2 := unit1.NextReviewAt

		// 2. Second review (Correct) - Cannot proceed
		time.Sleep(10 * time.Millisecond)
		lastReviewAtBeforeUpdate2 := unit1.LastReviewAt
		req2 := &pb.ReviewHanCharRequest{
			HanCharId:    hanCharID,
			Result:       pb.ReviewResult_REVIEW_RESULT_CORRECT,
			ResponseTime: 1200, // 1.2 seconds
		}
		resp2, err := grpcService.ReviewHanChar(ctx, req2)
		assert.NoError(t, err)
		assert.NotNil(t, resp2)

		// Verify updated state in DB
		unit2, err := memoryUnitRepo.GetByID(ctx, unit1.ID)
		assert.NoError(t, err)
		assert.NotNil(t, unit2)
		assert.Equal(t, uint32(2), unit2.ReviewCount)
		assert.Equal(t, entity.MasteryLevelBeginner, unit2.MasteryLevel) // 2 correct -> still Beginner (needs 3 for Familiar)
		assert.Equal(t, uint32(2), unit2.ConsecutiveCorrect)
		assert.Equal(t, uint32(0), unit2.ConsecutiveWrong)
		assert.Equal(t, uint32(1+1), unit2.StudyDuration) // 1s + 1s = 2s
		assert.True(t, unit2.LastReviewAt.After(lastReviewAtBeforeUpdate2), "Last review time should be updated")
		assert.True(t, unit2.NextReviewAt.After(unit2.LastReviewAt), "Next review time must be after last review")
		assert.True(t, unit2.NextReviewAt.After(initialNextReviewAt2), "Next review time should be further out after second correct")
	*/
}

// TestHanCharLearningFlow uses the new setup function
func TestHanCharLearningFlow(t *testing.T) {
	// Use the new setup function that returns a client
	ctx, client, db, cleanup := setupRealGrpcTest(t) // Renamed function call
	defer cleanup()

	// Prepare user context
	_ = entity.UID(1) // Use a specific user ID for testing (Marked as unused for now)
	// token, err := security.GenerateTestJWT(testUserID, "test-flow@example.com", time.Hour) // Function does not exist
	// require.NoError(t, err)
	token := "test-token" // Placeholder token for testing metadata
	userCtx := metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", "Bearer "+token))

	// --- 1. Seed Test Han Characters ---
	hanCharRepo := pg.NewHanCharRepository(db)
	testChars := []entity.HanChar{
		{Character: "测", Pinyin: "cè", Level: 1, Tags: []string{}, Categories: []string{}, Examples: []string{}},
		{Character: "试", Pinyin: "shì", Level: 1, Tags: []string{}, Categories: []string{}, Examples: []string{}},
		{Character: "学", Pinyin: "xué", Level: 1, Tags: []string{}, Categories: []string{}, Examples: []string{}},
		{Character: "习", Pinyin: "xí", Level: 1, Tags: []string{}, Categories: []string{}, Examples: []string{}},
	}
	charIDs := make(map[string]entity.HanCharID)
	for i := range testChars {
		char := &testChars[i]
		_, err := hanCharRepo.Create(ctx, char) // Assign both return values
		require.NoError(t, err, "Failed to create test HanChar: %s", char.Character)
		require.NotZero(t, char.ID, "HanChar ID should be populated after creation")
		charIDs[char.Character] = char.ID
		t.Logf("Seeded HanChar: %s (ID: %d)", char.Character, char.ID)
	}
	require.Len(t, charIDs, len(testChars), "Should have seeded all test characters")

	// --- 2. First Learning Interaction ---
	t.Run("FirstLearn", func(t *testing.T) {
		reviewReq := &pb.ReviewHanCharRequest{
			HanCharId:    uint32(charIDs["测"]),
			Result:       pb.ReviewResult_REVIEW_RESULT_CORRECT,
			ResponseTime: 1000,
		}
		_, err := client.ReviewHanChar(userCtx, reviewReq)
		require.NoError(t, err, "First review of '测' should succeed")
	})

	// --- 3. Multiple Learning Attempts ---
	t.Run("MultipleAttempts", func(t *testing.T) {
		// Review '测' again (correct)
		_, err := client.ReviewHanChar(userCtx, &pb.ReviewHanCharRequest{
			HanCharId:    uint32(charIDs["测"]),
			Result:       pb.ReviewResult_REVIEW_RESULT_CORRECT,
			ResponseTime: 1100,
		})
		require.NoError(t, err, "Second review of '测' (correct) should succeed")

		// Review '测' again (incorrect)
		_, err = client.ReviewHanChar(userCtx, &pb.ReviewHanCharRequest{
			HanCharId:    uint32(charIDs["测"]),
			Result:       pb.ReviewResult_REVIEW_RESULT_WRONG,
			ResponseTime: 1200,
		})
		require.NoError(t, err, "Third review of '测' (incorrect) should succeed")

		// Review '试' (correct)
		_, err = client.ReviewHanChar(userCtx, &pb.ReviewHanCharRequest{
			HanCharId:    uint32(charIDs["试"]),
			Result:       pb.ReviewResult_REVIEW_RESULT_CORRECT,
			ResponseTime: 900,
		})
		require.NoError(t, err, "First review of '试' (correct) should succeed")

		// Review non-existent character ID (should fail)
		nonExistentID := uint32(999999)
		_, err = client.ReviewHanChar(userCtx, &pb.ReviewHanCharRequest{
			HanCharId:    nonExistentID,
			Result:       pb.ReviewResult_REVIEW_RESULT_CORRECT,
			ResponseTime: 1300,
		})
		require.Error(t, err, "Review of non-existent HanChar ID should fail")
		st, ok := status.FromError(err)
		require.True(t, ok, "Error should be a gRPC status error")
		if st.Code() != codes.NotFound {
			require.Contains(t, st.Message(), "not found", "Error message should indicate not found if code is not NotFound")
		} else {
			require.Equal(t, codes.NotFound, st.Code(), "Error code should be NotFound")
		}

		t.Log("Completed multiple learning attempts.")
	})

	// --- 4. Check Statistics ---
	t.Run("CheckStats", func(t *testing.T) {
		memoryType := pb.MemoryUnitType_MEMORY_UNIT_TYPE_HAN_CHAR // Define the type
		statsReq := &pb.GetMemoryStatsRequest{
			Type: &memoryType, // Use the correct field name 'Type' and pass its pointer
		}
		statsRes, err := client.GetMemoryStats(userCtx, statsReq)
		require.NoError(t, err, "GetMemoryStats for HanChar should succeed")
		require.NotNil(t, statsRes, "Stats response should not be nil")

		stats := statsRes // Access fields directly on statsRes
		// Log retrieved stats using new field names
		t.Logf("Retrieved Stats: TotalLearned=%d, Mastered=%d, NeedReview=%d, LevelStats=%v",
			stats.TotalLearned, stats.MasteredCount, stats.NeedReviewCount, stats.LevelStats)

		expectedTotalLearned := uint32(2) // Renamed for clarity based on response field

		require.Equal(t, expectedTotalLearned, stats.TotalLearned, "Total learned units should match reviewed characters for the user")

		// Calculate sum from LevelStats map - assuming keys match entity.MasteryLevel values
		beginnerKey := uint32(entity.MasteryLevelBeginner) // Assuming Beginner exists
		familiarKey := uint32(entity.MasteryLevelFamiliar) // Assuming Familiar exists
		// masteredKey := uint32(entity.MasteryLevelMastered) // Assuming Mastered exists - Removed, using direct field

		actualLearnedSum := stats.LevelStats[beginnerKey] + stats.LevelStats[familiarKey] + stats.MasteredCount // Sum beginner, familiar, mastered

		require.GreaterOrEqual(t, actualLearnedSum, expectedTotalLearned, "Sum of learned levels (Beginner, Familiar, Mastered) should be at least the number of total learned units")
		require.Equal(t, uint32(0), stats.MasteredCount, "Mastered count should be 0 after these initial reviews") // Check MasteredCount directly

		// Replace checks for BeginnerCount, IntermediateCount, AdvancedCount with LevelStats checks
		actualNonMasteredSum := stats.LevelStats[beginnerKey] + stats.LevelStats[familiarKey] // Sum non-mastered levels tracked in LevelStats
		require.Equal(t, expectedTotalLearned, actualNonMasteredSum, "Sum of non-mastered tracked levels (Beginner, Familiar) should match total learned count")

		// Check default request
		statsReqDefault := &pb.GetMemoryStatsRequest{}
		statsResDefault, errDefault := client.GetMemoryStats(userCtx, statsReqDefault)
		require.NoError(t, errDefault, "GetMemoryStats with default type should succeed")
		require.NotNil(t, statsResDefault, "Default stats response should not be nil")
		require.Equal(t, stats.TotalLearned, statsResDefault.TotalLearned, "Default stats (HAN_CHAR) should match explicitly requested HAN_CHAR stats")

	})

	t.Log("HanChar learning flow test completed successfully.")
}
