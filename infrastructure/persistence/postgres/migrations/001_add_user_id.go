package migrations

import (
	"gorm.io/gorm"
)

// AddUserIDToWords 添加 user_id 字段到 words 表
func AddUserIDToWords(db *gorm.DB) error {
	// 1. 检查列是否已存在
	var count int64
	db.Raw(`SELECT COUNT(*) FROM information_schema.columns 
		WHERE table_name = 'words' AND column_name = 'user_id'`).Count(&count)
	if count > 0 {
		return nil // 列已存在，不需要执行迁移
	}

	// 2. 开启事务
	return db.Transaction(func(tx *gorm.DB) error {
		// 3. 添加可空的 user_id 列
		if err := tx.Exec(`ALTER TABLE words ADD COLUMN user_id bigint`).Error; err != nil {
			return err
		}

		// 4. 为现有记录设置默认值（比如设置为1，表示系统用户）
		if err := tx.Exec(`UPDATE words SET user_id = 1 WHERE user_id IS NULL`).Error; err != nil {
			return err
		}

		// 5. 将列设置为非空
		if err := tx.Exec(`ALTER TABLE words ALTER COLUMN user_id SET NOT NULL`).Error; err != nil {
			return err
		}

		// 6. 添加索引
		if err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_words_user_id ON words(user_id)`).Error; err != nil {
			return err
		}

		// 7. 添加联合唯一索引
		return tx.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_user_text ON words(user_id, text)`).Error
	})
}
