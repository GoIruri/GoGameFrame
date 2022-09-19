package gnet

type Message struct {
	Id      uint32 // 消息Id
	DataLen uint32 // 消息长度
	Data    []byte // 消息内容
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}

func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

func (m *Message) SetMsgLen(length uint32) {
	m.DataLen = length
}

func (m *Message) SetData(bytes []byte) {
	m.Data = bytes
}

func NewMessage(id uint32, data []byte) *Message {
	return &Message{Id: id, DataLen: uint32(len(data)), Data: data}
}
