package models

import (
	"github.com/Calmantara/lis-backend/internal/core/configurations"
	"github.com/Calmantara/lis-backend/internal/helpers/errors"
	"github.com/google/uuid"
)

type Header struct {
	RequestID     string `header:"X-Request-ID"`
	ApplicationID string `header:"X-Application-ID"`
	ClientID      string `header:"X-Client-ID"`
}

func (h *Header) Validate() error {
	if h.ApplicationID == "" {
		return errors.ERROR_MISSING_APPLICATION_ID
	}

	if h.ClientID == "" {
		return errors.ERROR_MISSING_CLIENT_ID
	}

	if h.RequestID == "" {
		h.RequestID = uuid.NewString()
	}

	return nil
}

type ApplicationKeyHeader struct {
	Header
	ApplicationKey string `header:"X-Application-Key"`
}

func (a *ApplicationKeyHeader) Validate() error {
	if err := a.Header.Validate(); err != nil {
		return err
	}

	if a.ApplicationKey != configurations.Config.Application.Key {
		return errors.ERROR_UNAUTHORIZED
	}

	return nil
}
