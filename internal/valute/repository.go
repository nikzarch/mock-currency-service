package valute

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed queries/get_daily_report.sql
var getDailyReportQuery string

//go:embed queries/get_valutes_by_report_id.sql
var getValutesByReportIDQuery string

//go:embed queries/insert_daily_report.sql
var insertDailyReportQuery string

//go:embed queries/insert_valute.sql
var insertValuteQuery string

var (
	ErrNotFound = errors.New("report not found")
)

type Repository interface {
	GetDailyReportByDate(date time.Time, ctx context.Context) (Currencies, error)
	AddDailyReport(report Currencies, ctx context.Context) error
}

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) Repository {
	return &PostgresRepository{pool: pool}
}

func (p *PostgresRepository) GetDailyReportByDate(date time.Time, ctx context.Context) (Currencies, error) {
	var (
		reportID   int64
		reportDate time.Time
		reportName string
	)

	err := p.pool.QueryRow(ctx, getDailyReportQuery, date).Scan(
		&reportID,
		&reportDate,
		&reportName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Currencies{}, ErrNotFound
		}
		return Currencies{}, fmt.Errorf("get daily report: %w", err)
	}

	rows, err := p.pool.Query(ctx, getValutesByReportIDQuery, reportID)
	if err != nil {
		return Currencies{}, fmt.Errorf("get valutes by report id: %w", err)
	}
	defer rows.Close()

	valutes := make([]Valute, 0, 8)

	for rows.Next() {
		var v Valute

		err = rows.Scan(
			&v.NumCode,
			&v.CharCode,
			&v.Nominal,
			&v.Name,
			&v.Value,
			&v.VunitRate,
		)
		if err != nil {
			return Currencies{}, fmt.Errorf("scan valute: %w", err)
		}

		valutes = append(valutes, v)
	}

	if err = rows.Err(); err != nil {
		return Currencies{}, fmt.Errorf("iterate valutes: %w", err)
	}

	return Currencies{
		Date:    reportDate,
		Name:    reportName,
		Valutes: valutes,
	}, nil
}

func (p *PostgresRepository) AddDailyReport(report Currencies, ctx context.Context) error {

	tx, err := p.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var reportID int64
	err = tx.QueryRow(ctx, insertDailyReportQuery, report.Date, report.Name).Scan(&reportID)
	if err != nil {
		return fmt.Errorf("insert daily report: %w", err)
	}

	for _, v := range report.Valutes {
		_, err = tx.Exec(
			ctx,
			insertValuteQuery,
			reportID,
			v.NumCode,
			v.CharCode,
			v.Nominal,
			v.Name,
			v.Value,
			v.VunitRate,
		)
		if err != nil {
			return fmt.Errorf("insert valute %s: %w", v.CharCode, err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func normalizeDate(dateReq time.Time) string {
	return dateReq.Format(time.DateOnly)
}
