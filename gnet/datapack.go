package gnet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/giface"
	"zinx/utils"
)

// DataPack 封包，拆包模块
type DataPack struct {
}

func (dp *DataPack) GetHeadLen() uint32 {
	// DataLen uint32 (4字节) + ID uint32 (4字节)
	return 8
}

// Pack 封包方法(dataLen/msgID/data)
func (dp *DataPack) Pack(msg giface.IMessage) ([]byte, error) {
	// 创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 将dataLen写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	// 将MsgId写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	// 将data数据写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// Unpack 拆包方法，只需要将包的Head信息读出来就可以了，再根据Head的信息里的data的长度再进行一次读
func (dp *DataPack) Unpack(binaryBuff []byte) (giface.IMessage, error) {
	// 创建一个从输入读取二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryBuff)

	// 只解压Head信息，得到dataLen和MsgID
	msg := &Message{}

	// 读DataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// 读MsgId
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断dataLen是否已经超出了我们允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data recv!")
	}

	return msg, nil
}

func NewDataPack() *DataPack {
	return &DataPack{}
}
