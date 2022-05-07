package bear_transport

import (
	"encoding/binary"
	"io"
)

type buffer struct {
	buf      []byte
	readIdx  int
	writeIdx int
}

// 最小buffer默认为256字节
func NewBuffer() *buffer {
	return &buffer{
		buf: make([]byte, 256),
	}
}

func (b *buffer) Read(reader io.Reader) {
	reader.Read(b.buf)
}

func (b *buffer) Write(writer io.Writer) {
	writer.Write(b.buf)
}

func (b *buffer) Length() int {
	return len(b.buf)
}
func (b *buffer) Next(n int) (buf []byte, err error) {
	buf, err = b.Peek(n)
	b.readIdx += n
	return buf, err
}

func (b *buffer) Peek(n int) (buf []byte, err error) {
	return b.buf[b.readIdx : b.readIdx+n], nil
}

func (b *buffer) Skip(n int) (err error) {
	b.readIdx += n
	return nil
}

func (b *buffer) Malloc(n int) (buf []byte, err error) {
	cur := b.writeIdx
	b.writeIdx += n
	return b.buf[cur: b.writeIdx], nil
}

// 判断buf是否足够 不够的花 double成长一次
func (b *buffer) IsNeedGrowup() {

}

func (b *buffer) ReadUint32() (uint32, error) {
	arr, err := b.Next(4)
	return binary.BigEndian.Uint32(arr), err
}

func (b *buffer) PeekUint32() (uint32, error) {
	arr, err := b.Peek(4)
	return binary.BigEndian.Uint32(arr), err
}

func (b *buffer) ReadUint16() (uint16, error) {
	arr, err := b.Next(2)
	return binary.BigEndian.Uint16(arr), err
}

func (b *buffer) PeekUint16() (uint16, error) {
	arr, err := b.Peek(2)
	return binary.BigEndian.Uint16(arr), err
}

func (b *buffer) WriteUint32(val uint32) error {
	buf, err := b.Malloc(4)
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint32(buf, val)
	return nil
}

func (b *buffer) WriteUint16(val uint16) error {
	buf, err := b.Malloc(4)
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint16(buf, val)
	return nil
}

func (b *buffer) WriteByte(val byte) error {
	buf, err := b.Malloc(1)
	if err != nil {
		return err
	}
	buf[0] = val
	return nil
}

func (b *buffer) WriteBinary(p []byte) error {
	n := len(p)
	arr, err := b.Malloc(n)
	copy(arr, p)
	b.writeIdx += n
	return err
}

// 关于string的处理 先写进len, 在读出对应大小的size
func (b *buffer) ReadString() (string, error) {
	strLen, err := b.ReadUint32()
	if err != nil {
		return "", err
	}
	arr, err := b.Next(int(strLen))
	return string(arr), err
}

func (b *buffer) WriteString(val string) error {
	strLen := len(val)
	b.WriteUint32(uint32(strLen))
	arr, err := b.Malloc(strLen)
	if err != nil {
		return err
	}
	copy(arr, val)
	return nil
}