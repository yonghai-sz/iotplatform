package model

import (
	"math"
	"net"
	"sync/atomic"

	"github.com/pkg/errors"
)

type (
	SessionCtxKey struct{}

	FrameCtxKey        struct{}
	PacketDecodeCtxKey struct{}
	ProcessDataCtxKey  struct{}
	PacketEncodeCtxKey struct{}
)

var (
	ErrDecodeMsg      = errors.New("Fail to decode msg")
	ErrEncodeMsg      = errors.New("Fail to encode msg")
	ErrGenOutgoingMsg = errors.New("Fail to generate outgoing msg")
)

/*
 *
 *
 */

type Session struct {
	ID   string
	Conn net.Conn

	serialNumber uint32
}

func (s *Session) GetNextSerialNum() uint16 {
	next := atomic.AddUint32(&s.serialNumber, 1)
	if next <= math.MaxUint16 {
		return uint16(next)
	}
	atomic.StoreUint32(&s.serialNumber, 0)
	return uint16(s.serialNumber)
}

/*
 *
 *
 */

type JT808Msg interface {
	Decode(*PacketData) error

	GetHeader() *MsgHeader
	GenOutgoing(incoming JT808Msg) error

	Encode() (pkt []byte, err error)
}

func writeHeader(m JT808Msg, pkt []byte) ([]byte, error) {
	m.GetHeader().Attr.BodyLength = uint16(len(pkt))
	headerPkt, err := m.GetHeader().Encode()
	if err != nil {
		return nil, err
	}
	return append(headerPkt, pkt...), nil
}

/*
 *
 *
 */

type PacketData struct {
	Header *MsgHeader
	Body   []byte

	SegCompleted bool
}

type ProcessData struct {
	Incoming JT808Msg
	Outgoing JT808Msg
}
