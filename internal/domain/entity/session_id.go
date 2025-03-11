package entity

import (
	"database/sql/driver"
	"strconv"
)

type SessionID uint64 // Session identifier type

// Value 实现GORM的Valuer接口
func (s SessionID) Value() (driver.Value, error) {
	return uint64(s), nil
}

// Scan 实现GORM的Scanner接口
func (s *SessionID) Scan(value interface{}) error {
	if value == nil {
		*s = 0
		return nil
	}
	switch v := value.(type) {
	case int64:
		*s = SessionID(v)
	case uint64:
		*s = SessionID(v)
	}
	return nil
}

// FromString 将字符串转换为SessionID
func (s *SessionID) FromString(str string) error {
	id, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return err
	}
	*s = SessionID(id)
	return nil
}

// String 将SessionID转换为字符串
func (s SessionID) String() string {
	return strconv.FormatUint(uint64(s), 10)
}
