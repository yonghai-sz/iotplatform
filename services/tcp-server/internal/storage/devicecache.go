package storage

import (
	"errors"
	"sync"

	"golang.org/x/exp/maps"

	"iotplatform/services/tcp-server/internal/protocol/model"
)

var ErrDeviceNotFound = errors.New("device not found")

type DeviceCache struct {
	cacheByPlate map[string]*model.Device
	cacheByPhone map[string]*model.Device
	mu           sync.Mutex
}

var deviceCacheSingleton *DeviceCache
var deviceCacheInitOnce sync.Once

func GetDeviceCache() *DeviceCache {
	deviceCacheInitOnce.Do(func() {
		deviceCacheSingleton = &DeviceCache{
			cacheByPlate: make(map[string]*model.Device),
			cacheByPhone: make(map[string]*model.Device),
		}
	})
	return deviceCacheSingleton
}

func (cache *DeviceCache) ListDevice() []*model.Device {
	return maps.Values(cache.cacheByPhone)
}

func (cache *DeviceCache) GetDeviceByPlate(carPlate string) (*model.Device, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if d, ok := cache.cacheByPlate[carPlate]; ok {
		return d, nil
	}
	return nil, ErrDeviceNotFound
}

func (cache *DeviceCache) HasPlate(carPlate string) bool {
	d, err := cache.GetDeviceByPlate(carPlate)
	return d != nil && err == nil
}

func (cache *DeviceCache) GetDeviceByPhone(phone string) (*model.Device, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if d, ok := cache.cacheByPhone[phone]; ok {
		return d, nil
	}

	return nil, ErrDeviceNotFound
}

func (cache *DeviceCache) HasPhone(phone string) bool {
	d, err := cache.GetDeviceByPhone(phone)
	return d != nil && err == nil
}

func (cache *DeviceCache) cacheDevice(d *model.Device) {
	cache.cacheByPlate[d.Plate] = d
	cache.cacheByPhone[d.Phone] = d
}

func (cache *DeviceCache) CacheDevice(d *model.Device) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.cacheDevice(d)
}

func (cache *DeviceCache) delDevice(carPlate, phone *string) {

	var d *model.Device
	var ok bool

	if carPlate != nil {
		d, ok = cache.cacheByPlate[*carPlate]
	}

	if phone != nil {
		d, ok = cache.cacheByPhone[*phone]
	}

	if !ok {
		return
	}

	delete(cache.cacheByPlate, d.Plate)
	delete(cache.cacheByPhone, d.Phone)
}

func (cache *DeviceCache) DelDeviceByCarPlate(carPlate string) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.delDevice(&carPlate, nil)
}

func (cache *DeviceCache) DelDeviceByPhone(phone string) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.delDevice(nil, &phone)
}
