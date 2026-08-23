package processor

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"iot-zero/services/tcp-server/internal/protocol/model"
	"iot-zero/services/tcp-server/internal/storage"
)

func processMsg0200(_ context.Context, data *model.ProcessData) error {

	in := data.Incoming.(*model.Msg0200)

	cache := storage.GetDeviceCache()
	device, err := cache.GetDeviceByPhone(in.Header.PhoneNumber)
	if errors.Is(err, storage.ErrDeviceNotFound) {
		return errors.Wrapf(err, "Fail to find device cache, phoneNumber=%s", in.Header.PhoneNumber)
	}

	dg := &model.DeviceGeo{}
	err = dg.Decode(device.Phone, in)
	if err != nil {
		return errors.Wrapf(err, "Fail to decode device geo, phoneNumber=%s", device.Phone)
	}

	if dg.Geo.ACCStatus == 0 {
		device.Status = model.DeviceStatusSleeping
		device.LastestComTime = time.Now()
		cache.CacheDevice(device)
	}

	return nil
}
