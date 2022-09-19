package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/gnet"
)

/*
模拟客户端
*/
func main() {
	fmt.Println("Client start ...")

	time.Sleep(1 * time.Second)

	//	1直接连接远程服务器,得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}

	for {
		// 发送封包的message消息,MsgId:0
		dp := gnet.NewDataPack()
		binaryMsg, err := dp.Pack(gnet.NewMessage(0, []byte("gamev0.5 client test message")))
		if err != nil {
			fmt.Println("pack error", err)
			return
		}
		_, err = conn.Write(binaryMsg)
		if err != nil {
			fmt.Println("write error ", err)
			return
		}

		// 服务器就应该给我们回复一个message数据，MsgID：1 pingpingping

		// 先读取流中的head部分 得到ID和DataLen
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error ", err)
			return
		}
		// 将二进制的head拆包到msg结构体中
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msgHead error", err)
			break
		}
		if msgHead.GetMsgLen() > 0 {
			// 再根据DataLen进行第二次读取，将Data读出来
			msg := msgHead.(*gnet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error,", err)
				return
			}

			fmt.Println("---->Recv Server Msg : ID = ", msg.Id, ", Len = ", msg.DataLen, ", Data = ", string(msg.Data))
		}

		//	cpu阻塞
		time.Sleep(1 * time.Second)
	}
}
