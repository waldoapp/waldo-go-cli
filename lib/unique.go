package lib

import (
	"encoding/hex"

	"github.com/google/uuid"
)

func NewUniqueID() (string, error) {
	uuid, err := uuid.NewRandom()

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(uuid[:]), nil
}
