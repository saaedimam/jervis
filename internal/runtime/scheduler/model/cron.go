package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"github.com/ioriimasu/jervis/internal/runtime/scheduler/contracts"
)

// CronSchedule triggers a job based on a cron expression.
// Format: "min hour dom month dow"
type CronSchedule struct {
	expression string
	min        []int
	hour       []int
	dom        []int
	month      []int
	dow        []int
}

func NewCronSchedule(expr string) (contracts.Schedule, error) {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return nil, fmt.Errorf("invalid cron expression: expected 5 fields, got %d", len(fields))
	}

	s := &CronSchedule{expression: expr}
	var err error

	if s.min, err = parseField(fields[0], 0, 59); err != nil {
		return nil, err
	}
	if s.hour, err = parseField(fields[1], 0, 23); err != nil {
		return nil, err
	}
	if s.dom, err = parseField(fields[2], 1, 31); err != nil {
		return nil, err
	}
	if s.month, err = parseField(fields[3], 1, 12); err != nil {
		return nil, err
	}
	if s.dow, err = parseField(fields[4], 0, 6); err != nil {
		return nil, err
	}

	return s, nil
}

func parseField(field string, min, max int) ([]int, error) {
	if field == "*" {
		return nil, nil // wildcard
	}
	val, err := strconv.Atoi(field)
	if err != nil {
		return nil, fmt.Errorf("invalid cron field %q: only '*' and single numbers supported in Phase 1.5", field)
	}
	if val < min || val > max {
		return nil, fmt.Errorf("field value %d out of range [%d, %d]", val, min, max)
	}
	return []int{val}, nil
}

func (s *CronSchedule) NextRun(ref time.Time) time.Time {
	// Start searching from the next minute
	t := ref.Add(1 * time.Minute).Truncate(time.Minute)

	// Simple search: check each minute for up to 1 year
	for i := 0; i < 525600; i++ {
		if s.matches(t) {
			return t
		}
		t = t.Add(1 * time.Minute)
	}

	return time.Time{}
}

func (s *CronSchedule) matches(t time.Time) bool {
	if s.min != nil && t.Minute() != s.min[0] {
		return false
	}
	if s.hour != nil && t.Hour() != s.hour[0] {
		return false
	}
	if s.dom != nil && t.Day() != s.dom[0] {
		return false
	}
	if s.month != nil && int(t.Month()) != s.month[0] {
		return false
	}
	if s.dow != nil && int(t.Weekday()) != s.dow[0] {
		return false
	}
	return true
}
