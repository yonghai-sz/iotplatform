package model

import (
	"iotplatform/pkg/utils"
)

type ResultCodeType byte

const (
	ResSuccess               ResultCodeType = 0
	ResCarAlreadyRegister    ResultCodeType = 1
	ResCarNotExist           ResultCodeType = 2
	ResDeviceAlreadyRegister ResultCodeType = 3
	ResDeviceNotExist        ResultCodeType = 4
)

type Msg8100 struct {
	Header *MsgHeader `json:"header"`

	AnswerSerialNumber uint16         `json:"answerSerialNumber"`
	Result             ResultCodeType `json:"result"`
	AuthCode           string         `json:"authCode"`
}

func (m *Msg8100) Decode(packet *PacketData) error {
	m.Header = packet.Header
	pkt, idx := packet.Body, 0
	m.AnswerSerialNumber = utils.ReadWordBigEndian(pkt, &idx)
	m.Result = ResultCodeType(utils.ReadByte(pkt, &idx))
	m.AuthCode = utils.ReadString(pkt, &idx, int(m.Header.Attr.BodyLength)-idx)
	return nil
}

func (m *Msg8100) GetHeader() *MsgHeader {
	return m.Header
}

func (m *Msg8100) GenOutgoing(incoming JT808Msg) error {

	in, ok := incoming.(*Msg0100)
	if !ok {
		return ErrGenOutgoingMsg
	}

	m.AnswerSerialNumber = in.Header.SerialNumber
	m.Result = 0
	m.AuthCode = "AuthCode"

	m.Header = in.Header
	m.Header.MsgID = 0x8100
	return nil
}

func (m *Msg8100) Encode() (pkt []byte, err error) {
	pkt = utils.WriteWordBigEndian(pkt, m.AnswerSerialNumber)
	pkt = utils.WriteByte(pkt, uint8(m.Result))
	pkt = utils.WriteString(pkt, m.AuthCode)

	pkt, err = writeHeader(m, pkt)
	return pkt, err
}
