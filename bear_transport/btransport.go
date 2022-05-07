package bear_transport

import "io"

type BearTransport interface {
	io.ReadWriter
	io.ByteReader
	io.ByteWriter
	stringWriter
	Flusher
	ReadSizeProvider
}

type Flusher interface {
	Flush() (err error)
}

type ReadSizeProvider interface {
	RemainingBytes() (num_bytes uint64)
}

type stringWriter interface {
	WriteString(s string) (n int, err error)
}