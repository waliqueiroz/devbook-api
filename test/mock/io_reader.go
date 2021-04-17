package mock

import "errors"

type errReader struct{}

// NewUserRepository creates a new user repository
func NewReader() *errReader {
	return &errReader{}
}

func (e errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}
