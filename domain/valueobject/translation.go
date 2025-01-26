package valueobject

import "github.com/lazyjean/sla2/domain/errors"

// Translation 翻译值对象
type Translation struct {
	Primary   string   // 主要翻译
	Secondary []string // 次要翻译
}

// NewTranslation 创建新的翻译值对象
func NewTranslation(primary string, secondary []string) (Translation, error) {
	if primary == "" {
		return Translation{}, errors.ErrEmptyPrimaryTranslation
	}

	if secondary == nil {
		secondary = make([]string, 0)
	}

	return Translation{
		Primary:   primary,
		Secondary: secondary,
	}, nil
}

// Equals 值对象的相等性比较
func (t Translation) Equals(other Translation) bool {
	if t.Primary != other.Primary {
		return false
	}
	if len(t.Secondary) != len(other.Secondary) {
		return false
	}
	for i, s := range t.Secondary {
		if s != other.Secondary[i] {
			return false
		}
	}
	return true
}
