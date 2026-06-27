package pagination

import "strconv"

const (
	// DefaultLimit is the default limit for pagination
	DefaultLimit = 10
	// MaxLimit is the maximum limit for pagination
	MaxLimit = 100
)

// Parse converts limit and offset query strings into validated pagination values.
func Parse(limitStr, offsetStr string) (limit, offset int) {
	limit = DefaultLimit
	offset = 0

	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
		offset = o
	}

	return limit, offset
}
