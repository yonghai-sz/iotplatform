package processor

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"iot-zero/services/tcp-server/internal/protocol/model"
	"iot-zero/services/tcp-server/internal/storage"
)

func processMsg0102(_ context.Context, data *model.ProcessData) error {

	in := data.Incoming.(*model.Msg0102)

	cache := storage.GetDeviceCache()
	device, err := cache.GetDeviceByPhone(in.Header.PhoneNumber)
	if errors.Is(err, storage.ErrDeviceNotFound) {
		return errors.Wrapf(err, "Fail to find device cache, phoneNumber=%s", in.Header.PhoneNumber)
	}

	out := data.Outgoing.(*model.Msg8001)

	if in.AuthCode != genAuthCode(device) {
		out.Result = model.ResultFail

		timer := NewKeepaliveTimer()
		timer.Cancel(device.Phone)

		cache.DelDeviceByPhone(device.Phone)
	} else {

		device.Status = model.DeviceStatusOnline
		device.LastestComTime = time.Now()
		device.AuthCode = in.AuthCode
		device.IMEI = in.IMEI
		device.SoftwareVersion = in.SoftwareVersion
		cache.CacheDevice(device)
	}

	return nil
}
