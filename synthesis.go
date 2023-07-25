package qkit

func coalesce[T comparable](preference, fallback T) T {
	if preference == Zero[T]() {
		return fallback
	}

	return preference
}

// Coalesce returns its left-most value if it's not zero value
func Coalesce[T comparable](vv ...T) T {
	final := Zero[T]()
	for _, v := range vv {
		final = coalesce(final, v)
	}

	return final
}
