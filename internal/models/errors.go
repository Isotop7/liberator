package models

import (
	"errors"
)

var ErrNoRecord = errors.New("models: no matching record found")
var ErrSumPageCount = errors.New("models: error while querying page count")
