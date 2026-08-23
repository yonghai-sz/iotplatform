package model

import (
	"iot-zero/pkg/utils"
)

type Msg0001 struct {
	Header *MsgHeader `json:"header"`

	AnswerSerialNumber uint16 `json:"answerSerialNumber"`
	AnswerMessageID    uint16 `json:"answerMessageId"`
	Result             uint8  `json:"result"`
}

func (m *Msg0001) Decode(packet *PacketData) error {
	m.Header = packet.Header
	pkt, idx := packet.Body, 0

	m.AnswerSerialNumber = utils.ReadWordBigEndian(pkt, &idx)
	m.AnswerMessageID = utils.ReadWordBigEndian(pkt, &idx)
	m.Result = utils.ReadByte(pkt, &idx)
	return nil
}

func (m *Msg0001) GetHeader() *MsgHeader {
	return m.Header
}

func (m *Msg0001) GenOutgoing(incoming JT808Msg) error {
	header := incoming.GetHeader()
	m.AnswerSerialNumber = header.SerialNumber
	m.AnswerMessageID = header.MsgID
	m.Result = 0

	m.Header = header
	m.Header.MsgID = 0x0001
	return nil
}

func (m *Msg0001) Encode() (pkt []byte, err error) {
	pkt = utils.WriteWordBigEndian(pkt, m.AnswerSerialNumber)
	pkt = utils.WriteWordBigEndian(pkt, m.AnswerMessageID)
	pkt = utils.WriteByte(pkt, m.Result)

	pkt, err = writeHeader(m, pkt)
	return pkt, err
}
