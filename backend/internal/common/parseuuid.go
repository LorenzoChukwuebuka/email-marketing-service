package common

import (
	"fmt"
	"github.com/google/uuid"
)

func ParseUUIDMap(inputs map[string]string) (map[string]uuid.UUID, error) {
	result := make(map[string]uuid.UUID)
	for key, val := range inputs {
		parsed, err := uuid.Parse(val)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", key, ErrInvalidUUID)
		}
		result[key] = parsed
	}
	return result, nil
}
