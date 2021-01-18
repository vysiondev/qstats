package err

type ReadError struct {
	Message string
}

func (e *ReadError) Error() string {
	return e.Message
}
