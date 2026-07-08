package model

import (
	"net"
	"time"
)

type DeviceStatus int8

const (
	DeviceStatusOffline  DeviceStatus = 0
	DeviceStatusOnline   DeviceStatus = 1
	DeviceStatusSleeping DeviceStatus = 2
)

type Device struct {
	ID    string `json:"id"`
	Plate string `json:"plate"`
	Phone string `json:"phone"`

	SessionID      string        `json:"sessionId"`
	Conn           net.Conn      `json:"-"`
	Keepalive      time.Duration `json:"keepalive"`
	LastestComTime time.Time     `json:"lastComTime"`
	Status         DeviceStatus  `json:"status"`

	ProtocolVersion uint8  `json:"protocolVersion"`
	AuthCode        string `json:"authcode"`
	IMEI            string `json:"imei"`
	SoftwareVersion string `json:"softwareVersion"`
}

func NewDevice(in *Msg0100, session *Session) *Device {
	return &Device{
		ID:    in.DeviceID,
		Plate: in.PlateNumber,
		Phone: in.Header.PhoneNumber,

		SessionID:      session.ID,
		Conn:           session.Conn,
		Keepalive:      time.Minute * 1,
		LastestComTime: time.Now(),
		Status:         DeviceStatusOffline,

		ProtocolVersion: in.Header.ProtocolVersion,
	}
}

func (d *Device) ShouleTurnOffline() bool {
	now := time.Now().UnixMilli()
	return d.Status != DeviceStatusOffline &&
		now > d.LastestComTime.UnixMilli()+d.Keepalive.Milliseconds()
}

func (d *Device) ShouldClear() bool {
	now := time.Now().UnixMilli()
	return d.Status == DeviceStatusOffline &&
		now > d.LastestComTime.UnixMilli()+d.Keepalive.Milliseconds()
}
