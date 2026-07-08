package model

type DeviceGeo struct {
	Phone string `json:"phone"`

	Location *Location `json:"location"`
	Drive    *Drive    `json:"drive"`

	Geo *GeoMeta `json:"gis"`
}

func (dg *DeviceGeo) Decode(phone string, m *Msg0200) error {
	dg.Phone = phone

	locInstance := &Location{}
	locInstance.Decode(m)
	dg.Location = locInstance

	driveInstance := &Drive{}
	driveInstance.Decode(m)
	dg.Drive = driveInstance

	geoMetaInstance := &GeoMeta{}
	geoMetaInstance.Decode(m.StatusSign)
	dg.Geo = geoMetaInstance

	return nil
}
