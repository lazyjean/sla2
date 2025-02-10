package learning

import (
	"context"
	"time"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LearningService struct {
	pb.UnimplementedLearningServiceServer
	learningRepo repository.LearningRepository
}

func NewLearningService(learningRepo repository.LearningRepository) *LearningService {
	return &LearningService{
		learningRepo: learningRepo,
	}
}

func (s *LearningService) UpdateLearningProgress(ctx context.Context, req *pb.UpdateLearningProgressRequest) (*pb.UpdateLearningProgressResponse, error) {
	now := time.Now()
	nextReviewAt := now.Add(24 * time.Hour) // 简单起见，固定24小时后复习

	progress, err := s.learningRepo.UpdateProgress(ctx, uint(req.UserId), uint(req.WordId), int(req.Familiarity), nextReviewAt)
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

func (s *LearningService) ListLearningProgress(ctx context.Context, req *pb.ListLearningProgressRequest) (*pb.ListLearningProgressResponse, error) {
	progresses, total, err := s.learningRepo.ListByUserID(ctx, uint(req.UserId), int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	var pbProgresses []*pb.LearningProgress
	for _, progress := range progresses {
		pbProgresses = append(pbProgresses, &pb.LearningProgress{
			UserId:         int64(progress.UserID),
			WordId:         int64(progress.WordID),
			Familiarity:    int32(progress.Familiarity),
			NextReviewAt:   timestamppb.New(progress.NextReviewAt),
			LastReviewedAt: timestamppb.New(progress.LastReviewedAt),
		})
	}

	return &pb.ListLearningProgressResponse{
		ProgressList: pbProgresses,
		Total:        int32(total),
	}, nil
}

func (s *LearningService) GetLearningStats(ctx context.Context, req *pb.GetLearningStatsRequest) (*pb.GetLearningStatsResponse, error) {
	stats, err := s.learningRepo.GetUserStats(ctx, uint(req.UserId))
	if err != nil {
		return nil, err
	}

	return &pb.GetLearningStatsResponse{
		Stats: &pb.LearningStats{
			UserId:          int64(stats.UserID),
			TotalWords:      int32(stats.TotalWords),
			MasteredWords:   int32(stats.MasteredWords),
			LearningWords:   int32(stats.LearningWords),
			ReviewDueCount:  int32(stats.ReviewDueCount),
			LastStudyTime:   timestamppb.New(stats.LastStudyTime),
			TodayStudyCount: int32(stats.TodayStudyCount),
			ContinuousDays:  int32(stats.ContinuousDays),
		},
	}, nil
}

func (s *LearningService) ListReviewWords(ctx context.Context, req *pb.ListReviewWordsRequest) (*pb.ListReviewWordsResponse, error) {
	words, progresses, total, err := s.learningRepo.ListReviewWords(ctx, uint(req.UserId), int(req.Page), int(req.PageSize))
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

	var pbProgresses []*pb.LearningProgress
	for _, progress := range progresses {
		pbProgresses = append(pbProgresses, &pb.LearningProgress{
			UserId:         int64(progress.UserID),
			WordId:         int64(progress.WordID),
			Familiarity:    int32(progress.Familiarity),
			NextReviewAt:   timestamppb.New(progress.NextReviewAt),
			LastReviewedAt: timestamppb.New(progress.LastReviewedAt),
		})
	}

	return &pb.ListReviewWordsResponse{
		Words:        pbWords,
		ProgressList: pbProgresses,
		Total:        int32(total),
	}, nil
}
