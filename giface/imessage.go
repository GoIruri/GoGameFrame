package giface

// IMessage 将请求的消息封装到一个IMessage中，定义抽象的接口
type IMessage interface {
	GetMsgId() uint32  // 获取消息Id
	GetMsgLen() uint32 // 获取消息的长度
	GetData() []byte   // 获取消息内容
	SetMsgId(uint32)   // 设置消息Id
	SetMsgLen(uint32)  // 设置消息长度
	SetData([]byte)    // 设置消息内容
}
