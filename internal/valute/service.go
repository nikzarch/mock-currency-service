package valute

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nikzarch/mock-currency-service/internal/config"
	"golang.org/x/sync/singleflight"
)

type Service interface {
	GetReportByDate(ctx context.Context, date time.Time) (Currencies, error)
}

type ValuteService struct {
	repository Repository
	generator  *Generator
	mode       config.ResponseMode
	group      singleflight.Group
}

func NewValuteService(repository Repository, generator *Generator, mode config.ResponseMode) *ValuteService {
	return &ValuteService{
		repository: repository,
		generator:  generator,
		mode:       mode,
	}
}

func (s *ValuteService) GetReportByDate(ctx context.Context, date time.Time) (Currencies, error) {
	if s.mode == config.ResponseModeError {
		return Currencies{}, fmt.Errorf("forced error mode: %w", ErrNotFound)
	}

	report, err := s.repository.GetDailyReportByDate(ctx, date)
	if err == nil {
		return report, nil
	}
	if !errors.Is(err, ErrNotFound) {
		return Currencies{}, err
	}

	key := date.UTC().Format(time.DateOnly)
	v, err, _ := s.group.Do(key, func() (any, error) {
		report, err := s.repository.GetDailyReportByDate(ctx, date)
		if err == nil {
			return report, nil
		}
		if !errors.Is(err, ErrNotFound) {
			return Currencies{}, err
		}

		currs, err := s.generator.Generate(date)
		if err != nil {
			return Currencies{}, err
		}

		saveCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := s.repository.AddDailyReport(saveCtx, currs); err != nil {
			return Currencies{}, err
		}

		return currs, nil
	})
	if err != nil {
		return Currencies{}, err
	}

	return v.(Currencies), nil
}
