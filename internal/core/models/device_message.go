package models

import (
	"encoding/base64"

	"github.com/Calmantara/lis-backend/internal/helpers/errors"
	"github.com/google/uuid"
)

type DeviceMessage struct {
	ID             *uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:uuid_generate_v4();"`
	DeviceID       string     `json:"device_id" gorm:"column:device_id"`
	DeviceTypeCode string     `json:"device_type_code" gorm:"column:device_type_code"`
	Message        string     `json:"message" gorm:"column:message"`
	Protocol       string     `json:"protocol" gorm:"column:protocol"`
	Default
}

func (DeviceMessage) TableName() string {
	return "device_messages"
}

type DeviceMessageParam struct {
	DeviceID       string `json:"device_id" form:"device_id"`
	DeviceTypeCode string `json:"device_type_code" form:"device_type_code"`
	Message        string `json:"message" form:"message"`
	Protocol       string `json:"protocol" form:"protocol"`
}

func (d *DeviceMessageParam) Validate() error {
	if d.DeviceID == "" || d.DeviceTypeCode == "" || d.Message == "" || d.Protocol == "" {
		return errors.ERROR_BAD_REQUEST
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(d.Message)
	if err != nil {
		return err
	}
	d.Message = string(decodedBytes)

	return nil
}

func (d DeviceMessageParam) ToInput() *DeviceMessageInput {
	return &DeviceMessageInput{
		DeviceID:       d.DeviceID,
		DeviceTypeCode: d.DeviceTypeCode,
		Message:        d.Message,
		Protocol:       d.Protocol,
	}
}

type DeviceMessageInput struct {
	DeviceID       string `json:"device_id"`
	DeviceTypeCode string `json:"device_type_code"`
	Message        string `json:"message"`
	Protocol       string `json:"protocol"`
}
