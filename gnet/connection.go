package gnet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/giface"
	"zinx/utils"
)

// Connection 链接
type Connection struct {
	// 当前conn隶属于哪个server
	TcpServer giface.Iserver
	//	当前链接的socket TCP套接字
	Conn *net.TCPConn
	//	链接的ID
	ConnID uint32
	//	当前的链接状态
	isClosed bool
	//	告知当前链接已经退出的停止channel（由Reader告知Writer退出）
	ExitChan chan bool
	// 无缓冲的管道，用于读写goroutine之间的消息通信
	msgChan chan []byte
	//	消息的管理MsgID和对应的处理业务API关系
	MsgHandler giface.IMsgHandle
	// 链接属性集合
	property map[string]interface{}
	// 保护链接属性的锁
	protertyLock sync.RWMutex
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
		// 得到当前conn的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池机制，将消息发送给worker工作池处理即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 从路由中找到注册绑定的Conn对应的Router调用
			// 根据绑定好的MsgID找到对应处理api业务 执行
			go c.MsgHandler.DoMsgHandler(&req)
		}
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

	// 按照开发者传递进来的 创建连接之后需要调用的处理业务，执行对应Hook函数
	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop... ConnID = ", c.ConnID)

	//	如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	// 调用开发者注册的在销毁链接之前 需要执行的业务Hook函数
	c.TcpServer.CallOnConnStop(c)

	//	关闭socket链接
	c.Conn.Close()

	// 告知Writer关闭
	c.ExitChan <- true

	// 将当前链接从ConnMgr中移除
	c.TcpServer.GetConnMgr().Remove(c)

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

// SetProperty 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.protertyLock.Lock()
	defer c.protertyLock.Unlock()

	// 添加一个链接属性
	c.property[key] = value
}

// GetProperty 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.protertyLock.RLock()
	defer c.protertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	}
	return nil, errors.New("no property found")
}

// RemoveProperty 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.protertyLock.Lock()
	defer c.protertyLock.Unlock()

	delete(c.property, key)
}

// NewConnection 初始化链接模块的方法
func NewConnection(server giface.Iserver, conn *net.TCPConn, connID uint32, handle giface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: handle,
		isClosed:   false,
		msgChan:    make(chan []byte),
		ExitChan:   make(chan bool),
		property:   make(map[string]interface{}),
	}

	// 将conn加入到ConnMgr中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}
