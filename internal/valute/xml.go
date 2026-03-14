package valute

import (
	"encoding/xml"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type xmlValCurs struct {
	XMLName xml.Name    `xml:"ValCurs"`
	Date    string      `xml:"Date,attr"`
	Name    string      `xml:"name,attr"`
	Valutes []xmlValute `xml:"Valute"`
}

type xmlValute struct {
	NumCode   string `xml:"NumCode"`
	CharCode  string `xml:"CharCode"`
	Nominal   int    `xml:"Nominal"`
	Name      string `xml:"Name"`
	Value     string `xml:"Value"`
	VunitRate string `xml:"VunitRate"`
}

func MarshalXMLDaily(currencies Currencies) ([]byte, error) {
	doc := xmlValCurs{
		Date:    currencies.Date.Format("02.01.2006"),
		Name:    currencies.Name,
		Valutes: make([]xmlValute, 0, len(currencies.Valutes)),
	}

	for _, v := range currencies.Valutes {
		doc.Valutes = append(doc.Valutes, xmlValute{
			NumCode:   v.NumCode,
			CharCode:  v.CharCode,
			Nominal:   v.Nominal,
			Name:      v.Name,
			Value:     floatToXMLDecimal(v.Value),
			VunitRate: floatToXMLDecimal(v.VunitRate),
		})
	}

	body, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, err
	}
	log.Printf("Generated XML\n%s", string(body))
	result := append([]byte(xml.Header), body...)
	return result, nil
}

func floatToXMLDecimal(v float64) string {
	s := strconv.FormatFloat(v, 'f', 4, 64)
	return strings.Replace(s, ".", ",", 1)
}

func DebugXMLName(c Currencies) string {
	return fmt.Sprintf("%s %s", c.Date.Format("2006-01-02"), c.Name)
}
