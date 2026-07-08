package model

type Msg0002 struct {
	Header *MsgHeader `json:"header"`
}

func (m *Msg0002) Decode(packet *PacketData) error {
	m.Header = packet.Header
	return nil
}

func (m *Msg0002) GetHeader() *MsgHeader {
	return m.Header
}

func (m *Msg0002) GenOutgoing(_ JT808Msg) error {
	return nil
}

func (m *Msg0002) Encode() (pkt []byte, err error) {
	pkt, err = writeHeader(m, pkt)
	return pkt, err
}
