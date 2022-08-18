package gnet

import "zinx/giface"

type Request struct {
	//	已经和客户端建立好的链接
	conn giface.IConnection
	//	客户端请求的数据
	data []byte
}

func (r *Request) GetConnection() giface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.data
}
