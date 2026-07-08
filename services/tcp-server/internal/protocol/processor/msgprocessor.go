package processor

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"iotplatform/services/tcp-server/internal/protocol/model"
	"iotplatform/services/tcp-server/internal/storage"
)

var (
	ErrMsgIDNotSupportted = errors.New("Msg id is not supportted")
	ErrNotAuthorized      = errors.New("Not authorized")
)

type action struct {
	genData func() *model.ProcessData
	process func(context.Context, *model.ProcessData) error
}

type processOptions map[uint16]*action

func initProcessOption() processOptions {

	options := make(processOptions)

	options[0x0001] = &action{
		genData: func() *model.ProcessData { return &model.ProcessData{Incoming: &model.Msg0001{}} },
	}

	options[0x0002] = &action{
		genData: func() *model.ProcessData {
			return &model.ProcessData{Incoming: &model.Msg0002{}, Outgoing: &model.Msg8001{}}
		},
		process: processMsg0002,
	}
	options[0x0003] = &action{
		genData: func() *model.ProcessData {
			return &model.ProcessData{Incoming: &model.Msg0003{}, Outgoing: &model.Msg8001{}}
		},
		process: processMsg0003,
	}

	options[0x0100] = &action{
		genData: func() *model.ProcessData {
			return &model.ProcessData{Incoming: &model.Msg0100{}, Outgoing: &model.Msg8100{}}
		},
		process: processMsg0100,
	}
	options[0x0102] = &action{
		genData: func() *model.ProcessData {
			return &model.ProcessData{Incoming: &model.Msg0102{}, Outgoing: &model.Msg8001{}}
		},
		process: processMsg0102,
	}

	options[0x0200] = &action{
		genData: func() *model.ProcessData {
			return &model.ProcessData{Incoming: &model.Msg0200{}, Outgoing: &model.Msg8001{}}
		},
		process: processMsg0200,
	}

	return options
}

/*
 *
 *
 *
 */

type JT808MsgProcessor struct {
	options processOptions
}

var jt808MsgProcessorSingleton *JT808MsgProcessor
var processorInitOnce sync.Once

func NewJT808MsgProcessor() *JT808MsgProcessor {
	processorInitOnce.Do(func() {
		jt808MsgProcessorSingleton = &JT808MsgProcessor{
			options: initProcessOption(),
		}
	})
	return jt808MsgProcessorSingleton
}

/*
 *
 *
 *
 */

func (mp *JT808MsgProcessor) Process(ctx context.Context, pkt *model.PacketData) (*model.ProcessData, error) {

	msgID := pkt.Header.MsgID
	act, ok := mp.options[msgID]
	if !ok {
		return nil, ErrMsgIDNotSupportted
	}

	if pkt.Header.IsFragmented() && !pkt.SegCompleted {
		return processSegmentPacket(ctx, pkt)
	}

	genDataFn := act.genData
	if genDataFn == nil {
		return nil, ErrMsgIDNotSupportted
	}
	data := genDataFn()

	in := data.Incoming
	err := in.Decode(pkt)
	if err != nil {
		return nil, errors.Wrap(err, "Fail to decode packet to jtmsg")
	}

	if data.Outgoing != nil {
		out := data.Outgoing
		err = out.GenOutgoing(in)
		if err != nil {
			return data, errors.Wrap(err, "Fail to generate outgoing msg")
		}
	}

	processFunc := act.process
	if processFunc == nil {
		return data, nil
	}
	err = processFunc(ctx, data)
	if err != nil {
		return data, errors.Wrap(err, "Fail to process data")
	}

	if data.Outgoing == nil {
		return nil, nil
	}
	return data, nil
}

func processSegmentPacket(_ context.Context, pkt *model.PacketData) (*model.ProcessData, error) {

	phone := pkt.Header.PhoneNumber

	cache := storage.GetDeviceCache()
	device, err := cache.GetDeviceByPhone(phone)
	if errors.Is(err, storage.ErrDeviceNotFound) {
		return nil, errors.Wrapf(err, "Fail to find device cache, phoneNumber=%s", phone)
	}

	session, err := storage.GetSession(device.SessionID)

	header := model.GenMsgHeader(device, 0x8001, session.GetNextSerialNum())

	outgoingMsg := &model.Msg8001{
		Header: header,

		AnswerSerialNumber: pkt.Header.SerialNumber,
		AnswerMessageID:    pkt.Header.MsgID,
		Result:             model.ResultSuccess,
	}

	return &model.ProcessData{Outgoing: outgoingMsg}, nil
}
