package utils

import "github.com/google/uuid"

func IsValidUUID(u uuid.UUID) bool {
	return u != uuid.Nil
}
