package model

import (
	"time"
)

type DeviceConfig struct {
	IsEnabled     *bool   `json:"isEnabled,omitempty"`
	IsInteractive *bool   `json:"isInteractive,omitempty"`
	Connection    *string `json:"connection,omitempty"`
	SendFrequency *string `json:"sendFrequency,omitempty"`
	Version       *string `json:"version,omitempty"`
}

type Device struct {
	ID        int          `json:"-"`
	PublicID  string       `json:"id" db:"public_id"`
	Status    string       `json:"status"`
	Class     string       `json:"class"`
	Name      string       `json:"name"`
	Config    DeviceConfig `json:"config"`
	CreatedAt time.Time    `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time    `json:"updatedAt" db:"updated_at"`
}

type CreateDeviceDTO struct {
	Class  string       `json:"class" validate:"required"`
	Name   string       `json:"name" validate:"required"`
	Config DeviceConfig `json:"config"`
}
