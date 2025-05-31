package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	//  dataLen uint32(4) + type uint32(4)
	return 8
}
func (dp *DataPack) Pack(message ziface.IMessage) ([]byte, error) {
	databuf := bytes.NewBuffer([]byte{})
	//写入长度
	err := binary.Write(databuf, binary.LittleEndian, message.GetDataLen())
	if err != nil {
		return nil, err
	}
	//写入消息ID（类型）
	err = binary.Write(databuf, binary.LittleEndian, message.GetMsgID())
	if err != nil {
		return nil, err
	}
	err = binary.Write(databuf, binary.LittleEndian, message.GetData())
	if err != nil {
		return nil, err
	}
	return databuf.Bytes(), nil
}
func (dp *DataPack) Unpack(data []byte) (ziface.IMessage, error) {
	databuf := bytes.NewReader(data)
	msg := &Message{}
	err := binary.Read(databuf, binary.LittleEndian, &msg.dataLen)
	if err != nil {
		return nil, err
	}
	err = binary.Read(databuf, binary.LittleEndian, &msg.ID)
	if err != nil {
		return nil, err
	}

	if msg.dataLen > utils.GlobalObject.MaxPacketSize && utils.GlobalObject.MaxPacketSize > 0 {
		return nil, errors.New("too Large msg data recv!")
	}
	return msg, nil
}
