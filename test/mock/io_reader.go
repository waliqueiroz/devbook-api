package mock

import "errors"

type reader struct{}

func NewReader() *reader {
	return &reader{}
}

func (r reader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}
