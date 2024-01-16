package qkit

import "encoding/json"

// Converts any type to a given type. If conversion fails, it returns the zero value of the given type.
func Cast[T any](val any) T {
	if val, ok := val.(T); ok {
		return val
	}
	var zero T
	return zero
}

// Converts any type to a given struct type.
func ToStruct[T any](m any) T {
	var value T
	bytes, _ := json.Marshal(m)
	json.Unmarshal(bytes, &value)
	return value
}

// Converts any type to a map[string]interface{}.
func ToMap(s any) map[string]interface{} {
	m := make(map[string]interface{})
	bytes, _ := json.Marshal(s)
	json.Unmarshal(bytes, &m)
	return m
}

// Generates a byte array from any type.
func ToBytes(v interface{}) []byte {
	bytes, _ := json.Marshal(v)
	return bytes
}

// Converts a byte array to a given type.
func FromBytes[T any](bytes []byte) T {
	var v T
	json.Unmarshal(bytes, &v)
	return v
}