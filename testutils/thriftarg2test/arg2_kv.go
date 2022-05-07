package thriftarg2test

import (
	"fmt"
	"testing"

	"github.com/bbbearxyz/another-tchannel-go/typed"
	"github.com/stretchr/testify/require"
)

// BuildKVBuffer builds an thrift Arg2 KV buffer.
func BuildKVBuffer(kv map[string]string) []byte {
	// Scan once to know size of buffer
	var bufSize int
	for k, v := range kv {
		// k~2 v~2
		bufSize += 2 + len(k) + 2 + len(v)
	}
	bufSize += 2 // nh:2
	buf := make([]byte, bufSize)
	wb := typed.NewWriteBuffer(buf)
	wb.WriteUint16(uint16(len(kv)))
	for k, v := range kv {
		wb.WriteLen16String(k)
		wb.WriteLen16String(v)
	}
	return buf[:wb.BytesWritten()]
}

// ReadKVBuffer converts an arg2 buffer to a string map
func ReadKVBuffer(b []byte) (map[string]string, error) {
	rbuf := typed.NewReadBuffer(b)
	nh := rbuf.ReadUint16()
	retMap := make(map[string]string, nh)
	for i := uint16(0); i < nh; i++ {
		key := rbuf.ReadLen16String()
		val := rbuf.ReadLen16String()
		retMap[key] = val
	}
	if rbuf.BytesRemaining() > 0 {
		return nil, fmt.Errorf("kv buffer wasn't fully consumed (%d bytes remaining)", rbuf.BytesRemaining())
	}
	return retMap, nil
}

// MustReadKVBuffer calls require.NoError on the error returned by ReadKVBuffer
func MustReadKVBuffer(tb testing.TB, b []byte) map[string]string {
	m, err := ReadKVBuffer(b)
	require.NoError(tb, err)
	return m
}
