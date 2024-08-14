package utils

import "bytes"

type ReadCloserSpy struct {
	IsClosed bool
	Data     *bytes.Buffer
}

func NewReadCloserSpy() *ReadCloserSpy {
	return &ReadCloserSpy{
		IsClosed: false,
		Data:     bytes.NewBufferString(""),
	}
}

func (r *ReadCloserSpy) Read(p []byte) (n int, err error) {
	return r.Data.Read(p)
}

func (r *ReadCloserSpy) Close() error {
	r.IsClosed = true
	return nil
}
