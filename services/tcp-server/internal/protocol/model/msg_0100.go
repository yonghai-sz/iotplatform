package model

import (
	"strings"

	"iot-zero/pkg/utils"
)

type Msg0100 struct {
	Header *MsgHeader `json:"header"`

	ProvinceID     uint16 `json:"provinceId"`
	CityID         uint16 `json:"cityId"`
	ManufacturerID string `json:"manufacturerId"`
	DeviceMode     string `json:"deviceMode"`
	DeviceID       string `json:"deviceId"`
	PlateColor     byte   `json:"plateColor"`
	PlateNumber    string `json:"plateNumber"`

	LocationDesc string `json:"locationDesc"`
}

func (m *Msg0100) Decode(packet *PacketData) error {
	m.Header = packet.Header
	pkt, idx := packet.Body, 0
	m.ProvinceID = utils.ReadWordBigEndian(pkt, &idx)
	m.CityID = utils.ReadWordBigEndian(pkt, &idx)

	var manuLen, modeLen, idLen int

	manuLen, modeLen, idLen = 11, 30, 30

	cutset := "\x00"
	m.ManufacturerID = strings.TrimRight(utils.ReadString(pkt, &idx, manuLen), cutset)
	m.DeviceMode = strings.TrimRight(utils.ReadString(pkt, &idx, modeLen), cutset)
	m.DeviceID = strings.TrimRight(utils.ReadString(pkt, &idx, idLen), cutset)

	m.PlateColor = utils.ReadByte(pkt, &idx)
	m.PlateNumber = "" // utils.ReadGBK(pkt, &idx, int(m.Header.Attr.BodyLength)-idx)

	// m.LocationDesc = region.Parse(fmt.Sprintf("%02d%04d", m.ProvinceID, m.CityID)).Name

	return nil
}

func (m *Msg0100) GetHeader() *MsgHeader {
	return m.Header
}

func (m *Msg0100) GenOutgoing(_ JT808Msg) error {
	return nil
}

func (m *Msg0100) Encode() (pkt []byte, err error) {
	pkt = utils.WriteWordBigEndian(pkt, m.ProvinceID)
	pkt = utils.WriteWordBigEndian(pkt, m.CityID)

	var manuLen, modeLen, idLen int

	manuLen, modeLen, idLen = 11, 30, 30

	var fillByte byte // '\x00'
	manu := []byte(m.ManufacturerID)
	toFillLen := manuLen - len(manu)
	if toFillLen < 0 {
		manu = manu[:manuLen]
	} else {
		for i := 0; i < toFillLen; i++ {
			manu = append(manu, fillByte)
		}
	}
	pkt = append(pkt, manu...)

	mode := []byte(m.DeviceMode)
	toFillLen = modeLen - len(mode)
	if toFillLen < 0 {
		mode = mode[:modeLen]
	} else {
		for i := 0; i < toFillLen; i++ {
			mode = append(mode, fillByte)
		}
	}
	pkt = append(pkt, mode...)

	id := []byte(m.DeviceID)
	toFillLen = idLen - len(id)
	if toFillLen < 0 {
		id = id[:idLen]
	} else {
		for i := 0; i < toFillLen; i++ {
			id = append(id, fillByte)
		}
	}
	pkt = append(pkt, id...)

	pkt = utils.WriteByte(pkt, m.PlateColor)
	// pkt = utils.WriteGBK(pkt, m.PlateNumber)

	pkt, err = writeHeader(m, pkt)
	return pkt, err
}
