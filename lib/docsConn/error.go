package docsConn

import (
	"fmt"
)

type DocsError struct {
	s          string
	StatusCode int
}

func (e *DocsError) Error() string {
	return e.s
}

func newDocsError(StatusCode int, str string, args ...interface{}) *DocsError {
	var error DocsError
	error.s = fmt.Sprintf(str, args...)
	error.StatusCode = StatusCode
	return &error
}
