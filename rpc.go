package h2

import "io"

type Codec struct {
	i io.Reader
	o io.Writer
}

func (c Codec) Read(b []byte) (n int, err error) {
	return c.i.Read(b)
}

func (c Codec) Write(b []byte) (n int, err error) {
	return c.o.Write(b)
}

func (Codec) Close() error {
	return nil
}
