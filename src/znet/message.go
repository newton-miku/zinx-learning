package znet

type Message struct {
	ID uint32

	dataLen uint32
	data    []byte
}

func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		ID:      id,
		dataLen: uint32(len(data)),
		data:    data,
	}
}
func (m *Message) GetMsgID() uint32 {
	return m.ID
}
func (m *Message) GetDataLen() uint32 {
	return m.dataLen
}
func (m *Message) GetData() []byte {
	return m.data
}
func (m *Message) SetMsgID(id uint32) {
	m.ID = id
}
func (m *Message) SetData(data []byte) {
	m.data = data
}
func (m *Message) SetDataLen(length uint32) {
	m.dataLen = length
}
