package wrappers

import (
	"net/http"

	"github.com/Calmantara/lis-backend/internal/helpers/errors"
	"github.com/labstack/echo/v4"
)

const (
	ERROR_INTERNAL_SERVER = iota + 8000
	ERROR_UNPROCESSABLE_ENTITY
	ERROR_BAD_REQUEST
	ERROR_NOT_FOUND
	ERROR_MAIN_COMMAND
	ERROR_CONNECTION
	ERROR_REQUIRED_USERNAME_EMAIL_INPUT
	ERROR_INVALID_PASSWORD
	ERROR_LOGIN_FAILED
	ERROR_CLIENT_ID_NOT_FOUND
	ERROR_CREATE_ACCESS_TOKEN
	ERROR_CREATE_SESSION_TOKEN
	ERROR_UNAUTHORIZED
	ERROR_MISSING_CLIENT_ID
	ERROR_MISSING_USER_ID
	ERROR_MISSING_APPLICATION_ID
	ERROR_MISSING_TOKEN_ID
	ERROR_DUPLICATED_KEY
)

var (
	ERROR_TO_CODE = map[error]int{
		errors.ERROR_INTERNAL_SERVER:               ERROR_INTERNAL_SERVER,
		errors.ERROR_UNPROCESSABLE_ENTITY:          ERROR_UNPROCESSABLE_ENTITY,
		errors.ERROR_BAD_REQUEST:                   ERROR_BAD_REQUEST,
		errors.ERROR_NOT_FOUND:                     ERROR_NOT_FOUND,
		errors.ERROR_MAIN_COMMAND:                  ERROR_MAIN_COMMAND,
		errors.ERROR_CONNECTION:                    ERROR_CONNECTION,
		errors.ERROR_REQUIRED_USERNAME_EMAIL_INPUT: ERROR_REQUIRED_USERNAME_EMAIL_INPUT,
		errors.ERROR_INVALID_PASSWORD:              ERROR_INVALID_PASSWORD,
		errors.ERROR_LOGIN_FAILED:                  ERROR_LOGIN_FAILED,
		errors.ERROR_CLIENT_ID_NOT_FOUND:           ERROR_CLIENT_ID_NOT_FOUND,
		errors.ERROR_CREATE_ACCESS_TOKEN:           ERROR_CREATE_ACCESS_TOKEN,
		errors.ERROR_CREATE_SESSION_TOKEN:          ERROR_CREATE_SESSION_TOKEN,
		errors.ERROR_UNAUTHORIZED:                  ERROR_UNAUTHORIZED,
		errors.ERROR_MISSING_CLIENT_ID:             ERROR_MISSING_CLIENT_ID,
		errors.ERROR_MISSING_USER_ID:               ERROR_MISSING_USER_ID,
		errors.ERROR_MISSING_APPLICATION_ID:        ERROR_MISSING_APPLICATION_ID,
		errors.ERROR_MISSING_TOKEN_ID:              ERROR_MISSING_TOKEN_ID,
		errors.ERROR_DUPLICATED_KEY:                ERROR_DUPLICATED_KEY,
	}
)

type (
	Failure struct {
		Success   bool      `json:"success"`
		Message   string    `json:"message"`
		Errors    []string  `json:"errors"`
		ErrorCode ErrorCode `json:"error_code"`
	}

	ErrorCode struct {
		ErrorID     int    `json:"error_id"`
		Description string `json:"description"`
	}

	ResponseFailureError struct {
		Code    int
		Failure Failure
	}

	errFn func(errorCode ErrorCode, err error) ResponseFailureError
)

const (
	PAYLOAD_REQUESTED_UNSUCCESSFULLY = "payload requested unsuccessfully"
)

func ErrorInternalServer(errorCode ErrorCode, err error) ResponseFailureError {
	errs := errors.Unpack(err)
	return ResponseFailureError{
		Code: http.StatusInternalServerError,
		Failure: Failure{
			Success:   false,
			Message:   PAYLOAD_REQUESTED_UNSUCCESSFULLY,
			Errors:    errs,
			ErrorCode: errorCode,
		},
	}
}

func ErrorBadRequest(errorCode ErrorCode, err error) ResponseFailureError {
	errs := errors.Unpack(err)
	return ResponseFailureError{
		Code: http.StatusBadRequest,
		Failure: Failure{
			Success:   false,
			Message:   PAYLOAD_REQUESTED_UNSUCCESSFULLY,
			Errors:    errs,
			ErrorCode: errorCode,
		},
	}
}

func ErrorNotFound(errorCode ErrorCode, err error) ResponseFailureError {

	errs := errors.Unpack(err)
	return ResponseFailureError{
		Code: http.StatusNotFound,
		Failure: Failure{
			Success:   false,
			Message:   PAYLOAD_REQUESTED_UNSUCCESSFULLY,
			Errors:    errs,
			ErrorCode: errorCode,
		},
	}
}

func ErrorUnprocessable(errorCode ErrorCode, err error) ResponseFailureError {
	errs := errors.Unpack(err)
	return ResponseFailureError{
		Code: http.StatusUnprocessableEntity,
		Failure: Failure{
			Success:   false,
			Message:   PAYLOAD_REQUESTED_UNSUCCESSFULLY,
			Errors:    errs,
			ErrorCode: errorCode,
		},
	}
}

func ErrorUnauthorized(errorCode ErrorCode, err error) ResponseFailureError {
	errs := errors.Unpack(err)
	return ResponseFailureError{
		Code: http.StatusUnauthorized,
		Failure: Failure{
			Success:   false,
			Message:   PAYLOAD_REQUESTED_UNSUCCESSFULLY,
			Errors:    errs,
			ErrorCode: errorCode,
		},
	}
}

func (e *ResponseFailureError) Error() string {
	return e.Failure.ErrorCode.Description
}

func ConstructResponseFailure(ctx echo.Context, e error) error {
	// check error group
	for eid, errGroup := range [][]error{
		errors.BAD_REQUEST,
		errors.NOT_FOUND,
		errors.UNAUTHORIZED,
		errors.UNPROCESSABLE_ENTITY,
		errors.INTERNAL_SERVER,
	} {
		err := matchErrorToCode(ctx, errGroup, eid, e)
		if err != nil {
			return err
		}
	}

	// default to internal server error
	errorCode := ErrorCode{
		ErrorID:     ERROR_INTERNAL_SERVER,
		Description: e.Error(),
	}
	res := ErrorInternalServer(errorCode, e)
	ctx.JSON(res.Code, res.Failure)

	return e
}

func matchErrorToCode(ctx echo.Context, errorGroup []error, eid int, e error) error {
	errKeyVal := map[int]errFn{
		0: ErrorBadRequest,
		1: ErrorNotFound,
		2: ErrorUnauthorized,
		3: ErrorUnprocessable,
		4: ErrorInternalServer,
	}

	for _, err := range errorGroup {
		if errors.Is(e, err) {
			errorCode := ErrorCode{
				ErrorID:     ERROR_TO_CODE[err],
				Description: err.Error(),
			}
			res := errKeyVal[eid](errorCode, e)
			ctx.JSON(res.Code, res.Failure)

			return e
		}
	}

	return nil
}
