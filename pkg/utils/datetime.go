package utils

import (
	"strings"
	"time"
)

func ParseQueryTime(v string) (*time.Time, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil, nil
	}

	layouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02",
	}

	var lastErr error
	for _, l := range layouts {
		t, err := time.Parse(l, v)
		if err == nil {
			if l == "2006-01-02" {
				t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
			}
			return &t, nil
		}
		lastErr = err
	}
	return nil, lastErr
}
