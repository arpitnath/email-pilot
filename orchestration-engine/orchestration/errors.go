package orchestration

import "fmt"

type EngineError struct {
	Message string
}

func (e *EngineError) Error() string {
	return fmt.Sprintf("EngineError: %s", e.Message)
}

func NewEngineError(message string) *EngineError {
	return &EngineError{Message: message}
}
