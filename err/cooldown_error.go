package err

type CooldownError struct {
	TimeLeft int64
	Message  string
}

func (e *CooldownError) Error() string {
	return e.Message
}

func (e *CooldownError) GetTimeLeft() int64 {
	return e.TimeLeft
}
