package asset

import (
	"errors"
)

var (
	SortTypeID        = SortType("id")
	SortTypeName      = SortType("name")
	SortTypeSize      = SortType("size")
	SortTypeCreatedat = SortType("createdat")

	ErrInvalidSortType = errors.New("invalid sort type")
)

type SortType string
