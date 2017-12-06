package exif

import (
	"fmt"
	"encoding/binary"
)

const (
	FIELD_SIZE = 12
	IFD_EXIF_TAG_NUMBER = 0x8769
	IFD_GPS_TAG_NUMBER = 0x8825
	IFD_INTEROPERABILITY_TAG_NUMBER = 0xA005
)

type ExifTagNumber uint16
type ExifFormatId uint16

type ExifField interface {
	GetTag() ExifTagNumber
	GetType() ExifFormatId
	GetCount() uint32
	GetData() []byte
	String() string // for Stringer I/F
}

type ExifCommonField struct {
	tag_ ExifTagNumber
	type_ ExifFormatId
	count_ uint32
	data_ []byte
	byteOrder_ binary.ByteOrder
}

func (f *ExifCommonField) GetTag() ExifTagNumber {
	return f.tag_
}
func (f *ExifCommonField) GetType() ExifFormatId {
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

// Exifのタイプ定義　ちょうど連番なのでiotaを使用
const (
	_ = iota // 0は未使用なので捨てる
	IFD_FORMAT_UBYTE
	IFD_FORMAT_STRING
	IFD_FORMAT_USHORT
	IFD_FORMAT_ULONG
	IFD_FORMAT_URATIONAL // 5
	IFD_FORMAT_SBYTE
	IFD_FORMAT_UNDEFINED
	IFD_FORMAT_SSHORT
	IFD_FORMAT_SLONG
	IFD_FORMAT_SRATIONAL //10
	IFD_FORMAT_SINGLE
	IFD_FORMAT_DOUBLE
	IFD_FORMAT_IFD //13
)

type FormatType struct {
	BytePerFormat uint32
	FormatName string
}

var FormatTypeMap = map[ExifFormatId]FormatType{
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

func NewExifField(payload []byte, offset uint32, byteOrder binary.ByteOrder) ExifField {
	common := ExifCommonField {
		tag_ : ExifTagNumber(byteOrder.Uint16(payload[offset:offset+2])),
		type_ : ExifFormatId(byteOrder.Uint16(payload[offset+2:offset+4])),
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

