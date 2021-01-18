package err

type SafeError struct {
	Message string
}

func (e *SafeError) Error() string {
	return e.Message
}
