package exif

import (
	"encoding/binary"
)

const FIELD_SIZE = 12

type ExifField struct {
	Tag uint16
	FormatType uint16
	Count uint32
	Data []byte
}

const IFD_FORMAT_BYTE = 1
const IFD_FORMAT_STRING = 2
const IFD_FORMAT_USHORT = 3
const IFD_FORMAT_ULONG = 4
const IFD_FORMAT_URATIONAL = 5
const IFD_FORMAT_SBYTE = 6
const IFD_FORMAT_UNDEFINED = 7
const IFD_FORMAT_SSHORT = 8
const IFD_FORMAT_SLONG = 9
const IFD_FORMAT_SRATIONAL = 10
const IFD_FORMAT_SINGLE = 11
const IFD_FORMAT_DOUBLE = 12
const IFD_FORMAT_IFD = 13

type FormatType struct {
	BytePerFormat uint32
	FormatName string
}

var FormatTypeMap = map[uint16]FormatType{
	IFD_FORMAT_BYTE:{1, "BYTE"},
	IFD_FORMAT_STRING:{1, "ASCII"},
	IFD_FORMAT_USHORT:{2, "SHORT"},
	IFD_FORMAT_ULONG:{4, "LONG"},
	IFD_FORMAT_URATIONAL:{8, "RATIONAL"},
	IFD_FORMAT_SBYTE:{1, "SBYTE"},
	IFD_FORMAT_UNDEFINED:{1, "UNDEFINED"},
	IFD_FORMAT_SSHORT:{2, "SSHORT"},
	IFD_FORMAT_SLONG:{4, "SLONG"},
	IFD_FORMAT_SRATIONAL:{8, "SRATIONAL"},
	IFD_FORMAT_SINGLE:{4, "FLOAT"},
	IFD_FORMAT_DOUBLE:{8, "DFLOAT"},
	IFD_FORMAT_IFD:{1, "IFD"},
}

func CreateExifField(payload []byte, offset uint32, byteOrder binary.ByteOrder) *ExifField {
	e := ExifField {
		byteOrder.Uint16(payload[offset:offset+2]),
		byteOrder.Uint16(payload[offset+2:offset+4]),
		byteOrder.Uint32(payload[offset+4:offset+8]),
		payload[offset+8:offset+12],
	}
	byteCount := e.Count * FormatTypeMap[e.FormatType].BytePerFormat
	if byteCount > 4 {
		dataOffset := byteOrder.Uint32(e.Data)
		e.Data = payload[dataOffset:dataOffset+byteCount]
	}
	return &e
}

func (f *ExifField) GetFormatTypeString() string {
	return FormatTypeMap[f.FormatType].FormatName
}