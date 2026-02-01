package wrappers

import (
	"net/http"
	"time"

	"github.com/Calmantara/lis-backend/internal/utils"
	"github.com/labstack/echo/v4"
)

type (
	Success struct {
		Success  bool   `json:"success"`
		Message  string `json:"message"`
		Data     any    `json:"data"`
		Metadata any    `json:"metadata"`
		Topic    string `json:"topic,omitempty"`
	}
	ResponseSuccess struct {
		Code    int
		Success Success
	}

	Metadata struct {
		TimeFrom  *time.Time `json:"time_from,omitempty"`
		TimeTo    *time.Time `json:"time_to,omitempty"`
		SortBy    string     `json:"sort_by,omitempty"`
		Page      int        `json:"page,omitempty"`
		Limit     int        `json:"limit,omitempty"`
		Total     int64      `json:"total,omitempty"`
		TotalData int64      `json:"total_data,omitempty"`
	}
)

const (
	PAYLOAD_CREATED_SUCCESSFULLY   = "payload created successfully"
	PAYLOAD_REQUESTED_SUCCESSFULLY = "payload requested successfully"
	PAYLOAD_FETCHED_SUCCESSFULLY   = "payload fetched successfully"
	PAYLOAD_UPDATED_SUCCESSFULLY   = "payload updated successfully"
	PAYLOAD_DELETED_SUCCESSFULLY   = "payload deleted successfully"
)

func SuccessCreated(data, metadata any) ResponseSuccess {
	return ResponseSuccess{
		Code: http.StatusCreated,
		Success: Success{
			Success:  true,
			Message:  PAYLOAD_CREATED_SUCCESSFULLY,
			Data:     data,
			Metadata: TransformMetadata(metadata),
		},
	}
}

func SuccessOK(data, metadata any) ResponseSuccess {
	return ResponseSuccess{
		Code: http.StatusOK,
		Success: Success{
			Success:  true,
			Message:  PAYLOAD_REQUESTED_SUCCESSFULLY,
			Data:     data,
			Metadata: TransformMetadata(metadata),
		},
	}
}

func SuccessFetched(data, metadata any) ResponseSuccess {
	return ResponseSuccess{
		Code: http.StatusOK,
		Success: Success{
			Success:  true,
			Message:  PAYLOAD_FETCHED_SUCCESSFULLY,
			Data:     data,
			Metadata: TransformMetadata(metadata),
		},
	}
}

func SuccessFetchedWithTopic(data, metadata any, topic string) ResponseSuccess {
	return ResponseSuccess{
		Code: http.StatusOK,
		Success: Success{
			Success:  true,
			Message:  PAYLOAD_FETCHED_SUCCESSFULLY,
			Data:     data,
			Metadata: TransformMetadata(metadata),
			Topic:    topic,
		},
	}
}

func SuccessUpdated(data, metadata any) ResponseSuccess {
	return ResponseSuccess{
		Code: http.StatusOK,
		Success: Success{
			Success:  true,
			Message:  PAYLOAD_UPDATED_SUCCESSFULLY,
			Data:     data,
			Metadata: TransformMetadata(metadata),
		},
	}
}

func SuccessDeleted(data, metadata any) ResponseSuccess {
	return ResponseSuccess{
		Code: http.StatusOK,
		Success: Success{
			Success:  true,
			Message:  PAYLOAD_DELETED_SUCCESSFULLY,
			Data:     data,
			Metadata: TransformMetadata(metadata),
		},
	}
}

func TransformMetadata(metadata any) *Metadata {
	res := &Metadata{}
	utils.ObjectMapper(&metadata, res)

	if !res.IsValidMetadata() {
		return nil
	}

	return res
}

func (m Metadata) IsValidMetadata() bool {
	if m.Page <= 0 && m.Limit <= 0 && m.Total <= 0 && m.TotalData <= 0 {
		return false
	}

	return true
}

func ConstructResponseSuccess(ctx echo.Context, res ResponseSuccess) error {
	ctx.JSON(res.Code, res.Success)

	return nil
}
