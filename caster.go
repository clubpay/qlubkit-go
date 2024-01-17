package qkit

import (
	"encoding/json"
)

// Converts any type to a given type. If conversion fails, it returns the zero value of the given type.
func Cast[T any](val any) T {
	if val, ok := val.(T); ok {
		return val
	}
	var zero T
	return zero
}

// Converts any type to a given type. This partially fills the target in case the 2 types are not directly compatible.
func CastPartial[T any](val any) T {
	return FromBytes[T](Ok(json.Marshal(val)))
}

// Converts any type to a map[string]interface{}.
func ToMap(s any) map[string]interface{} {
	m := make(map[string]interface{})
	json.Unmarshal(Ok(json.Marshal(s)), &m)
	return m
}

// Converts a byte array to a given type.
func FromBytes[T any](bytes []byte) T {
	var v T
	json.Unmarshal(bytes, &v)
	return v
}