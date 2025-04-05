package word

import (
	"context"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type WordService struct {
	pb.UnimplementedWordServiceServer
	wordService *service.WordService
}

func NewWordService(wordService *service.WordService) *WordService {
	return &WordService{
		wordService: wordService,
	}
}

func (s *WordService) GetWord(ctx context.Context, req *pb.WordServiceGetRequest) (*pb.WordServiceGetResponse, error) {
	word, err := s.wordService.GetWord(ctx, uint(req.WordId))
	if err != nil {
		return nil, err
	}

	var definitions []string
	for _, def := range word.Definitions {
		definitions = append(definitions, def.Meaning)
	}

	return &pb.WordServiceGetResponse{
		Word: &pb.WordInfo{
			Id:            uint32(word.ID),
			Spelling:      word.Text,
			Pronunciation: word.Phonetic,
			Definitions:   definitions,
			Examples:      word.Examples,
			CreatedAt:     timestamppb.New(word.CreatedAt),
			UpdatedAt:     timestamppb.New(word.UpdatedAt),
		},
	}, nil
}

func (s *WordService) ListWords(ctx context.Context, req *pb.WordServiceListRequest) (*pb.WordServiceListResponse, error) {
	words, total, err := s.wordService.ListWords(ctx, 0, int(req.Page)*int(req.PageSize), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	var pbWords []*pb.WordInfo
	for _, word := range words {
		var definitions []string
		for _, def := range word.Definitions {
			definitions = append(definitions, def.Meaning)
		}

		pbWords = append(pbWords, &pb.WordInfo{
			Id:            uint32(word.ID),
			Spelling:      word.Text,
			Pronunciation: word.Phonetic,
			Definitions:   definitions,
			Examples:      word.Examples,
			CreatedAt:     timestamppb.New(word.CreatedAt),
			UpdatedAt:     timestamppb.New(word.UpdatedAt),
		})
	}

	return &pb.WordServiceListResponse{
		Words: pbWords,
		Total: uint32(total),
	}, nil
}
