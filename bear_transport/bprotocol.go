package bear_transport

// 目前是参考thrift的设计

type BearProtocol interface {
	WriteMessageBegin(name string, typeId BearMessageType, seqid int32) error
	WriteMessageEnd() error
	WriteStructBegin(name string) error
	WriteStructEnd() error
	WriteFieldBegin(name string, typeId BearType, id int16) error
	WriteFieldEnd() error
	WriteFieldStop() error
	WriteMapBegin(keyType BearType, valueType BearType, size int) error
	WriteMapEnd() error
	WriteListBegin(elemType BearType, size int) error
	WriteListEnd() error
	WriteSetBegin(elemType BearType, size int) error
	WriteSetEnd() error
	WriteBool(value bool) error
	WriteByte(value int8) error
	WriteI16(value int16) error
	WriteI32(value int32) error
	WriteI64(value int64) error
	WriteDouble(value float64) error
	WriteString(value string) error
	WriteBinary(value []byte) error

	ReadMessageBegin() (name string, typeId BearMessageType, seqid int32, err error)
	ReadMessageEnd() error
	ReadStructBegin() (name string, err error)
	ReadStructEnd() error
	ReadFieldBegin() (name string, typeId BearType, id int16, err error)
	ReadFieldEnd() error
	ReadMapBegin() (keyType BearType, valueType BearType, size int, err error)
	ReadMapEnd() error
	ReadListBegin() (elemType BearType, size int, err error)
	ReadListEnd() error
	ReadSetBegin() (elemType BearType, size int, err error)
	ReadSetEnd() error
	ReadBool() (value bool, err error)
	ReadByte() (value int8, err error)
	ReadI16() (value int16, err error)
	ReadI32() (value int32, err error)
	ReadI64() (value int64, err error)
	ReadDouble() (value float64, err error)
	ReadString() (value string, err error)
	ReadBinary() (value []byte, err error)

	Skip(fieldType BearType) (err error)
	Flush() (err error)

	Transport() BearTransport
}

