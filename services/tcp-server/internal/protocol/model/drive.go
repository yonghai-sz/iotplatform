package model

const SpeedAccuracy = 10

type Drive struct {
	Speed     float64 `json:"speed"`
	Direction uint16  `json:"direction"`
}

func (d *Drive) Decode(m *Msg0200) {
	d.Speed = float64(m.Speed) / SpeedAccuracy
	d.Direction = m.Direction
}
