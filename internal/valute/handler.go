package valute

import (
	"errors"
	"net/http"
	"time"
)

const dateFormat = "02/01/2006"

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dateReq := r.URL.Query().Get("date_req")
	if dateReq == "" {
		http.Error(w, "date_req is required", http.StatusBadRequest)
		return
	}

	dateParsed, err := time.Parse(dateFormat, dateReq)
	if err != nil {
		http.Error(w, ErrInvalidDateReq.Error(), http.StatusBadRequest)
		return
	}

	report, err := h.service.GetReportByDate(r.Context(), dateParsed)
	if err != nil {
		if errors.Is(err, ErrInvalidDateReq) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	xmlReport, err := MarshalXMLDaily(report)
	if err != nil {
		http.Error(w, "marshal xml error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(xmlReport)
}
