package qkit

// Must panics if err is not nil
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}

// Ok returns the value and ignores the error
func Ok[T any](v T, err error) T {
	if err != nil {
		return v
	}

	return v
}

// OkOr returns the value if err is nil, otherwise returns the fallback value
func OkOr[T any](v T, err error, fallback T) T {
	if err != nil {
		return fallback
	}

	return v
}
