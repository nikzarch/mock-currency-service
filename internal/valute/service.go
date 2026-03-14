package valute

import (
	"context"
	"errors"
	"time"
)

type Service interface {
	GetReportByDate(ctx context.Context, date time.Time) (Currencies, error)
}

type ValuteService struct {
	repository Repository
	generator  Generator
}

func NewService(repository Repository, generator Generator) Service {
	return &ValuteService{repository: repository, generator: generator}
}

func (s *ValuteService) GetReportByDate(ctx context.Context, date time.Time) (Currencies, error) {
	report, err := s.repository.GetDailyReportByDate(ctx, date)
	if errors.Is(err, ErrNotFound) {
		currs, err := s.generator.Generate(date)
		if err != nil {
			return Currencies{}, err
		}
		// TODO: validate
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		err = s.repository.AddDailyReport(ctx, currs)
		if err != nil {
			return Currencies{}, err
		}
	}
	if err != nil {
		return Currencies{}, err
	}
	return report, nil
}
