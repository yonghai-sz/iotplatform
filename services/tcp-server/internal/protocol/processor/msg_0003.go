package processor

import (
	"context"

	"github.com/pkg/errors"

	"iotplatform/services/tcp-server/internal/protocol/model"
	"iotplatform/services/tcp-server/internal/storage"
)

func processMsg0003(_ context.Context, data *model.ProcessData) error {

	cache := storage.GetDeviceCache()
	device, err := cache.GetDeviceByPhone(data.Incoming.GetHeader().PhoneNumber)
	if errors.Is(err, storage.ErrDeviceNotFound) {
		return errors.Wrapf(err, "Fail to find device cache, phoneNumber=%s", data.Incoming.GetHeader().PhoneNumber)
	}

	timer := NewKeepaliveTimer()
	timer.Cancel(device.Phone)

	cache.DelDeviceByPhone(device.Phone)

	return nil
}
