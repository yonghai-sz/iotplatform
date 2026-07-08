package model

import (
	"github.com/pkg/errors"

	"iotplatform/pkg/utils"
)

var (
	ErrDecodeHeader = errors.New("Fail to decode header")
	ErrEncodeHeader = errors.New("Fail to encode header")
)

type MsgHeader struct {
	MsgID           uint16            `json:"msgID"`
	Attr            *MsgBodyAttr      `json:"attr"`
	ProtocolVersion uint8             `json:"protocolVersion"`
	PhoneNumber     string            `json:"phoneNumber"`
	SerialNumber    uint16            `json:"serialNumber"`
	Frag            *MsgFragmentation `json:"frag"`

	Idx int `json:"-"`
}

func (h *MsgHeader) Decode(pkt []byte) error {
	var idx int
	h.MsgID = utils.ReadWordBigEndian(pkt, &idx)

	attr := &MsgBodyAttr{}
	err := attr.Decode(utils.ReadWordBigEndian(pkt, &idx))
	if err != nil {
		return ErrDecodeHeader
	}
	h.Attr = attr

	h.ProtocolVersion = utils.ReadByte(pkt, &idx)
	// h.PhoneNumber = utils.ReadBCD(pkt, &idx, 10)
	h.SerialNumber = utils.ReadWordBigEndian(pkt, &idx)

	if h.Attr.PacketFragmented == 1 {
		frag := &MsgFragmentation{}
		err = frag.Decode(pkt, &idx)
		if err != nil {
			return ErrDecodeHeader
		}
		h.Frag = frag
	}

	h.Idx = idx
	return nil
}

func (h *MsgHeader) Encode() (pkt []byte, err error) {
	pkt = utils.WriteWordBigEndian(pkt, h.MsgID)
	pkt = utils.WriteWordBigEndian(pkt, h.Attr.Encode())

	pkt = utils.WriteByte(pkt, h.ProtocolVersion)

	// pkt = utils.WriteBCD(pkt, h.PhoneNumber)
	pkt = utils.WriteWordBigEndian(pkt, h.SerialNumber)
	if h.Frag != nil {
		pkt = append(pkt, h.Frag.Encode()...)
	}

	return pkt, nil
}

func (h *MsgHeader) IsFragmented() bool {
	return h.Attr.PacketFragmented == 1
}

/*
 *
 *
 *
 */

const (
	bodyLengthBit    uint16 = 0b0000001111111111
	encryptionBit    uint16 = 0b0001110000000000
	fragmentationBit uint16 = 0b0010000000000000
	versionSignBit   uint16 = 0b0100000000000000
	extraBit         uint16 = 0b1000000000000000
)

type EncryptionType int8

const (
	EncryptionUnknown EncryptionType = -1
	EncryptionNone    EncryptionType = 0b000
	EncryptionRSA     EncryptionType = 0b001
)

type MsgBodyAttr struct {
	BodyLength       uint16 `json:"bodyLength"`
	Encryption       uint8  `json:"encryption"`
	PacketFragmented uint8  `json:"packetFragmented"`
	VersionSign      uint8  `json:"versionSign"`
	Extra            uint8  `json:"extra"`
}

func (attr *MsgBodyAttr) Decode(bitNum uint16) error {
	attr.BodyLength = bitNum & bodyLengthBit
	attr.Encryption = uint8((bitNum & encryptionBit) >> 10)
	attr.PacketFragmented = uint8(bitNum & fragmentationBit >> 13)
	attr.VersionSign = uint8(bitNum & versionSignBit >> 14)
	attr.Extra = uint8(bitNum & extraBit >> 15)
	return nil
}

func (attr *MsgBodyAttr) Encode() uint16 {
	var bitNum uint16
	bitNum += attr.BodyLength
	bitNum += uint16(attr.Encryption) << 10
	bitNum += uint16(attr.PacketFragmented) << 13
	bitNum += uint16(attr.VersionSign) << 14
	bitNum += uint16(attr.Extra) << 15
	return bitNum
}

/*
 *
 *
 *
 *
 */

type MsgFragmentation struct {
	Total uint16 `json:"total"`
	Index uint16 `json:"index"`
}

func (frag *MsgFragmentation) Decode(pkt []byte, idx *int) error {
	frag.Total = utils.ReadWordBigEndian(pkt, idx)
	frag.Index = utils.ReadWordBigEndian(pkt, idx)
	return nil
}

func (frag *MsgFragmentation) Encode() []byte {
	pkt := make([]byte, 0)
	pkt = utils.WriteWordBigEndian(pkt, frag.Total)
	pkt = utils.WriteWordBigEndian(pkt, frag.Index)
	return pkt
}

func GenMsgHeader(d *Device, msgID, serialNumber uint16) *MsgHeader {
	return &MsgHeader{
		MsgID: msgID,
		Attr: &MsgBodyAttr{
			Encryption:       uint8(EncryptionNone),
			PacketFragmented: 0,
			VersionSign:      1,
			Extra:            0,
		},
		ProtocolVersion: d.ProtocolVersion,
		PhoneNumber:     d.Phone,
		SerialNumber:    serialNumber,
	}
}
