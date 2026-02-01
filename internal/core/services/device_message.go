package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Calmantara/lis-backend/internal/core/configurations"
	"github.com/Calmantara/lis-backend/internal/core/models"
	"github.com/Calmantara/lis-backend/internal/core/ports"
	"github.com/go-resty/resty/v2"
)

type deviceMessageSvcImpl struct {
	deviceMessageCommand ports.DeviceMessageCommand
	client               *resty.Client
	lisPlatformConfig    configurations.LisPlatform
}

func NewDeviceMessageService(
	deviceMessageCommand ports.DeviceMessageCommand,
	lisPlatform configurations.LisPlatform,
) ports.DeviceMessageService {
	return &deviceMessageSvcImpl{
		lisPlatformConfig:    lisPlatform,
		deviceMessageCommand: deviceMessageCommand,
		client:               setHTTPClient(),
	}
}

func (d *deviceMessageSvcImpl) Process(ctx context.Context, inputs *models.DeviceMessageInput) (err error) {
	// store to database
	deviceMessage := &models.DeviceMessage{
		DeviceID:       inputs.DeviceID,
		DeviceTypeCode: inputs.DeviceTypeCode,
		Message:        inputs.Message,
		Protocol:       inputs.Protocol,
	}
	err = d.deviceMessageCommand.Create(ctx, deviceMessage)
	if err != nil {
		return err
	}

	// routing message
	message, err := d.routingMessage(ctx, *deviceMessage)
	if err != nil || message == nil {
		return
	}

	// publish message
	err = d.publish(ctx, deviceMessage, message)
	if err != nil {
		return nil
	}

	return err
}

func (d *deviceMessageSvcImpl) routingMessage(ctx context.Context, deviceMessage models.DeviceMessage) (parsedMessage models.Message, err error) {

	switch deviceMessage.Protocol {
	case "hl7":
		parsedMessage, err = parseHL7Message(deviceMessage)
		if err != nil {
			return nil, err
		}
	case "rs232":
		switch deviceMessage.DeviceTypeCode {
		default:
			parsedMessage, err = parseUrineTestText(deviceMessage.Message)
			if err != nil {
				return nil, err
			}
		}
	}

	return
}

func (d *deviceMessageSvcImpl) publish(ctx context.Context, deviceMessage *models.DeviceMessage, message models.Message) (err error) {
	serializer := message.Serialize(*deviceMessage)
	cln := d.client.SetTimeout(time.Duration(10) * time.Second)

	// Make a request
	jsonData, _ := json.Marshal(serializer)
	fmt.Println("JSON output:", string(jsonData))

	_, err = cln.R().
		SetHeader("Accept", "application/json").
		SetBody(serializer).
		Post(d.lisPlatformConfig.Url)

	return err
}
