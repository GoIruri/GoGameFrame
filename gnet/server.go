package gnet

import (
	"errors"
	"fmt"
	"net"
	"zinx/giface"
)

// Server IServer接口的实现,定义一个Server的服务器模块
type Server struct {
	//	服务器名称
	Name string
	//	服务器绑定的IP版本
	IPVersion string
	//	服务器监听的IP
	IP string
	//	服务器监听的端口
	Port int
	//	路由,Server注册的链接对应的处理业务
	Router giface.IRouter
}

// CallBackToClient 定义当前客户端链接的所绑定handle api(目前这个handle是写死的,以后优化应用应该由用户自定义handle方法)
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	//	回显的业务
	fmt.Println("[Conn Handle] callbacktoclient ...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err", err)
		return errors.New("CallBackToClient error")
	}

	return nil
}

// Start 启动服务器
func (s *Server) Start() {
	fmt.Println("start server ...")

	go func() {
		//	1获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}

		//	2监听服务端的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "error", err)
			return
		}
		fmt.Println("start Zinx server succ", s.Name, "listenning...")

		var cid uint32 = 0
		//	3阻塞的等待客户端连接,处理客户端的业务
		for {
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			//将处理新链接的业务方法 和 conn 进行绑定,得到我们的链接模块
			dealConn := NewConnection(conn, cid, CallBackToClient)
			cid++

			//	启动当前的链接业务
			go dealConn.Start()
		}
	}()
}

// Stop 停止服务器
func (s *Server) Stop() {
	//todo 将一些服务器的资源,状态或者一些已经开辟的链接信息,进行停止或者回收
}

// Serve 运行服务器
func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	//todo 做一些启动服务器之后的额外业务

	//	阻塞状态
	select {}
}

// AddRouter 添加路由功能
func (s *Server) AddRouter(router giface.IRouter) {

}

// NewServer 初始化
func NewServer(name string) giface.Iserver {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "127.0.0.1",
		Port:      8999,
	}
	return s
}
