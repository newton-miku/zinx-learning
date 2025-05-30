package znet

type Message struct {
	ID uint32

	DataLen uint32
	Data    []byte
}

func (m *Message) GetMsgID() uint32 {
	return m.ID
}
func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}
func (m *Message) GetData() []byte {
	return m.Data
}
func (m *Message) SetMsgID(id uint32) {
	m.ID = id
}
func (m *Message) SetData(data []byte) {
	m.Data = data
}
func (m *Message) SetDataLen(length uint32) {
	m.DataLen = length
}
