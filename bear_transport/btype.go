package bear_transport

type BearMessageType uint32

const (
	INVALID_TMESSAGE_TYPE BearMessageType = 0
	// support two types
	CALL                  BearMessageType = 1
	REPLY                 BearMessageType = 2

	EXCEPTION             BearMessageType = 3
	ONEWAY                BearMessageType = 4
	STREAMING			  BearMessageType = 5
)

// 目前支持call和reply两种模式 对异常并没有捕捉
func TypeToUint32(messageType BearMessageType) uint32 {
	if messageType == CALL {
		return 1
	} else if messageType == REPLY {
		return 2
	}
	return 0
}

func Uint32ToType(index uint32) BearMessageType {
	if index == 1 {
		return CALL
	} else if index == 2 {
		return REPLY
	}
	return INVALID_TMESSAGE_TYPE
}


