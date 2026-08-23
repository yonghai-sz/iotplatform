package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"

	"iot-zero/services/tcp-server/internal/protocol"
	"iot-zero/services/tcp-server/internal/protocol/model"
	"iot-zero/services/tcp-server/internal/storage"
	"runtime/debug"

	"github.com/pkg/errors"
)

type TCPServer struct {
	listener net.Listener
}

func NewTCPServer() *TCPServer {
	return &TCPServer{}
}

func (serv *TCPServer) Listen(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	serv.listener = l
	return nil
}

func (serv *TCPServer) Start() {
	for {
		conn, err := serv.listener.Accept()
		if err != nil {
			continue
		}
		session := &model.Session{
			Conn: conn,
			ID:   conn.RemoteAddr().String(),
		}
		storage.StoreSession(session)
		go serve(session)
	}
}

func (serv *TCPServer) Stop() {
	serv.listener.Close()
}

func remove(session *model.Session) {
	session.Conn.Close()
	storage.ClearSession(session.ID)
}

func serve(session *model.Session) {

	defer func() {
		if p := recover(); p != nil {
			s := string(debug.Stack())
			fmt.Printf("err=%s, stack=%s", fmt.Sprint(p), s)
		}
	}()

	defer remove(session)

	pg := protocol.NewPipeline(session.Conn)
	for {
		ctx := context.WithValue(context.Background(), model.SessionCtxKey{}, session)
		err := pg.ProcessConnRead(ctx)
		if err == nil {
			continue
		}

		switch {
		case errors.Is(err, io.EOF), errors.Is(err, io.ErrClosedPipe), errors.Is(err, net.ErrClosed),
			errors.Is(err, storage.ErrDeviceNotFound):
			return // close connection when EOF or closed
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func Send(id string, msg model.JT808Msg) {

	session, err := storage.GetSession(id)
	if err != nil && errors.Is(err, storage.ErrSessionClosed) {
		return
	}

	pg := protocol.NewPipeline(session.Conn)

	ctx := context.WithValue(context.Background(), model.ProcessDataCtxKey{}, &model.ProcessData{Outgoing: msg})

	err = pg.ProcessConnWrite(ctx)
	if err == nil {
		return
	}

	if errors.Is(err, io.ErrClosedPipe) || errors.Is(err, net.ErrClosed) {
		remove(session)
	}

}
