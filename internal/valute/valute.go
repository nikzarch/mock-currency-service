package valute

import "time"

type Currencies struct {
	Date    time.Time
	Name    string
	Valutes []Valute
}
type Valute struct {
	NumCode   string
	CharCode  string
	Nominal   int
	Name      string
	Value     float64
	VunitRate float64
}
