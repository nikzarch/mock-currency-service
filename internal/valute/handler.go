package valute

import (
	"log"
	"net/http"
	"time"
)

const dateFormat = "02/01/2006"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	reqDate := r.URL.Query()["date"]
	dateParsed, err := time.Parse(dateFormat, reqDate[0])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	currs, err := h.service.GetReportByDate(r.Context(), dateParsed)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	report, err := MarshalXMLDaily(currs)
	if err != nil {
		log.Printf("MarshalXMLDaily error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.Write(report)
}
