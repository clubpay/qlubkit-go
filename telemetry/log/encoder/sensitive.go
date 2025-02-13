package encoder

import (
	"bytes"
	"errors"
	"regexp"

	qkit "github.com/clubpay/qlubkit-go"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var _ zapcore.Encoder = (*Sensitive)(nil)

type SensitiveConfig struct {
	zapcore.EncoderConfig

	ForbiddenHeaders []string
	AutoDetect       bool
}

type Sensitive struct {
	zapcore.Encoder
	forbiddenHeaders [][]byte
}

func NewSensitive(cfg SensitiveConfig) *Sensitive {
	cfg.NewReflectedEncoder = newJSONEncoder
	return &Sensitive{
		Encoder:          zapcore.NewJSONEncoder(cfg.EncoderConfig),
		forbiddenHeaders: qkit.Map(func(src string) []byte { return []byte(src) }, cfg.ForbiddenHeaders),
	}
}

const (
	minPayloadSize = 128
)

var (
	requestRegex  = regexp.MustCompile(`^(GET|POST|PUT|DELETE|PATCH|OPTIONS|HEAD) [^ ]+`)
	responseRegex = regexp.MustCompile(`^HTTP/\d\.\d \d+ [^\r\n]+`)
)

func (s Sensitive) Clone() zapcore.Encoder {
	return Sensitive{Encoder: s.Encoder.Clone()}
}

func (s Sensitive) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	return s.Encoder.EncodeEntry(entry, fields)
}

func (s Sensitive) sanitizeHTTPRequest(data []byte) []byte {
	nextID := s.parseRequestLine(data)
	if nextID < 0 {
		return data
	}

	s.parseHeaders(data, nextID)

	return data
}

func (s Sensitive) parseRequestLine(src []byte) int {
	lineIdx, nextIdx, err := splitLine(src, 0)
	line := src[:lineIdx]
	if err != nil {
		return -1
	}
	MethodIndex := bytes.IndexByte(line, ' ')
	if MethodIndex < 0 || MethodIndex < 3 {
		return -1
	}
	URIIndex := bytes.IndexByte(line[MethodIndex+1:], ' ')
	if URIIndex < 0 {
		return -1
	}

	return nextIdx
}

func (s Sensitive) parseHeaders(src []byte, next int) int {
	var (
		err  error
		line int
	)

	for {
		start := next
		line, next, err = splitLine(src, next)
		if err != nil {
			return -1
		}
		if line > start {
			break
		}

		s.parseHeaderLine(src[start:line])
	}

	return next
}

func (s Sensitive) parseHeaderLine(src []byte) {
	idx := bytes.IndexByte(src, ':')
	if idx < 0 {
		return
	}
	// RFC2616 Section 4.2
	// Remove all leading and trailing LWS on field contents

	// skip leading LWS
	var i int = idx + 1
	for ; i < len(src); i++ {
		if src[i] != ' ' && src[i] != '\t' {
			break
		}
	}
	// skip trailing LWS
	var j int = len(src) - 1
	for ; j > i; j-- {
		if src[j] != ' ' && src[j] != '\t' {
			break
		}
	}

	for _, h := range s.forbiddenHeaders {
		if bytes.EqualFold(h, src[:idx]) {
			copy(src[i:j+1], bytes.Repeat([]byte{'x'}, j-i+1))

			break
		}
	}

	return
}

func splitLine(src []byte, i int) (line, rest int, err error) {
	idx := bytes.IndexByte(src[i:], '\n') + i
	if idx < 1 { // 0: cr 1: lf
		return 0, 0, errors.New("buffer too small")
	}

	if src[idx-1] == '\r' {
		return idx - 1, idx + 1, nil
	}

	return idx, idx + 1, nil
}

func (s Sensitive) sanitizeHTTPResponse(value []byte) []byte { return value }
