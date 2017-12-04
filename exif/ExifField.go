package exif

import (
	"fmt"
	"encoding/binary"
)

const FIELD_SIZE = 12

type ExifField interface {
	GetTag() uint16
	GetType() uint16
	GetCount() uint32
	GetData() []byte
	String() string // for Stringer I/F
}

type ExifCommonField struct {
	tag_ uint16
	type_ uint16
	count_ uint32
	data_ []byte
	byteOrder_ binary.ByteOrder
}

func (f *ExifCommonField) GetTag() uint16 {
	return f.tag_
}
func (f *ExifCommonField) GetType() uint16 {
	return f.type_
}
func (f *ExifCommonField) GetCount() uint32 {
	return f.count_
}
func (f *ExifCommonField) GetData() []byte {
	return f.data_
}
func (f *ExifCommonField) String() string {
	return fmt.Sprintf("%04X\t%-9s\t%10d\t", f.GetTag(), f.GetFormatTypeString(), f.GetCount())
}
func (f *ExifCommonField) GetFormatTypeString() string {
	return FormatTypeMap[f.GetType()].FormatName
}

const IFD_FORMAT_UBYTE = 1
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
	IFD_FORMAT_UBYTE:{1, "BYTE"},
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

func CreateExifField(payload []byte, offset uint32, byteOrder binary.ByteOrder) ExifField {
	common := ExifCommonField {
		tag_ : byteOrder.Uint16(payload[offset:offset+2]),
		type_ : byteOrder.Uint16(payload[offset+2:offset+4]),
		count_ : byteOrder.Uint32(payload[offset+4:offset+8]),
		data_ : payload[offset+8:offset+12],
		byteOrder_ : byteOrder,
	}
	byteCount := common.GetCount() * FormatTypeMap[common.GetType()].BytePerFormat
	if byteCount > 4 {
		dataOffset := byteOrder.Uint32(common.GetData())
		common.data_ = payload[dataOffset:dataOffset+byteCount]
	}

	var e ExifField
	switch common.GetType() {
	case IFD_FORMAT_UBYTE: e = &ExifUByteField{&common}
	case IFD_FORMAT_STRING: e = &ExifStringField{&common}
	case IFD_FORMAT_USHORT: e = &ExifUShortField{&common}
	case IFD_FORMAT_ULONG: e = &ExifULongField{&common}
	case IFD_FORMAT_URATIONAL: e = &ExifURationalField{&common}
	case IFD_FORMAT_SBYTE: e = &ExifSByteField{&common}
	case IFD_FORMAT_UNDEFINED: e = &ExifUndefinedField{&common}
	case IFD_FORMAT_SSHORT: e = &ExifSShortField{&common}
	case IFD_FORMAT_SLONG: e = &ExifSLongField{&common}
	case IFD_FORMAT_SRATIONAL: e = &ExifSRationalField{&common}
	case IFD_FORMAT_SINGLE: e = &ExifSingleFloatField{&common}
	case IFD_FORMAT_DOUBLE: e = &ExifDoubleFloatField{&common}
	case IFD_FORMAT_IFD: e = &ExifIFDField{&common}
	default:
	}
	
	return e
}

