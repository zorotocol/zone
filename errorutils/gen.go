package errorutils

func Throw(err error) {
	if err != nil {
		panic(err)
	}
}
func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
