package bear_transport

import (
	"io"
)


// 直接通过pb做序列化
// WriteStruct writes the given Thrift struct to a writer. It pools TProtocols.
func WriteStruct(writer io.Writer, s pbStruct) error {
	// 序列化
	res, err := s.XXX_Marshal([]byte{}, true)
	writer.Write(res)
	return err
}

// ReadStruct reads the given Thrift struct. It pools TProtocols.
func ReadStruct(reader io.Reader, s pbStruct) error {
	var res []byte
	reader.Read(res)
	err := s.XXX_Unmarshal(res)
	return err
}
