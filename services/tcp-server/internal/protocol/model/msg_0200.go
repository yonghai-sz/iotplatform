package model

import (
	"iot-zero/pkg/utils"
)

type Msg0200 struct {
	Header *MsgHeader `json:"header"`

	AlarmSign  uint32 `json:"alarmSign"`
	StatusSign uint32 `json:"statusSign"`

	Latitude  uint32 `json:"latitude"`
	Longitude uint32 `json:"longitude"`
	Altitude  uint16 `json:"altitude"`

	Speed     uint16 `json:"speed"`
	Direction uint16 `json:"direction"`
	Time      string `json:"time"`
}

func (m *Msg0200) Decode(packet *PacketData) error {
	m.Header = packet.Header
	pkt, idx := packet.Body, 0
	m.AlarmSign = utils.ReadDoubleWordBigEndian(pkt, &idx)
	m.StatusSign = utils.ReadDoubleWordBigEndian(pkt, &idx)
	m.Latitude = utils.ReadDoubleWordBigEndian(pkt, &idx)
	m.Longitude = utils.ReadDoubleWordBigEndian(pkt, &idx)
	m.Altitude = utils.ReadWordBigEndian(pkt, &idx)
	m.Speed = utils.ReadWordBigEndian(pkt, &idx)
	m.Direction = utils.ReadWordBigEndian(pkt, &idx)
	// m.Time = utils.ReadBCD(pkt, &idx, 6)
	return nil
}

func (m *Msg0200) GetHeader() *MsgHeader {
	return m.Header
}

func (m *Msg0200) GenOutgoing(_ JT808Msg) error {
	return nil
}

func (m *Msg0200) Encode() (pkt []byte, err error) {
	pkt = utils.WriteDoubleWordBigEndian(pkt, m.AlarmSign)
	pkt = utils.WriteDoubleWordBigEndian(pkt, m.StatusSign)
	pkt = utils.WriteDoubleWordBigEndian(pkt, m.Latitude)
	pkt = utils.WriteDoubleWordBigEndian(pkt, m.Longitude)
	pkt = utils.WriteWordBigEndian(pkt, m.Altitude)
	pkt = utils.WriteWordBigEndian(pkt, m.Speed)
	pkt = utils.WriteWordBigEndian(pkt, m.Direction)
	// pkt = utils.WriteBCD(pkt, m.Time)

	pkt, err = writeHeader(m, pkt)
	return pkt, err
}
