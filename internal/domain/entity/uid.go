package entity

import (
	"database/sql/driver"
)

type UID uint64 // Custom type for domain clarity

// Value 实现GORM的Valuer接口
func (u UID) Value() (driver.Value, error) {
	return uint64(u), nil
}

// Scan 实现GORM的Scanner接口
func (u *UID) Scan(value interface{}) error {
	if value == nil {
		*u = 0
		return nil
	}
	switch v := value.(type) {
	case int64:
		*u = UID(v)
	case uint64:
		*u = UID(v)
	}
	return nil
}
