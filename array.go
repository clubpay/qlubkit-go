package qkit

func Filter[T any](match func(src T) bool, tt []T) []T {
	ftt := make([]T, 0, len(tt))
	for _, t := range tt {
		if match(t) {
			ftt = append(ftt, t)
		}
	}

	return ftt[:len(ftt):len(ftt)]
}

func Map[S, D any](transformer func(src S) D, ss []S) []D {
	dd := make([]D, len(ss))
	for k, src := range ss {
		dd[k] = transformer(src)
	}

	return dd
}

func Reduce[T, R any](reducer func(r R, t T) R, tt []T) R {
	var r R
	for _, t := range tt {
		r = reducer(r, t)
	}

	return r
}

func Paginate[T any](arr []T, pageSize int, fn func(start, end int) error) error {
	start := 0
	for {
		end := start + pageSize - 1
		if end > len(arr)-1 {
			end = len(arr) - 1
		}
		err := fn(start, end)
		if err != nil {
			return err
		}
		start = end + 1
		if start >= len(arr) {
			break
		}
	}

	return nil
}
