package exif

import (
	"encoding/binary"
)

type ExifDirectory struct {
	fieldList     []*ExifField
	offsetNextIfd uint32
}

func CreateExifDirectory(payload []byte, offset uint32, byteOrder binary.ByteOrder) *ExifDirectory {
	count := int(byteOrder.Uint16(payload[offset:offset+2]))
	var e ExifDirectory
	e.fieldList = make([]*ExifField, count)
	offset += 2
	for i := 0; i < count; i++ {
		e.fieldList[i] = CreateExifField(payload, offset, byteOrder)
		offset += FIELD_SIZE
	}
	e.offsetNextIfd = byteOrder.Uint32(payload[offset:offset+4])
	return &e
}

func (e *ExifDirectory) findField(tag uint16) *ExifField {
	for i := 0; i < len(e.fieldList); i++ {
		if e.fieldList[i].Tag == tag {
			return e.fieldList[i]
		}
	}
	return nil
}

func (e *ExifDirectory) GetFieldList() []*ExifField {
	return e.fieldList
}