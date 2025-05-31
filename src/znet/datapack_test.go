package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 测试DataPack  封包、拆包 单元测试
func TestDataPack(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("[ERROR] Listener err: ", err)
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("[ERROR] Accept err:", err)
			}
			go func(conn net.Conn) {
				dp := NewDataPack()
				for {
					headPack := make([]byte, dp.GetHeadLen())
					io.ReadFull(conn, headPack)
					msgHead, err := dp.Unpack(headPack)
					if err != nil && err != io.EOF {
						fmt.Println("[ERROR] datapack head Unpack err:", err)
						return
					}
					if msgHead.GetDataLen() > 0 {
						msg := msgHead.(*Message)
						msg.data = make([]byte, msg.GetDataLen())

						//读取data数据部分
						_, err := io.ReadFull(conn, msg.data)
						if err != nil {
							fmt.Println("[ERROR] datapack data ReadFull err:", err)
							break
						}
						fmt.Println("------>[INFO] Recv msg ID:", msg.GetMsgID(), "DataLen:", msg.GetDataLen(), "Data:", string(msg.GetData()))
					}
				}
			}(conn)
		}
	}()

	client, clerr := net.Dial("tcp", "127.0.0.1:7777")
	if clerr != nil {
		fmt.Println("[ERROR] Dial err:", clerr)
	}
	dp := NewDataPack()
	// 模拟TCP粘包
	msg1 := &Message{
		ID:      1,
		dataLen: 5,
		data:    []byte{'H', 'e', 'l', 'l', 'o'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("[ERROR] datapack Pack1 err:", err)
	}
	//第二个包
	msg2 := &Message{
		ID:      2,
		dataLen: 7,
		data:    []byte{'n', 'i', ',', 'h', 'a', 'o', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("[ERROR] datapack Pack2 err:", err)
	}
	sendData1 = append(sendData1, sendData2...)
	client.Write(sendData1)

	select {}
}
