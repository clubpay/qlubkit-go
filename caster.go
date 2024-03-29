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

// Converts any type to a given type based on their json representation. It partially fills the target in case they are not directly compatible.
func CastJSON[T any](val any) T {
	return FromJSON[T](ToJSON(val))
}

// Converts a given value to a byte array.
func ToJSON(val any) []byte {
	return Ok(json.Marshal(val))
}

// Converts a byte array to a given type.
func FromJSON[T any](bytes []byte) T {
	var v T
	json.Unmarshal(bytes, &v)
	return v
}

// Converts any type to a map[string]interface{}.
func ToMap(s any) map[string]interface{} {
	m := make(map[string]interface{})
	json.Unmarshal(Ok(json.Marshal(s)), &m)
	return m
}
