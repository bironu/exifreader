package exif

import(
	"reflect"
	"encoding/binary"
)
type ExifHeader struct {
	byteOrder binary.ByteOrder
	ifdOffset uint32
}

func CreateExifHeader(payload []byte) *ExifHeader{
	var e ExifHeader
	e.byteOrder = getByteOrder(payload[0:2])
	e.ifdOffset = e.byteOrder.Uint32(payload[4:8])
	return &e
}

func getByteOrder(buf[]byte) binary.ByteOrder {
	BE := []byte{0x4D, 0x4D}
	LE := []byte{0x49, 0x49}

	if reflect.DeepEqual(buf[:2], BE) {
		return binary.BigEndian
	} else if reflect.DeepEqual(buf[:2], LE) {
		return binary.LittleEndian
	} else {
		panic("endian invalid!")
	}
}
