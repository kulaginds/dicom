package parse

import "encoding/binary"

func UL(data []byte, bo binary.ByteOrder) uint32 {
	return bo.Uint32(data)
}
