package libs

import "github.com/google/uuid"

func ParseUUID(id string) (*uuid.UUID, error) {

	parsedId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return &parsedId, nil
}
