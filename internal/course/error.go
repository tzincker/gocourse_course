package course

import (
	"errors"
	"fmt"
)

var ErrNameRequired = errors.New("name is required")

type ErrNotFound struct {
	CourseId string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("course '%s' doesn't exist", e.CourseId)
}
