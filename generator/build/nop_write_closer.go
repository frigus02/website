package build

import "io"

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error {
	return nil
}

func newNopWriteCloser(w io.Writer) io.WriteCloser {
	return nopWriteCloser{w}
}
