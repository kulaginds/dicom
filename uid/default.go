package uid

import "encoding/binary"

// ParseTransferSyntaxUID parse TransferSyntaxUID and
// return byte order and implicit.
func ParseTransferSyntaxUID(id string) (binary.ByteOrder, bool) {
	switch id {
	case "1.2.840.10008.1.2":
		return binary.LittleEndian, true
	case "1.2.840.10008.1.2.1", "1.2.840.10008.1.2.1.99":
		return binary.LittleEndian, false
	case "1.2.840.10008.1.2.2":
		return binary.BigEndian, false
	}

	// Big Endian byte ordering was previously described but has been retired, See PS3.5 2016b.
	return binary.LittleEndian, false
}
