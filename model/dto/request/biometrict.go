package request

import "backend-mobile-api/model/enum"

type BiometricRequest struct {
	ServiceType enum.BiometricType `json:"service_type"`
	Token       string             `json:"token"`
	UUID        string             `json:"uuid"`
	DeviceID    string             `json:"device_id"`
}
