package mock

import "errors"

type readerMock struct{}

func NewReader() *readerMock {
	return &readerMock{}
}

func (r readerMock) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}
