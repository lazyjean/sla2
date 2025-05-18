package learning_test

import (
	"context"
	"net"
	"testing"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/config"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/infrastructure/test"
	"github.com/lazyjean/sla2/internal/wire"
	"github.com/lazyjean/sla2/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
)

func setupTest(t *testing.T) (context.Context, pb.LearningServiceClient, pb.VocabularyServiceClient, func()) {
	// 1. Setup
	app, err := wire.InitializeTestApp(t)
	config.GetConfig().Log.Level = "warn"
	config.GetConfig().Log.Format = "console"
	// config.GetConfig().Database.Debug = true
	require.NoError(t, err)

	// 2. Start app
	ctx := logger.WithContext(context.Background(), logger.NewAppLogger(&config.GetConfig().Log))
	go func() {
		t.Logf("Starting app for test: %s", t.Name())
		err = app.Start(ctx)
		require.NoError(t, err)
	}()

	// 3. Create client
	lis := app.GetGRPCListener().GetGRPCListener().(*bufconn.Listener)
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	require.NoError(t, err)

	client := pb.NewLearningServiceClient(conn)
	vocabClient := pb.NewVocabularyServiceClient(conn)

	// Add mock token to context
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer mock-token")

	// Initialize test data
	db := test.GetTestDB()
	err = db.Exec(`
		TRUNCATE TABLE memory_units CASCADE;
		TRUNCATE TABLE users CASCADE;
		TRUNCATE TABLE courses CASCADE;
		TRUNCATE TABLE course_sections CASCADE;
		TRUNCATE TABLE course_section_units CASCADE;
		TRUNCATE TABLE course_learning_progresses CASCADE;
		TRUNCATE TABLE course_section_progresses CASCADE;
		TRUNCATE TABLE han_chars CASCADE;
		INSERT INTO users (id, username, email, created_at, updated_at) VALUES (1, 'test_user', 'test@example.com', NOW(), NOW());
		INSERT INTO courses (id, title, description, level, created_at, updated_at) VALUES (1, 'Test Course', 'Test Description', 'BEGINNER', NOW(), NOW());
		INSERT INTO course_sections (id, course_id, title, "desc", created_at, updated_at) VALUES (1, 1, 'Test Section', 'Test Description', NOW(), NOW());
		INSERT INTO course_section_units (id, section_id, title, "desc", created_at, updated_at) VALUES (1, 1, 'Test Unit', 'Test Content', NOW(), NOW());
		INSERT INTO course_learning_progresses (user_id, course_id, status, score, progress, created_at, updated_at) 
		VALUES (1, 1, 'not_started', 0, 0.0, NOW(), NOW());
		INSERT INTO course_section_progresses (user_id, course_id, section_id, status, progress, created_at, updated_at)
		VALUES (1, 1, 1, 'not_started', 0, NOW(), NOW());
		INSERT INTO han_chars (id, character, pinyin, tags, categories, examples, level, created_at, updated_at)
		VALUES (1, '测', 'cè', '[]'::jsonb, '[]'::jsonb, '[]'::jsonb, 1, NOW(), NOW());
	`).Error
	require.NoError(t, err)

	// Return cleanup function
	cleanup := func() {
		t.Logf("=== Cleaning up test: %s ===", t.Name())
		conn.Close()
		app.Stop(ctx)
		test.CleanupTestDB()
	}

	return ctx, client, vocabClient, cleanup
}

func TestGetCourseProgress(t *testing.T) {
	ctx, client, _, cleanup := setupTest(t)
	defer cleanup()

	// 4. Execute test
	req := &pb.LearningServiceGetCourseProgressRequest{
		CourseId: 1,
	}
	ctx = service.WithUserID(ctx, entity.UID(1))
	resp, err := client.GetCourseProgress(ctx, req)
	require.NoError(t, err)

	// 5. Verify results
	assert.Equal(t, float32(0.0), resp.Progress.Progress)
	assert.Equal(t, uint32(0), resp.Progress.CompletedItems)
	assert.Equal(t, uint32(1), resp.Progress.TotalItems)
}

func TestGetSectionProgress(t *testing.T) {
	ctx, client, _, cleanup := setupTest(t)
	defer cleanup()

	// 4. Execute test
	req := &pb.LearningServiceGetSectionProgressRequest{
		SectionId: 1,
	}
	ctx = service.WithUserID(ctx, entity.UID(1))
	resp, err := client.GetSectionProgress(ctx, req)
	require.NoError(t, err)

	// 5. Verify results
	assert.Equal(t, float32(0.0), resp.Progress.Progress)
	assert.Equal(t, uint32(0), resp.Progress.CompletedItems)
	assert.Equal(t, uint32(0), resp.Progress.TotalItems)
}

func TestUpdateUnitProgress(t *testing.T) {
	ctx, client, _, cleanup := setupTest(t)
	defer cleanup()

	// 4. Execute test
	req := &pb.LearningServiceUpdateUnitProgressRequest{
		UnitId:    1,
		SectionId: 1,
		Completed: true,
	}
	ctx = service.WithUserID(ctx, entity.UID(1))
	resp, err := client.UpdateUnitProgress(ctx, req)
	require.NoError(t, err)
	assert.NotNil(t, resp)
}
