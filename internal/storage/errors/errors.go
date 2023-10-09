package errors

import (
	"fmt"
)

var ErrIncorrectNewNote = fmt.Errorf("incorrect fields in new note")
