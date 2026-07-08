package model

const LocationAccuracy = 1000000

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  uint16  `json:"altitude"`
}

func (l *Location) Decode(m *Msg0200) {
	l.Latitude = float64(m.Latitude) / LocationAccuracy
	l.Longitude = float64(m.Longitude) / LocationAccuracy
	l.Altitude = m.Altitude
}
