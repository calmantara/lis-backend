package mysql

import (
	"context"

	"github.com/Calmantara/lis-backend/internal/core/models"
	"github.com/Calmantara/lis-backend/internal/core/ports"
	"github.com/Calmantara/lis-backend/internal/helpers/errors"
)

type deviceMessageCommandImpl struct {
	mysql MySql
}

func NewDeviceMessageCommand(mysql MySql) ports.DeviceMessageCommand {
	return &deviceMessageCommandImpl{mysql: mysql}
}

func (u *deviceMessageCommandImpl) Create(ctx context.Context, deviceMessage *models.DeviceMessage) error {
	ctx, txn := begin(ctx, u.mysql.Master())

	err := txn.
		WithContext(ctx).
		Create(deviceMessage).Error
	if err != nil {

		return errors.Wrap(err, errors.ERROR_INTERNAL_SERVER.Error())
	}

	return nil
}
