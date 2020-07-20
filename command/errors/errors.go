package errors

type ErrInvalidArgs struct{ Msg string }

func (e *ErrInvalidArgs) Error() string { return e.Msg }
