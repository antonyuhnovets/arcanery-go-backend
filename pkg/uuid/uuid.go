package uuid

import (
	"strconv"

	"github.com/google/uuid"
)

func GenerateId() string {
	id := strconv.Itoa(int(uuid.New().ID()))
	return id
}
