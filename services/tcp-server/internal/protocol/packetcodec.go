package protocol

import (
	"sync"

	"github.com/pkg/errors"

	"iotplatform/services/tcp-server/internal/protocol/model"
	"iotplatform/services/tcp-server/internal/storage"
)

const (
	boundaryMark = 0x7e
	escapeMark   = 0x7d
	escapeOne    = 0x01
	escapeTwo    = 0x02
)

var (
	ErrEmptyPacket  = errors.New("Empty packet")
	ErrVerifyFailed = errors.New("Verify failed")
	ErrEncodeType   = errors.New("Error data type")
)

type JT808PacketCodec struct {
}

var jt808PacketCodec *JT808PacketCodec
var codecOnce sync.Once

func NewJT808PacketCodec() *JT808PacketCodec {
	codecOnce.Do(func() {
		jt808PacketCodec = &JT808PacketCodec{}
	})
	return jt808PacketCodec
}

func (pc *JT808PacketCodec) Decode(payload []byte) (*model.PacketData, error) {

	pkt := unescape(payload)
	pkt, err := verify(pkt)
	if err != nil {
		return nil, err
	}

	header := &model.MsgHeader{}
	err = header.Decode(pkt)
	if err != nil {
		return nil, errors.Wrap(err, "Fail to decode packet")
	}
	pd := &model.PacketData{}
	pd.Header = header
	pd.Body = pkt[header.Idx:]

	if pd.Header.IsFragmented() {
		seg := model.NewSegment(pd)
		pd.SegCompleted = storage.CacheSegment(seg)
		pd.Body = seg.Data
	}

	pd.Header.Idx = 0

	return pd, nil
}

func (pc *JT808PacketCodec) Encode(data any) (pkt []byte, err error) {

	out, ok := data.(model.JT808Msg)
	if !ok {
		return nil, ErrEncodeType
	}
	pkt, err = out.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "Fail to encode jtmsg")
	}

	pkt = genVerifier(pkt)

	payload := escape(pkt)

	return payload, nil
}

func unescape(src []byte) []byte {

	dst := make([]byte, 0)

	i := 1
	n := len(src)

	for i < n-1 {
		if i < n-2 && src[i] == 0x7d && src[i+1] == 0x02 {
			dst = append(dst, boundaryMark)
			i += 2
		} else if i < n-2 && src[i] == 0x7d && src[i+1] == 0x01 {
			dst = append(dst, escapeMark)
			i += 2
		} else {
			dst = append(dst, src[i])
			i++
		}
	}
	return dst
}

func escape(src []byte) []byte {

	dst := make([]byte, 0)
	dst = append(dst, boundaryMark)

	for _, v := range src {
		switch v {
		case boundaryMark:
			dst = append(dst, escapeMark, escapeTwo)
		case escapeMark:
			dst = append(dst, escapeMark, escapeOne)
		default:
			dst = append(dst, v)
		}
	}

	dst = append(dst, boundaryMark)
	return dst
}

func verify(pkt []byte) ([]byte, error) {

	n := len(pkt)
	if n == 0 {
		return nil, ErrEmptyPacket
	}

	expected := pkt[n-1]

	var actual byte
	for _, v := range pkt[:n-1] {
		actual ^= v
	}

	if expected == actual {
		return pkt[:n-1], nil
	}

	return nil, ErrVerifyFailed
}

func genVerifier(pkt []byte) []byte {
	var code byte
	for _, v := range pkt {
		code ^= v
	}
	pkt = append(pkt, code)
	return pkt
}
