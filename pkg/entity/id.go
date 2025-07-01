package entity

import (
	"fmt"

	"github.com/google/uuid"
)

type ID = uuid.UUID

func NewID() ID {
	return ID(uuid.New())
}

func ParseID(s string) (ID, error) {
	_, err := uuid.Parse(s)
	if err != nil {
		return ID{}, fmt.Errorf("invalid UUID: %w", err)
	}
	return ID(uuid.MustParse(s)), nil
}
