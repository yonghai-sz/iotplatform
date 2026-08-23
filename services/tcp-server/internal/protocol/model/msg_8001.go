package model

import (
	"iot-zero/pkg/utils"
)

type ResultCode uint8

const (
	ResultSuccess      ResultCode = 0
	ResultFail         ResultCode = 1
	ResultErrMsg       ResultCode = 2
	ResultNotSupported ResultCode = 3
)

type Msg8001 struct {
	Header *MsgHeader `json:"header"`

	AnswerSerialNumber uint16     `json:"answerSerialNumber"`
	AnswerMessageID    uint16     `json:"answerMessageId"`
	Result             ResultCode `json:"result"`
}

func (m *Msg8001) Decode(packet *PacketData) error {
	m.Header = packet.Header
	pkt, idx := packet.Body, 0
	m.AnswerSerialNumber = utils.ReadWordBigEndian(pkt, &idx)
	m.AnswerMessageID = utils.ReadWordBigEndian(pkt, &idx)
	m.Result = ResultCode(utils.ReadByte(pkt, &idx))

	return nil
}

func (m *Msg8001) GetHeader() *MsgHeader {
	return m.Header
}

func (m *Msg8001) GenOutgoing(incoming JT808Msg) error {
	header := incoming.GetHeader()
	m.AnswerSerialNumber = header.SerialNumber
	m.AnswerMessageID = header.MsgID
	m.Result = 0

	m.Header = header
	m.Header.MsgID = 0x8001

	return nil
}

func (m *Msg8001) Encode() (pkt []byte, err error) {
	pkt = utils.WriteWordBigEndian(pkt, m.AnswerSerialNumber)
	pkt = utils.WriteWordBigEndian(pkt, m.AnswerMessageID)
	pkt = utils.WriteByte(pkt, byte(m.Result))

	pkt, err = writeHeader(m, pkt)
	return pkt, err
}
