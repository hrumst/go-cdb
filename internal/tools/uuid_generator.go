package tools

import (
	cryptoRand "crypto/rand"
	"encoding/binary"
	"fmt"
	"time"
)

func putUint48(b []byte, v uint64) {
	b[0] = byte(v >> 40)
	b[1] = byte(v >> 32)
	b[2] = byte(v >> 24)
	b[3] = byte(v >> 16)
	b[4] = byte(v >> 8)
	b[5] = byte(v)
}

func GenerateUUIDv7() string {
	var uuid [16]byte

	// Encode the timestamp as 48 bits for
	putUint48(uuid[:6], uint64(time.Now().UnixMilli()))
	// Fill remaining 10 bytes with random data
	_, _ = cryptoRand.Read(uuid[6:])

	// Set version to 7 (bits 4-7 in byte 6)
	uuid[6] = (uuid[6] & 0x0F) | 0x70
	// Set the variant to RFC 4122 (bits 6-7 in byte 8)
	uuid[8] = (uuid[8] & 0x3F) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		binary.BigEndian.Uint32(uuid[0:4]),
		binary.BigEndian.Uint16(uuid[4:6]),
		binary.BigEndian.Uint16(uuid[6:8]),
		binary.BigEndian.Uint16(uuid[8:10]),
		uuid[10:],
	)
}
