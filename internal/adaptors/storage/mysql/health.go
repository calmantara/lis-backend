package mysql

import (
	"context"

	"github.com/Calmantara/lis-backend/internal/core/ports"
	"github.com/Calmantara/lis-backend/internal/helpers/errors"
)

type healthRepoImpl struct {
	mysql MySql
}

func NewHealthRepo(mysql MySql) ports.HealthRepo {
	return &healthRepoImpl{mysql: mysql}
}

func (h *healthRepoImpl) CheckMaster(ctx context.Context) (err error) {
	master := h.mysql.Master()
	data := map[string]any{}
	err = master.Raw("SELECT now()").Scan(&data).Error
	if err != nil {
		return errors.Wrap(err, errors.ERROR_CONNECTION.Error())
	}

	return
}

func (h *healthRepoImpl) CheckSlave(ctx context.Context) (err error) {
	slave := h.mysql.Slave()
	data := map[string]any{}
	err = slave.Raw("SELECT now()").Scan(&data).Error
	if err != nil {
		return errors.Wrap(err, errors.ERROR_CONNECTION.Error())
	}

	return
}
