// Copyright (c) 2015 Uber Technologies, Inc.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package pb

import (
	tchannel "github.com/bbbearxyz/another-tchannel-go"
	"github.com/bbbearxyz/another-tchannel-go/typed"
)

// 要求interface可以序列化和反序列化
type PbStruct interface {
	XXX_Unmarshal(b []byte) error
	XXX_Marshal(b []byte, deterministic bool) ([]byte, error)
}

// 直接通过pb做序列化
// WriteStruct writes the given pb struct to a writer.
func WriteStruct(w tchannel.ArgWriter, s PbStruct) error {
	// 序列化
	res, err := s.XXX_Marshal([]byte{}, true)
	writer := typed.NewWriter(w)
	writer.WriteUint16(uint16(len(res)))
	writer.WriteBytes(res)
	return err
}

// ReadStruct reads the given Thrift struct.
// 目前没有通过buf来做缓存
// 可能会出现half read
func ReadStruct(r tchannel.ArgReader, s PbStruct) error {
	reader := typed.NewReader(r)
	length := reader.ReadUint16()

	reader.Release()
	res := make([]byte, length)
	// there is some bug in buf
	num, _ := r.Read(res)
	if num != int(length) {
		println("want length(%d), but read(%d)", int(length), num)
	}
	err := s.XXX_Unmarshal(res)
	return err
}
