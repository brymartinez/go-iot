package model

import (
	"time"
)

type DeviceConfig struct {
	IsEnabled     bool   `json:"isEnabled"`
	IsInteractive bool   `json:"isInteractive"`
	Connection    string `json:"connection"`
	SendFrequency string `json:"sendFrequency"`
	Version       string `json:"version"`
}

type Device struct {
	ID        string       `json:"id"`
	Status    string       `json:"status"`
	Class     string       `json:"class"`
	Config    DeviceConfig `json:"config"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
}
