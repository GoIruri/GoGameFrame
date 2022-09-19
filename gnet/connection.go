package gnet

import (
	"errors"
	"fmt"
	"io"
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
	//	告知当前链接已经退出的停止channel（由Reader告知Writer退出）
	ExitChan chan bool
	// 无缓冲管道，用于读写goroutine之间的消息通信
	msgChan chan []byte
	//	消息的管理MsgID和对应的处理业务API关系
	MsgHandler giface.IMsgHandle
}

// StartReader 链接的读业务的方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running...")
	defer fmt.Println("connID = ", c.ConnID, "Reader is exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//	读取客户端的数据到buf中
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("recv buf err", err)
		//	continue
		//}
		// 创建一个拆包解包对象
		dp := NewDataPack()

		// 读取客户端的Msg Head 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTcpConnection(), headData); err != nil {
			fmt.Println("read MsgHead error", err)
			break
		}
		// 拆包，得到MsgID 和 MsgDataLen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}
		// 根据datalen，再次读取Data，放在msg.Data属性中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTcpConnection(), data); err != nil {
				fmt.Println("read msgdata error", err)
				break
			}
		}
		msg.SetData(data)
		//	得到当前conn的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		//	从路由中找到注册绑定的Conn对应的Router调用
		// 根据绑定好的MsgID找到对应处理api业务 执行
		go c.MsgHandler.DoMsgHandler(&req)
	}
}

// StartWriter 写消息goroutine，专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("writer goroutine is running...")
	defer fmt.Println(c.RemoteAddr().String(), "conn writer exit...")

	// 不断的阻塞的等待channel的消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			_, err := c.Conn.Write(data)
			if err != nil {
				fmt.Println("send data error ", err)
				return
			}
		case <-c.ExitChan:
			// 代表Reader已经退出，此时Writer也要退出
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start... ConnID = ", c.ConnID)

	//	启动从当前链接的读数据的业务
	go c.StartReader()
	//	启动从当前链接写数据的业务
	go c.StartWriter()
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

	// 告知Writer关闭
	c.ExitChan <- true

	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
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

// SendMsg 提供一个SendMsg方法，将我们要发送给客户端的数据，先进行封包，再进行发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}

	// 将data进行封包 MsgDataLen/MsgID/Data
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg")
	}

	// 将数据发送给客户端
	c.msgChan <- binaryMsg

	return nil
}

// NewConnection 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, handle giface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: handle,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool),
	}
	return c
}
