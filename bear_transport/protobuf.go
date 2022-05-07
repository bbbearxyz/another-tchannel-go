package bear_transport

import (
	"context"
	"fmt"
)

// Protobuf Protocol
// Len 4Byte, MsgType 4Byte, MethodName String, SeqID 4Byte, Protobuf Body String

// 要求interface可以序列化和反序列化
type pbStruct interface {
	XXX_Unmarshal(b []byte) error
	XXX_Marshal(b []byte, deterministic bool) ([]byte, error)
}

func NewProtobufCodec() *protobufCodec {
	return &protobufCodec{}
}

type protobufCodec struct{}

type RPCInfo struct {
	MsgType BearMessageType
	MethodName string
	SeqId uint32
}
// content可以是args 也可以是result
func (c protobufCodec) Marshal(rpcInfo RPCInfo, content pbStruct, buf *buffer) error {
	// 保持4个byte的长度位置

	buf.WriteUint32(TypeToUint32(rpcInfo.MsgType))
	buf.WriteString(rpcInfo.MethodName)
	buf.WriteUint32(rpcInfo.SeqId)

	// encode pb struct
	msgBuf, err := content.XXX_Marshal([]byte{}, true)
	if err != nil {
		return err
	}
	buf.WriteBinary(msgBuf)
	return nil
}

func (c protobufCodec) Unmarshal(ctx context.Context, message remote.Message, in remote.ByteBuffer) error {
	payloadLen := message.PayloadLen()
	magicAndMsgType, err := codec.ReadUint32(in)
	if err != nil {
		return err
	}
	if magicAndMsgType&codec.MagicMask != codec.ProtobufV1Magic {
		return perrors.NewProtocolErrorWithType(perrors.BadVersion, "Bad version in protobuf Unmarshal")
	}
	msgType := magicAndMsgType & codec.FrontMask
	if err := codec.UpdateMsgType(msgType, message); err != nil {
		return err
	}

	methodName, methodFieldLen, err := codec.ReadString(in)
	if err != nil {
		return perrors.NewProtocolErrorWithErrMsg(err, fmt.Sprintf("protobuf unmarshal, read method name failed: %s", err.Error()))
	}
	if err = codec.SetOrCheckMethodName(methodName, message); err != nil && msgType != uint32(remote.Exception) {
		return err
	}
	seqID, err := codec.ReadUint32(in)
	if err != nil {
		return perrors.NewProtocolErrorWithErrMsg(err, fmt.Sprintf("protobuf unmarshal, read seqID failed: %s", err.Error()))
	}
	if err = codec.SetOrCheckSeqID(int32(seqID), message); err != nil && msgType != uint32(remote.Exception) {
		return err
	}
	actualMsgLen := payloadLen - metaInfoFixLen - methodFieldLen
	actualMsgBuf, err := in.Next(actualMsgLen)
	if err != nil {
		return perrors.NewProtocolErrorWithErrMsg(err, fmt.Sprintf("protobuf unmarshal, read message buffer failed: %s", err.Error()))
	}
	// exception message
	if message.MessageType() == remote.Exception {
		var exception pbError
		if err := exception.Unmarshal(actualMsgBuf); err != nil {
			return perrors.NewProtocolErrorWithMsg(fmt.Sprintf("protobuf unmarshal Exception failed: %s", err.Error()))
		}
		return &exception
	}

	if err = codec.NewDataIfNeeded(methodName, message); err != nil {
		return err
	}
	data := message.Data()
	// fast read
	if msg, ok := data.(bprotoc.FastRead); ok {
		_, err := bprotoc.Binary.ReadMessage(actualMsgBuf, bprotoc.SkipTypeCheck, msg)
		if err != nil {
			return remote.NewTransErrorWithMsg(remote.ProtocolError, err.Error())
		}
		return nil
	}
	msg, ok := data.(protobufMsgCodec)
	if !ok {
		return remote.NewTransErrorWithMsg(remote.InvalidProtocol, "decode failed, codec msg type not match with protobufCodec")
	}
	if err = msg.Unmarshal(actualMsgBuf); err != nil {
		return remote.NewTransErrorWithMsg(remote.ProtocolError, err.Error())
	}
	return err
}
