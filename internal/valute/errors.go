package valute

import "errors"

var (
	ErrNotFound       = errors.New("report not found")
	ErrInvalidDateReq = errors.New("invalid date_req format, expected DD/MM/YYYY")
)
