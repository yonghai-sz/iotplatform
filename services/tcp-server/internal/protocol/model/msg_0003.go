package model

type Msg0003 struct {
	Header *MsgHeader `json:"header"`
}

func (m *Msg0003) Decode(packet *PacketData) error {
	m.Header = packet.Header
	return nil
}

func (m *Msg0003) GetHeader() *MsgHeader {
	return m.Header
}

func (m *Msg0003) GenOutgoing(_ JT808Msg) error {
	return nil
}

func (m *Msg0003) Encode() (pkt []byte, err error) {
	pkt, err = writeHeader(m, pkt)
	return pkt, err
}
