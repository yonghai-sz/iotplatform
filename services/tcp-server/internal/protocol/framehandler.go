package protocol

import (
	"bufio"
	"context"
	"io"
	"net"

	"github.com/pkg/errors"
)

const MaxFrameLen = 2 + 21 + 1023*2 + 1 + 2

var ErrFrameReadEmpty = errors.New("Read empty frame")

type FramePayload []byte

type JT808FrameHandler struct {
	rbuf   *bufio.Reader
	writer io.Writer
}

func NewJT808FrameHandler(conn net.Conn) *JT808FrameHandler {
	return &JT808FrameHandler{
		rbuf:   bufio.NewReader(conn),
		writer: conn,
	}
}

func (fh *JT808FrameHandler) Recv(ctx context.Context) (FramePayload, error) {
	var inMessage bool
	buf := make([]byte, 0, MaxFrameLen)
	readBuf := make([]byte, 1)
	for {
		_, err := fh.rbuf.Read(readBuf)
		if err != nil {
			return nil, errors.Wrap(err, "Fail to read stream to framePayload")
		}
		if readBuf[0] == 0x7e {
			buf = append(buf, readBuf[0])
			if inMessage {
				break
			}
			inMessage = !inMessage
		} else {
			if inMessage {
				buf = append(buf, readBuf[0])
			}
		}
	}

	if len(buf) == 0 {
		return nil, ErrFrameReadEmpty
	}
	return FramePayload(buf), nil
}

func (fh *JT808FrameHandler) Send(payload FramePayload) error {

	var p = payload
	if len(p) == 0 {
		return nil
	}

	for {
		n, err := fh.writer.Write([]byte(p))
		if err != nil {
			return errors.Wrap(err, "Failed to send payload")
		}
		if n >= len(p) {
			break
		}
		if n < len(p) {
			p = p[n:]
		}
	}

	return nil
}
