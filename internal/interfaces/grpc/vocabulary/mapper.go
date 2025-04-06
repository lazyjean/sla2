package vocabulary

import (
	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/valueobject"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToEntity 将 DTO 转换为实体
func ToEntity(createDTO *dto.WordCreateDTO, userID entity.UID) (*entity.Word, error) {
	return entity.NewWord(
		createDTO.Text,
		createDTO.Phonetic,
		createDTO.Definitions,
		createDTO.Examples,
		createDTO.Tags,
	)
}

// ToDTO 将实体转换为 DTO
func ToDTO(word *entity.Word) *dto.WordResponseDTO {
	return &dto.WordResponseDTO{
		ID:          uint32(word.ID),
		Text:        word.Text,
		Definitions: word.Definitions,
		Phonetic:    word.Phonetic,
		Examples:    word.Examples,
		Tags:        word.Tags,
		CreatedAt:   word.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   word.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// ToProto 将实体转换为 Proto 消息
func ToProto(word *entity.Word) *pb.WordInfo {
	// 将 Definitions 转换为字符串数组
	var definitions []string
	for _, def := range word.Definitions {
		definitions = append(definitions, def.Meaning)
	}

	return &pb.WordInfo{
		Id:            uint32(word.ID),
		Spelling:      word.Text,
		Pronunciation: word.Phonetic,
		Definitions:   definitions,
		Examples:      word.Examples,
		CreatedAt:     timestamppb.New(word.CreatedAt),
		UpdatedAt:     timestamppb.New(word.UpdatedAt),
	}
}

// ToProtoWord 将实体转换为 Proto Word 消息
func ToProtoWord(word *entity.Word) *pb.Word {
	return &pb.Word{
		Id:            uint32(word.ID),
		Word:          word.Text,
		Spelling:      word.Phonetic,
		Pronunciation: word.Phonetic,
		Difficulty:    pb.WordDifficultyLevel(word.Difficulty),
		Examples:      word.Examples,
		Tags:          word.Tags,
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
		Level:      ConvertStringToLevel(hanChar.Level.String()),
	}
}

// ToEntityFromProto 将 Proto 消息转换为实体
func ToEntityFromProto(proto *pb.WordInfo, userID entity.UID) (*entity.Word, error) {
	// 将字符串数组转换为 Definition 数组
	var definitions []entity.Definition
	for _, meaning := range proto.Definitions {
		definitions = append(definitions, entity.Definition{
			Meaning: meaning,
		})
	}

	return entity.NewWord(
		proto.Spelling,
		proto.Pronunciation,
		definitions,
		proto.Examples,
		nil, // 标签为空
	)
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
