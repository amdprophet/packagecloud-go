package errors

type ErrInvalidArgs struct{ Msg string }

func (e *ErrInvalidArgs) Error() string { return e.Msg }

type ErrWithUsage struct {
	Msg     string
	UsageFn func() error
}

func NewErrorWithUsageFactory(usageFn func() error) func(string) error {
	return func(msg string) error {
		err := ErrWithUsage{
			Msg:     msg,
			UsageFn: usageFn,
		}
		return &err
	}
}

func (e *ErrWithUsage) Error() string { return e.Msg }

func (e *ErrWithUsage) Usage() { e.UsageFn() }
