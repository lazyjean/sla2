package postgres

import (
	"os"
	"testing"

	"github.com/lazyjean/sla2/config"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

func SetupTestDB(t *testing.T) (*gorm.DB, func()) {
	// 从环境变量获取测试数据库配置
	cfg := &config.DatabaseConfig{
		Host:     getEnvOrDefault("TEST_DB_HOST", "localhost"),
		Port:     getEnvOrDefault("TEST_DB_PORT", "5432"),
		User:     getEnvOrDefault("TEST_DB_USER", "sla"),
		Password: getEnvOrDefault("TEST_DB_PASSWORD", "sla1234"),
		DBName:   getEnvOrDefault("TEST_DB_NAME", "sla2_test"),
		Debug:    false,
	}

	// 连接数据库
	db, err := NewDB(cfg)
	require.NoError(t, err)

	// 确保表结构正确
	err = db.AutoMigrate(
		&entity.Word{},
		&entity.CourseLearningProgress{},
		&entity.CourseSectionProgress{},
		&entity.CourseSectionUnitProgress{},
		&entity.HanChar{},
	)
	require.NoError(t, err)

	// 清理测试数据
	cleanTestData(t, db)

	// 返回清理函数
	cleanup := func() {
		cleanTestData(t, db)
	}

	return db, cleanup
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func cleanTestData(t *testing.T, db *gorm.DB) {
	// 清理所有相关表的数据
	tables := []string{
		"words",
		"course_learning_progresses",
		"course_section_progresses",
		"course_section_unit_progresses",
		"han_chars",
	}

	for _, table := range tables {
		err := db.Exec("TRUNCATE TABLE " + table + " CASCADE").Error
		require.NoError(t, err)
	}
}
