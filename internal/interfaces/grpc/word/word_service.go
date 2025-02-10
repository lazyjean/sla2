package word

import (
	"context"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type WordService struct {
	pb.UnimplementedWordServiceServer
	wordRepo repository.WordRepository
}

func NewWordService(wordRepo repository.WordRepository) *WordService {
	return &WordService{
		wordRepo: wordRepo,
	}
}

func (s *WordService) GetWord(ctx context.Context, req *pb.GetWordRequest) (*pb.GetWordResponse, error) {
	word, err := s.wordRepo.FindByID(ctx, uint(req.WordId))
	if err != nil {
		return nil, err
	}

	return &pb.GetWordResponse{
		Word: &pb.Word{
			Id:            int64(word.ID),
			Spelling:      word.Text,
			Pronunciation: word.Phonetic,
			Definitions:   []string{word.Translation},
			Examples:      word.Examples,
			CreatedAt:     timestamppb.New(word.CreatedAt),
			UpdatedAt:     timestamppb.New(word.UpdatedAt),
		},
	}, nil
}

func (s *WordService) ListWords(ctx context.Context, req *pb.ListWordsRequest) (*pb.ListWordsResponse, error) {
	words, total, err := s.wordRepo.List(ctx, 0, int(req.Page)*int(req.PageSize), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	var pbWords []*pb.Word
	for _, word := range words {
		pbWords = append(pbWords, &pb.Word{
			Id:            int64(word.ID),
			Spelling:      word.Text,
			Pronunciation: word.Phonetic,
			Definitions:   []string{word.Translation},
			Examples:      word.Examples,
			CreatedAt:     timestamppb.New(word.CreatedAt),
			UpdatedAt:     timestamppb.New(word.UpdatedAt),
		})
	}

	return &pb.ListWordsResponse{
		Words: pbWords,
		Total: int32(total),
	}, nil
}
