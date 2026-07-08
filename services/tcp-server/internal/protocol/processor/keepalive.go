package processor

import (
	"sync"

	"github.com/fakeyanss/gron"
	"github.com/pkg/errors"

	"iotplatform/services/tcp-server/internal/protocol/model"
	"iotplatform/services/tcp-server/internal/storage"
)

type KeepaliveTimer struct {
	cron *gron.Cron
}

var timerSingleton *KeepaliveTimer
var timerInitOnce sync.Once

func NewKeepaliveTimer() *KeepaliveTimer {

	timerInitOnce.Do(func() {
		timerSingleton = &KeepaliveTimer{
			cron: gron.New(),
		}
		timerSingleton.cron.Start()
	})

	return timerSingleton
}

/*
 *
 *
 */

func (t *KeepaliveTimer) Register(devicePhone string) {

	cache := storage.GetDeviceCache()
	device, err := cache.GetDeviceByPhone(devicePhone)
	if errors.Is(err, storage.ErrDeviceNotFound) {
		return
	}

	job := &CheckDeviceJob{phone: devicePhone}

	t.cron.Add(gron.Every(device.Keepalive), job)
}

func (t *KeepaliveTimer) Cancel(devicePhone string) {
	t.cron.Cancel(devicePhone)
}

func (t *KeepaliveTimer) Jobs() []*gron.Entry {
	return t.cron.Entries()
}

/*
 *
 *
 *
 */

type CheckDeviceJob struct {
	phone string
}

func (j *CheckDeviceJob) JobID() string {
	return j.phone
}

func (j *CheckDeviceJob) Run() {
	checkDeviceKeepalive(timerSingleton, j.phone)
}

func checkDeviceKeepalive(t *KeepaliveTimer, devicePhone string) {

	cache := storage.GetDeviceCache()
	d, err := cache.GetDeviceByPhone(devicePhone)
	if errors.Is(err, storage.ErrDeviceNotFound) {
		t.Cancel(devicePhone)
	}

	if d.ShouleTurnOffline() {
		d.Status = model.DeviceStatusOffline
		cache.CacheDevice(d)

	} else if d.ShouldClear() {

		d.Conn.Close()

		cache.DelDeviceByPhone(devicePhone)

		t.Cancel(devicePhone)
	}
}
