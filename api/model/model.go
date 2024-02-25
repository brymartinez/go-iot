package model

import (
	"time"
)

type DeviceConfig struct {
	IsEnabled     bool   `json:"isEnabled,omitempty"`
	IsInteractive bool   `json:"isInteractive,omitempty"`
	Connection    string `json:"connection,omitempty"`
	SendFrequency string `json:"sendFrequency,omitempty"`
	Version       string `json:"version,omitempty"`
}

type Device struct {
	ID        string       `json:"id"`
	Status    string       `json:"status"`
	Class     string       `json:"class"`
	Config    DeviceConfig `json:"config"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
}

type CreateDeviceDTO struct {
	Class  string       `json:"class" validate:"required"`
	Config DeviceConfig `json:"config"`
}
