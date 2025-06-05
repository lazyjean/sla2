package converter

import (
	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/valueobject"
)

type VocabularyConverter struct{}

func NewVocabularyConverter() *VocabularyConverter {
	return &VocabularyConverter{}
}

// ToProtoWord 将实体转换为 Proto Word 消息
func (c *VocabularyConverter) ToProtoWord(word *entity.Word) *pb.Word {
	if word == nil {
		return nil
	}

	var definitions []*pb.WordDefinition
	for _, def := range word.Definitions {
		definitions = append(definitions, &pb.WordDefinition{
			PartOfSpeech: pb.WordPartOfSpeech(pb.WordPartOfSpeech_value[def.PartOfSpeech]),
			Meaning:      def.Meaning,
			Example:      def.Example,
			Synonyms:     def.Synonyms,
			Antonyms:     def.Antonyms,
		})
	}

	return &pb.Word{
		Id:          uint32(word.ID),
		Word:        word.Text,
		Spelling:    word.Phonetic,
		Level:       c.ConvertLevelToProto(word.Level),
		Definitions: definitions,
		Examples:    word.Examples,
		Tags:        word.Tags,
	}
}

// ToProtoHanChar 将汉字实体转换为 Proto 消息
func (c *VocabularyConverter) ToProtoHanChar(hanChar *entity.HanChar) *pb.HanChar {
	return &pb.HanChar{
		Id:         uint32(hanChar.ID),
		Character:  hanChar.Character,
		Pinyin:     hanChar.Pinyin,
		Tags:       hanChar.Tags,
		Categories: hanChar.Categories,
		Examples:   hanChar.Examples,
		Level:      c.ConvertLevelToProto(hanChar.Level),
	}
}

// PbToEntityHanChar 将 Proto 汉字消息转换为实体
func (c *VocabularyConverter) PbToEntityHanChar(hanCharPb *pb.HanChar) *entity.HanChar {
	return &entity.HanChar{
		Character:  hanCharPb.Character,
		Pinyin:     hanCharPb.Pinyin,
		Level:      c.ConvertLevelToValueObject(hanCharPb.Level),
		Tags:       hanCharPb.Tags,
		Categories: hanCharPb.Categories,
		Examples:   hanCharPb.Examples,
	}
}

// ConvertLevelToString 将 protobuf 枚举类型转换为字符串
func (c *VocabularyConverter) ConvertLevelToString(level pb.WordDifficultyLevel) string {
	switch level {
	case pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_A1:
		return valueobject.WORD_DIFFICULTY_LEVEL_A1.String()
	case pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_A2:
		return valueobject.WORD_DIFFICULTY_LEVEL_A2.String()
	case pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_B1:
		return valueobject.WORD_DIFFICULTY_LEVEL_B1.String()
	case pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_B2:
		return valueobject.WORD_DIFFICULTY_LEVEL_B2.String()
	case pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_C1:
		return valueobject.WORD_DIFFICULTY_LEVEL_C1.String()
	case pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_C2:
		return valueobject.WORD_DIFFICULTY_LEVEL_C2.String()
	case pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_HSK1:
		return valueobject.WORD_DIFFICULTY_LEVEL_HSK1.String()
	case pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_HSK2:
		return valueobject.WORD_DIFFICULTY_LEVEL_HSK2.String()
	case pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_HSK3:
		return valueobject.WORD_DIFFICULTY_LEVEL_HSK3.String()
	case pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_HSK4:
		return valueobject.WORD_DIFFICULTY_LEVEL_HSK4.String()
	case pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_HSK5:
		return valueobject.WORD_DIFFICULTY_LEVEL_HSK5.String()
	case pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_HSK6:
		return valueobject.WORD_DIFFICULTY_LEVEL_HSK6.String()
	default:
		return valueobject.WORD_DIFFICULTY_LEVEL_A1.String()
	}
}

// ConvertStringToLevel 将字符串转换为 protobuf 枚举类型
func (c *VocabularyConverter) ConvertStringToLevel(level string) pb.WordDifficultyLevel {
	switch level {
	case valueobject.WORD_DIFFICULTY_LEVEL_A1.String():
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_A1
	case valueobject.WORD_DIFFICULTY_LEVEL_A2.String():
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_A2
	case valueobject.WORD_DIFFICULTY_LEVEL_B1.String():
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_B1
	case valueobject.WORD_DIFFICULTY_LEVEL_B2.String():
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_B2
	case valueobject.WORD_DIFFICULTY_LEVEL_C1.String():
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_C1
	case valueobject.WORD_DIFFICULTY_LEVEL_C2.String():
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_C2
	case valueobject.WORD_DIFFICULTY_LEVEL_HSK1.String():
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_HSK1
	case valueobject.WORD_DIFFICULTY_LEVEL_HSK2.String():
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_HSK2
	case valueobject.WORD_DIFFICULTY_LEVEL_HSK3.String():
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_HSK3
	case valueobject.WORD_DIFFICULTY_LEVEL_HSK4.String():
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_HSK4
	case valueobject.WORD_DIFFICULTY_LEVEL_HSK5.String():
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_HSK5
	case valueobject.WORD_DIFFICULTY_LEVEL_HSK6.String():
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_HSK6
	default:
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_A1
	}
}

