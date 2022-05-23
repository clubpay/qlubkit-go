package qkit

import (
	"encoding/json"
	"fmt"
)

type Stringer struct {
	s   string
	n   json.Number
	isn bool // is number
}

func (s *Stringer) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &s.s)
	if err == nil {
		return nil
	}

	s.isn = true

	return json.Unmarshal(data, &s.n)
}

func (s Stringer) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", s.String())), nil
}

func (s *Stringer) String() string {
	if s.isn {
		return s.n.String()
	}

	return s.s
}

func StrToStringer(s string) Stringer {
	return Stringer{
		s: s,
	}
}
