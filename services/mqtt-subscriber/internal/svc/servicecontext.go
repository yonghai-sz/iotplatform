package svc

import "iot-zero/services/mqtt-subscriber/internal/config"

// ServiceContext carries dependencies for handlers and logic.
type ServiceContext struct {
	Config config.Config
}

// NewServiceContext builds a service context from config.
func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