// ConvertLevelToValueObject 将Proto难度等级转换为值对象
func (c *VocabularyConverter) ConvertLevelToValueObject(level pb.WordDifficultyLevel) valueobject.WordDifficultyLevel {
	switch level {
	case pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_UNSPECIFIED:
		return valueobject.WORD_DIFFICULTY_LEVEL_UNSPECIFIED
	case pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_A1:
		return valueobject.WORD_DIFFICULTY_LEVEL_A1
	case pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_A2:
		return valueobject.WORD_DIFFICULTY_LEVEL_A2
	case pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_B1:
		return valueobject.WORD_DIFFICULTY_LEVEL_B1
	case pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_B2:
		return valueobject.WORD_DIFFICULTY_LEVEL_B2
	case pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_C1:
		return valueobject.WORD_DIFFICULTY_LEVEL_C1
	case pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_C2:
		return valueobject.WORD_DIFFICULTY_LEVEL_C2
	case pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_HSK1:
		return valueobject.WORD_DIFFICULTY_LEVEL_HSK1
	case pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_HSK2:
		return valueobject.WORD_DIFFICULTY_LEVEL_HSK2
	case pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_HSK3:
		return valueobject.WORD_DIFFICULTY_LEVEL_HSK3
	default:
		return valueobject.WORD_DIFFICULTY_LEVEL_UNSPECIFIED
	}
}

// ConvertLevelToProto 将领域值对象转换为Proto的难度等级
func (c *VocabularyConverter) ConvertLevelToProto(level valueobject.WordDifficultyLevel) pb.WordDifficultyLevel {
	switch level {
	case valueobject.WORD_DIFFICULTY_LEVEL_UNSPECIFIED:
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_UNSPECIFIED
	case valueobject.WORD_DIFFICULTY_LEVEL_A1:
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_A1
	case valueobject.WORD_DIFFICULTY_LEVEL_A2:
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_A2
	case valueobject.WORD_DIFFICULTY_LEVEL_B1:
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_B1
	case valueobject.WORD_DIFFICULTY_LEVEL_B2:
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_B2
	case valueobject.WORD_DIFFICULTY_LEVEL_C1:
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_C1
	case valueobject.WORD_DIFFICULTY_LEVEL_C2:
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_C2
	case valueobject.WORD_DIFFICULTY_LEVEL_HSK1:
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_HSK1
	case valueobject.WORD_DIFFICULTY_LEVEL_HSK2:
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_HSK2
	case valueobject.WORD_DIFFICULTY_LEVEL_HSK3:
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_HSK3
	case valueobject.WORD_DIFFICULTY_LEVEL_HSK4:
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_HSK4
	case valueobject.WORD_DIFFICULTY_LEVEL_HSK5:
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_HSK5
	case valueobject.WORD_DIFFICULTY_LEVEL_HSK6:
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_HSK6
	default:
		return pb.WordDifficultyLevel_WORD_DIFFICULTY_LEVEL_UNSPECIFIED
	}
}

// PbWordsToEntities 将 []*pb.Word 转换为 []*entity.Word
func (c *VocabularyConverter) PbWordsToEntities(wordPbs []*pb.Word) []*entity.Word {
	words := make([]*entity.Word, 0, len(wordPbs))
	for _, wordPb := range wordPbs {
		words = append(words, c.PbToEntityWord(wordPb))
	}
	return words
}

// PbToEntityWord 将 pb.Word 转换为 entity.Word
func (c *VocabularyConverter) PbToEntityWord(wordPb *pb.Word) *entity.Word {
	if wordPb == nil {
		return nil
	}
	definitions := make([]entity.Definition, 0, len(wordPb.Definitions))
	for _, def := range wordPb.Definitions {
		definitions = append(definitions, entity.Definition{
			PartOfSpeech: def.PartOfSpeech.String(),
			Meaning:      def.Meaning,
			Example:      def.Example,
			Synonyms:     def.Synonyms,
			Antonyms:     def.Antonyms,
		})
	}
	return &entity.Word{
		Text:        wordPb.Word,
		Phonetic:    wordPb.Spelling,
		Definitions: definitions,
		Examples:    wordPb.Examples,
		Tags:        wordPb.Tags,
		Level:       c.ConvertLevelToValueObject(wordPb.Level),
	}
}
