package valute

import (
	"errors"
	"fmt"
	"hash/fnv"
	"math"
	"math/rand"
	"sort"
	"time"
)

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(date time.Time) (ValuteCurrencies, error) {
	if date.IsZero() {
		return ValuteCurrencies{}, errors.New("date is required")
	}

	normalizedDate := truncateToDate(date.UTC())
	seed := stableSeed(normalizedDate)
	rnd := rand.New(rand.NewSource(seed))

	base := baseCurrencies()
	extraPool := extraCurrencies()

	extraCount := 2 + rnd.Intn(4)

	perm := rnd.Perm(len(extraPool))
	selected := make([]currencyMeta, 0, len(base)+extraCount)
	selected = append(selected, base...)
	for i := 0; i < extraCount && i < len(extraPool); i++ {
		selected = append(selected, extraPool[perm[i]])
	}
	sort.Slice(selected, func(i, j int) bool {
		return selected[i].CharCode < selected[j].CharCode
	})

	valutes := make([]Valute, 0, len(selected))
	for idx, meta := range selected {
		valuePerUnit := generateUnitRate(rnd, meta, normalizedDate, idx)
		value := round4(valuePerUnit * float64(meta.Nominal))

		valutes = append(valutes, Valute{
			NumCode:   meta.NumCode,
			CharCode:  meta.CharCode,
			Nominal:   meta.Nominal,
			Name:      meta.Name,
			Value:     value,
			VunitRate: round4(valuePerUnit),
		})
	}

	return ValuteCurrencies{
		Date:    normalizedDate,
		Name:    "Foreign Currency Market",
		Valutes: valutes,
	}, nil
}

type currencyMeta struct {
	NumCode  string
	CharCode string
	Nominal  int
	Name     string

	MinRate float64
	MaxRate float64
}

func baseCurrencies() []currencyMeta {
	return []currencyMeta{
		{
			NumCode:  "036",
			CharCode: "AUD",
			Nominal:  1,
			Name:     "Австралийский доллар",
			MinRate:  45,
			MaxRate:  75,
		},
		{
			NumCode:  "978",
			CharCode: "EUR",
			Nominal:  1,
			Name:     "Евро",
			MinRate:  85,
			MaxRate:  115,
		},
		{
			NumCode:  "840",
			CharCode: "USD",
			Nominal:  1,
			Name:     "Доллар США",
			MinRate:  75,
			MaxRate:  105,
		},
		{
			NumCode:  "156",
			CharCode: "CNY",
			Nominal:  10,
			Name:     "Китайских юаней",
			MinRate:  10,
			MaxRate:  16,
		},
	}
}

func extraCurrencies() []currencyMeta {
	return []currencyMeta{
		{
			NumCode:  "933",
			CharCode: "BYN",
			Nominal:  1,
			Name:     "Белорусский рубль",
			MinRate:  20,
			MaxRate:  40,
		},
		{
			NumCode:  "124",
			CharCode: "CAD",
			Nominal:  1,
			Name:     "Канадский доллар",
			MinRate:  50,
			MaxRate:  80,
		},
		{
			NumCode:  "756",
			CharCode: "CHF",
			Nominal:  1,
			Name:     "Швейцарский франк",
			MinRate:  85,
			MaxRate:  120,
		},
		{
			NumCode:  "826",
			CharCode: "GBP",
			Nominal:  1,
			Name:     "Фунт стерлингов",
			MinRate:  95,
			MaxRate:  130,
		},
		{
			NumCode:  "392",
			CharCode: "JPY",
			Nominal:  100,
			Name:     "Японских иен",
			MinRate:  45,
			MaxRate:  80,
		},
		{
			NumCode:  "398",
			CharCode: "KZT",
			Nominal:  100,
			Name:     "Тенге",
			MinRate:  12,
			MaxRate:  25,
		},
		{
			NumCode:  "498",
			CharCode: "MDL",
			Nominal:  10,
			Name:     "Молдавских леев",
			MinRate:  40,
			MaxRate:  65,
		},
		{
			NumCode:  "578",
			CharCode: "NOK",
			Nominal:  10,
			Name:     "Норвежских крон",
			MinRate:  70,
			MaxRate:  110,
		},
		{
			NumCode:  "985",
			CharCode: "PLN",
			Nominal:  1,
			Name:     "Польский злотый",
			MinRate:  18,
			MaxRate:  30,
		},
		{
			NumCode:  "752",
			CharCode: "SEK",
			Nominal:  10,
			Name:     "Шведских крон",
			MinRate:  70,
			MaxRate:  110,
		},
		{
			NumCode:  "949",
			CharCode: "TRY",
			Nominal:  10,
			Name:     "Турецких лир",
			MinRate:  20,
			MaxRate:  45,
		},
		{
			NumCode:  "980",
			CharCode: "UAH",
			Nominal:  10,
			Name:     "Украинских гривен",
			MinRate:  18,
			MaxRate:  35,
		},
	}
}

func generateUnitRate(rnd *rand.Rand, meta currencyMeta, date time.Time, idx int) float64 {
	span := meta.MaxRate - meta.MinRate
	base := meta.MinRate + rnd.Float64()*span

	monthFactor := 1.0 + (float64(date.Month())-6.5)/100.0

	yearFactor := 1.0 + float64(date.Year()-2000)/500.0

	dayFactor := 1.0 + float64((date.Day()+idx)%7-3)/300.0

	rate := base * monthFactor * yearFactor * dayFactor
	if rate < 0.0001 {
		rate = 0.0001
	}

	return round4(rate)
}

func stableSeed(date time.Time) int64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(date.Format("2006-01-02")))
	return int64(h.Sum64())
}

func truncateToDate(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func round4(v float64) float64 {
	return math.Round(v*10000) / 10000
}

func DebugString(v ValuteCurrencies) string {
	return fmt.Sprintf("date=%s name=%s valutes=%d", v.Date.Format("2006-01-02"), v.Name, len(v.Valutes))
}
