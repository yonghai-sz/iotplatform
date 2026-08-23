package processor

import (
	"context"

	"iot-zero/services/tcp-server/internal/protocol/model"
	"iot-zero/services/tcp-server/internal/storage"
)

func processMsg0100(ctx context.Context, data *model.ProcessData) error {

	in := data.Incoming.(*model.Msg0100)
	out := data.Outgoing.(*model.Msg8100)

	cache := storage.GetDeviceCache()
	if cache.HasPlate(in.PlateNumber) {
		out.Result = model.ResCarAlreadyRegister
		return nil
	}
	if cache.HasPhone(in.Header.PhoneNumber) {
		out.Result = model.ResDeviceAlreadyRegister
		return nil
	}

	session := ctx.Value(model.SessionCtxKey{}).(*model.Session)

	device := model.NewDevice(in, session)

	out.AuthCode = genAuthCode(device)

	cache.CacheDevice(device)

	timer := NewKeepaliveTimer()
	timer.Register(device.Phone)
	return nil
}
