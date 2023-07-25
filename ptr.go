package qkit

// Zero returns zero value of generic type
func Zero[T any]() T {
	return *new(T)
}

func PtrVal[T any](src *T) T {
	if src == nil {
		return Zero[T]()
	}

	return *src
}

func ValPtr[T any](src T) *T {
	return &src
}

func String(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func StringPtr(in string) (out *string) {
	if in == "" {
		return nil
	}

	out = new(string)
	*out = in

	return
}

func Int64(s *int64) int64 {
	if s == nil {
		return 0
	}

	return *s
}

func Int64Ptr(in int64) (out *int64) {
	if in == 0 {
		return nil
	}

	out = new(int64)
	*out = in

	return
}

func BoolPtr(in bool) (out *bool) {
	if !in {
		return nil
	}

	out = new(bool)
	*out = in

	return
}

func BoolPtrStrict(in bool) (out *bool) {
	out = new(bool)
	*out = in

	return
}

func Int32Ptr(in int32) (out *int32) {
	if in == 0 {
		return
	}
	out = new(int32)
	*out = in

	return
}

func Uint64Ptr(in uint64) (out *uint64) {
	if in == 0 {
		return
	}
	out = new(uint64)
	*out = in

	return
}

func Uint32Ptr(in uint32) (out *uint32) {
	if in == 0 {
		return
	}
	out = new(uint32)
	*out = in

	return
}

func Int(in *int) int {
	if in == nil {
		return 0
	}

	return *in
}

func IntPtr(in int) (out *int) {
	if in == 0 {
		return
	}
	out = new(int)
	*out = in

	return
}

func UInt(in *uint) uint {
	if in == nil {
		return 0
	}

	return *in
}

func UintPtr(in uint) (out *uint) {
	if in == 0 {
		return
	}
	out = new(uint)
	*out = in

	return
}

func Float64(f *float64) float64 {
	if f == nil {
		return 0
	}

	return *f
}

func Float64Ptr(f float64) *float64 {
	if f == 0 {
		return nil
	}

	return &f
}

func Float64PtrStrict(f float64) *float64 {
	return &f
}
