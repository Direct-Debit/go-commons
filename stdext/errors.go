package stdext

import "fmt"

// WrapError wraps a given error with a message and returns it.
// If the given error is nil, nil is also returned.
// The message may be a format string with the arguments specified in a.
func WrapError(err error, message string, a ...any) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", fmt.Sprintf(message, a...), err)
}

type Monad[R any] struct {
	Result R
	Error  error
}

func (m *Monad[R]) Success() bool {
	return m.Error == nil
}

func (m *Monad[R]) ValueSafe() R {
	if m.Error != nil {
		var r R
		return r
	}
	return m.Result
}
