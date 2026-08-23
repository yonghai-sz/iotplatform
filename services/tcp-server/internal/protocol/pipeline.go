package protocol

import (
	"context"
	"net"

	"iot-zero/services/tcp-server/internal/protocol/model"
	"iot-zero/services/tcp-server/internal/protocol/processor"
)

type Pipeline struct {
	fh *JT808FrameHandler
	pc *JT808PacketCodec
	mp *processor.JT808MsgProcessor
}

func NewPipeline(conn net.Conn) *Pipeline {
	return &Pipeline{
		fh: NewJT808FrameHandler(conn),
		pc: NewJT808PacketCodec(),
		mp: processor.NewJT808MsgProcessor(),
	}
}

type delegateFunc func(context.Context, *Pipeline) (context.Context, error)

func (p *Pipeline) ProcessConnRead(ctx context.Context) error {
	actions := []delegateFunc{recv(), decode(), process(), encode(), send()}
	return p.callWithBlocking(ctx, actions)
}

func (p *Pipeline) ProcessConnWrite(ctx context.Context) error {
	actions := []delegateFunc{encode(), send()}
	return p.callWithBlocking(ctx, actions)
}

func (p *Pipeline) callWithBlocking(ctx context.Context, funcs []delegateFunc) error {
	curCtx := ctx
	var err error
	for _, f := range funcs {
		curCtx, err = f(curCtx, p)
		if curCtx == nil || err != nil {
			break
		}
	}
	return err
}

/*
 *
 *
 *
 *
 */

func recv() delegateFunc {
	return delegateFunc(
		func(ctx context.Context, p *Pipeline) (context.Context, error) {
			framePayload, err := p.fh.Recv(ctx)

			wvCtx := context.WithValue(ctx, model.FrameCtxKey{}, framePayload)
			return wvCtx, err
		})
}

func decode() delegateFunc {
	return delegateFunc(
		func(ctx context.Context, p *Pipeline) (context.Context, error) {
			framePayload := ctx.Value(model.FrameCtxKey{}).(FramePayload)

			packet, err := p.pc.Decode(framePayload)

			wvCtx := context.WithValue(ctx, model.PacketDecodeCtxKey{}, packet)
			return wvCtx, err
		})
}

func process() delegateFunc {
	return delegateFunc(
		func(ctx context.Context, p *Pipeline) (context.Context, error) {
			packet := ctx.Value(model.PacketDecodeCtxKey{}).(*model.PacketData)
			if packet == nil {
				return nil, nil
			}

			pd, err := p.mp.Process(ctx, packet)

			wvCtx := context.WithValue(ctx, model.ProcessDataCtxKey{}, pd)
			return wvCtx, err
		})
}

func encode() delegateFunc {
	return delegateFunc(
		func(ctx context.Context, p *Pipeline) (context.Context, error) {

			pd := ctx.Value(model.ProcessDataCtxKey{}).(*model.ProcessData)
			if pd == nil || pd.Outgoing == nil {
				return nil, nil
			}

			pkt, err := p.pc.Encode(pd.Outgoing)
			wvCtx := context.WithValue(ctx, model.PacketEncodeCtxKey{}, pkt)
			return wvCtx, err
		})
}

func send() delegateFunc {
	return delegateFunc(
		func(ctx context.Context, p *Pipeline) (context.Context, error) {
			packet := ctx.Value(model.PacketEncodeCtxKey{}).([]byte)
			err := p.fh.Send(packet)
			return ctx, err
		})
}
