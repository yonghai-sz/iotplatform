package processor

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"iot-zero/services/tcp-server/internal/protocol/model"
	"iot-zero/services/tcp-server/internal/storage"
)

func processMsg0002(_ context.Context, data *model.ProcessData) error {

	cache := storage.GetDeviceCache()
	device, err := cache.GetDeviceByPhone(data.Incoming.GetHeader().PhoneNumber)
	if errors.Is(err, storage.ErrDeviceNotFound) {
		return errors.Wrapf(err, "Fail to find device cache, phoneNumber=%s", data.Incoming.GetHeader().PhoneNumber)
	}

	device.LastestComTime = time.Now()
	cache.CacheDevice(device)

	return nil
}
