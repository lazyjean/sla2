package word

import (
	"context"
	"time"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type WordService struct {
	pb.UnimplementedWordServiceServer
	wordRepo         repository.WordRepository
	wordLearningRepo repository.WordLearningRepository
}

func NewWordService(wordRepo repository.WordRepository, wordLearningRepo repository.WordLearningRepository) *WordService {
	return &WordService{
		wordRepo:         wordRepo,
		wordLearningRepo: wordLearningRepo,
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

func (s *WordService) UpdateLearningProgress(ctx context.Context, req *pb.UpdateLearningProgressRequest) (*pb.UpdateLearningProgressResponse, error) {
	now := time.Now()
	nextReviewAt := now.Add(24 * time.Hour) // 简单起见，固定24小时后复习

	progress, err := s.wordLearningRepo.UpdateProgress(ctx, uint(req.UserId), uint(req.WordId), int(req.Familiarity), nextReviewAt)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateLearningProgressResponse{
		Progress: &pb.LearningProgress{
			UserId:         int64(progress.UserID),
			WordId:         int64(progress.WordID),
			Familiarity:    int32(progress.Familiarity),
			NextReviewAt:   timestamppb.New(progress.NextReviewAt),
			LastReviewedAt: timestamppb.New(progress.LastReviewedAt),
		},
	}, nil
}
