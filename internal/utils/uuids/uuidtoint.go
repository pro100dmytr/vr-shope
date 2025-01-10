package uuids

import (
	"github.com/google/uuid"
	"hash/crc64"
)

func UUIDToInt(u uuid.UUID) uint64 {
	table := crc64.MakeTable(crc64.ECMA)
	return crc64.Checksum(u[:], table)
}
