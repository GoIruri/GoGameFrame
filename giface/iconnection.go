package giface

import "net"

// IConnection 定义链接模块的抽象层
type IConnection interface {
	Start()

	Stop()
	//	获取当前链接绑定的socket conn
	GetTcpConnection() *net.TCPConn
	//	获取当前模块的链接ID
	GetConnID() uint32
	//	获取远程客户端的TCP状态 IP port
	RemoteAddr() net.Addr
	//	发送数据，将数据发送给远程的客户端
	SendMsg(msgId uint32, data []byte) error
}

// HandleFunc 定义一个处理链接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
