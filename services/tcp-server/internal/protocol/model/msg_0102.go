package model

import (
	"strings"

	"iot-zero/pkg/utils"
)

type Msg0102 struct {
	Header *MsgHeader `json:"header"`

	AuthCodeLen     uint8  `json:"authCodeLen"`
	AuthCode        string `json:"authCode"`
	IMEI            string `json:"imei"`
	SoftwareVersion string `json:"softwareVersion"`
}

func (m *Msg0102) Decode(packet *PacketData) error {
	m.Header = packet.Header
	pkt, idx := packet.Body, 0

	m.AuthCodeLen = utils.ReadByte(pkt, &idx)
	m.AuthCode = utils.ReadString(pkt, &idx, int(m.AuthCodeLen))

	m.IMEI = utils.ReadString(pkt, &idx, 15)
	m.SoftwareVersion = strings.TrimRight(utils.ReadString(pkt, &idx, 20), "\x00")

	return nil
}

func (m *Msg0102) GetHeader() *MsgHeader {
	return m.Header
}

func (m *Msg0102) GenOutgoing(incoming JT808Msg) error {

	in, ok := incoming.(*Msg8100)
	if !ok {
		return ErrGenOutgoingMsg
	}

	m.Header = in.Header
	m.Header.MsgID = 0x0102

	m.AuthCode = in.AuthCode
	m.AuthCodeLen = uint8(len(m.AuthCode))

	return nil
}

func (m *Msg0102) Encode() (pkt []byte, err error) {

	if m.AuthCodeLen != uint8(len(m.AuthCode)) {
		m.AuthCodeLen = uint8(len(m.AuthCode))
	}
	pkt = utils.WriteByte(pkt, m.AuthCodeLen)
	pkt = utils.WriteString(pkt, m.AuthCode)

	var fillByte byte // '\x00'

	im := []byte(m.IMEI)
	toFillLen := 15 - len(m.IMEI)
	for i := 0; i < toFillLen; i++ {
		im = append(im, fillByte)
	}
	pkt = append(pkt, im...)

	sv := []byte(m.SoftwareVersion)
	toFillLen = 20 - len(m.SoftwareVersion)
	for i := 0; i < toFillLen; i++ {
		sv = append(sv, fillByte)
	}
	pkt = append(pkt, sv...)

	pkt, err = writeHeader(m, pkt)
	return pkt, err
}
