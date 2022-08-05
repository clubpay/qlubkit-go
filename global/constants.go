package global

import "time"

/*
	Globally accessible constants should be defined here
*/

const (
	DefaultIdempotancyTtl = 24 * time.Hour
)
