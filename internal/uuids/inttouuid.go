package uuids

import (
	"encoding/binary"
	"github.com/google/uuid"
)

func IntToUUID(id int64) uuid.UUID {
	u := uuid.Nil

	u[8] = (u[8] & 0x3f) | 0x80

	binary.BigEndian.PutUint64(u[8:], uint64(id))

	return u
}
