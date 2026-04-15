package handlers

import "errors"

var ErrInvalidEndDateFormat error = errors.New("invalid end_date format, expected YYYYMMDD")
