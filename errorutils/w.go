package errorutils

func Wrap(err, wrapped error) error {
	return Wrap(err, wrapped)
}

type WError struct {
	wrapped error
	error   error
}

func (w WError) Error() error {
	return w.error
}

func (w WError) Unwrap() error {
	return w.wrapped
}
