package postgres

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lazyjean/sla2/config"
	"github.com/lazyjean/sla2/domain/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB 创建并初始化数据库连接
func NewDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)

	// 配置 GORM 日志
	logConfig := logger.Config{
		SlowThreshold:             time.Second,  // 慢 SQL 阈值
		LogLevel:                  logger.Error, // 默认只记录错误
		IgnoreRecordNotFoundError: true,         // 忽略记录未找到的错误
		Colorful:                  false,        // 禁用彩色输出
	}

	if cfg.Debug {
		logConfig.LogLevel = logger.Info // 调试模式下记录所有 SQL
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logConfig,
		),
	})
	if err != nil {
		return nil, err
	}

	// 自动迁移
	if err := autoMigrate(db); err != nil {
		return nil, err
	}

	// 获取底层 *sql.DB 并配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置连接池
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return db, nil
}

// autoMigrate 自动迁移数据库结构
func autoMigrate(db *gorm.DB) error {
	// 删除现有表
	db.Exec("DROP TABLE IF EXISTS word_tags CASCADE")
	db.Exec("DROP TABLE IF EXISTS tags CASCADE")
	db.Exec("DROP TABLE IF EXISTS examples CASCADE")
	db.Exec("DROP TABLE IF EXISTS words CASCADE")

	// 创建表
	err := db.AutoMigrate(
		&entity.Word{},
		&entity.Example{},
		&entity.Tag{},
	)
	if err != nil {
		return err
	}

	// 创建索引
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_words_text ON words (text)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_tags_name ON tags (name)`)

	return nil
}
