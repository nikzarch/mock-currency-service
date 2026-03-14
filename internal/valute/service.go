package valute

import (
	"context"
	"errors"
	"time"
)

type Service interface {
	GetReportByDate(date time.Time, ctx context.Context) (Currencies, error)
}

type ValuteService struct {
	repository Repository
	generator  Generator
}

func NewService(repository Repository, generator Generator) Service {
	return &ValuteService{repository: repository, generator: generator}
}

func (s *ValuteService) GetReportByDate(date time.Time, ctx context.Context) (Currencies, error) {
	report, err := s.repository.GetDailyReportByDate(date, ctx)
	if errors.Is(err, ErrNotFound) {
		currs, err := s.generator.Generate(date)
		if err != nil {
			return Currencies{}, err
		}
		// TODO: validate
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		err = s.repository.AddDailyReport(currs, ctx)
		if err != nil {
			return Currencies{}, err
		}
	}
	if err != nil {
		return Currencies{}, err
	}
	return report, nil
}
