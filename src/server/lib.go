package server

import (
	"encoding/binary"
)

// 4bytes to uint32
func BytesToUInt32(buf []byte) uint32 {
	return uint32(binary.LittleEndian.Uint32(buf))
}

// uint32 to 4bytes
func UInt32ToBytes(num uint32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, num)
	return bytes
}
