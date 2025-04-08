package vocabulary

import (
	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/valueobject"
)

// ToProtoWord 将实体转换为 Proto Word 消息
func ToProtoWord(word *entity.Word) *pb.Word {
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
		Level:       ConvertLevelToProto(word.Level),
		Definitions: definitions,
		Examples:    word.Examples,
		Tags:        word.Tags,
	}
}

// ToProtoHanChar 将汉字实体转换为 Proto 消息
func ToProtoHanChar(hanChar *entity.HanChar) *pb.HanChar {
	return &pb.HanChar{
		Id:         uint32(hanChar.ID),
		Character:  hanChar.Character,
		Pinyin:     hanChar.Pinyin,
		Tags:       hanChar.Tags,
		Categories: hanChar.Categories,
		Examples:   hanChar.Examples,
		Level:      ConvertLevelToProto(hanChar.Level),
	}
}

// ConvertLevelToString 将 protobuf 枚举类型转换为字符串
func ConvertLevelToString(level pb.WordDifficultyLevel) string {
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
func ConvertStringToLevel(level string) pb.WordDifficultyLevel {
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
func ConvertLevelToValueObject(level pb.WordDifficultyLevel) valueobject.WordDifficultyLevel {
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
func ConvertLevelToProto(level valueobject.WordDifficultyLevel) pb.WordDifficultyLevel {
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