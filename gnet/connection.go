package gnet

import (
	"fmt"
	"net"
	"zinx/giface"
)

// Connection 链接
type Connection struct {
	//	当前链接的socket TCP套接字
	Conn *net.TCPConn
	//	链接的ID
	ConnID uint32
	//	当前的链接状态
	isClosed bool
<<<<<<< HEAD
=======
	//	当前绑定所绑定的处理业务的方法
	handleAPI giface.HandleFunc
>>>>>>> a097f1c616f3b3e12b2139584d06eab8616b0a8b
	//	告知当前链接已经退出的停止channel
	ExitChan chan bool
	//	该链接处理的方法Router
	Router giface.IRouter
}

// StartReader 链接的读业务的方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, "Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//	读取客户端的数据到buf中,最大是512字节
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err", err)
			continue
		}

		//	得到当前conn的Request请求数据
		req := Request{
			conn: c,
			data: buf,
		}
		//	从路由中找到注册绑定的Conn对应的Router调用
		go func(request giface.IRequest) {
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start... ConnID = ", c.ConnID)

	//	启动从当前链接的读数据的业务
	go c.StartReader()
	//	todo 启动从当前链接写数据的业务
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop... ConnID = ", c.ConnID)

	//	如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//	关闭socket链接
	c.Conn.Close()

	close(c.ExitChan)
}

func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) error {
	return nil
}

// NewConnection 初始化链接模块的方法
<<<<<<< HEAD
func NewConnection(conn *net.TCPConn, connID uint32, router giface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool),
=======
func NewConnection(conn *net.TCPConn, connID uint32, callbackApi giface.HandleFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		handleAPI: callbackApi,
		isClosed:  false,
		ExitChan:  make(chan bool),
>>>>>>> a097f1c616f3b3e12b2139584d06eab8616b0a8b
	}
	return c
}
