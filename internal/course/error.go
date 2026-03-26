package course

import (
	"errors"
	"fmt"
	"time"
)

var ErrNameRequired = errors.New("name is required")

type ErrNotFound struct {
	CourseId string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("course '%s' doesn't exist", e.CourseId)
}

type ErrEndDateNotValid struct {
	StartDate time.Time
	EndDate   time.Time
}

func (e *ErrEndDateNotValid) Error() string {
	start := e.StartDate.Format("2006-01-02")
	end := e.EndDate.Format("2006-01-02")
	return fmt.Sprintf("course start date '%s' is greater than %s", start, end)
}
