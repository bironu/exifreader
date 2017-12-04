package exif

import (
	"fmt"
	"bytes"
	"encoding/binary"
)

type ExifUByteField struct {
	*ExifCommonField
}
func (f *ExifUByteField) String() string {
	result := f.ExifCommonField.String()
	buf := bytes.NewBuffer(f.GetData())
	v := make([]HexByte, f.GetCount())
	binary.Read(buf, f.byteOrder_, v)
	result += fmt.Sprint(v)
	return result
}

type ExifStringField struct {
	*ExifCommonField
}
func (f *ExifStringField) String() string {
	result := f.ExifCommonField.String()
	result += "\"" + string(f.GetData()) + "\""
	return result
}

type ExifUShortField struct {
	*ExifCommonField
}
func (f *ExifUShortField) String() string {
	result := f.ExifCommonField.String()
	buf := bytes.NewBuffer(f.GetData())
	v := make([]uint16, f.GetCount())
	binary.Read(buf, f.byteOrder_, v)
	result += fmt.Sprint(v)
	return result
}

type ExifULongField struct {
	*ExifCommonField
}
func (f *ExifULongField) String() string {
	result := f.ExifCommonField.String()
	buf := bytes.NewBuffer(f.GetData())
	v := make([]uint32, f.GetCount())
	binary.Read(buf, f.byteOrder_, v)
	result += fmt.Sprint(v)
	return result
}

type ExifURationalField struct {
	*ExifCommonField
}
func (f *ExifURationalField) String() string {
	result := f.ExifCommonField.String()
	buf := bytes.NewBuffer(f.GetData())
	v := make([]URational, f.GetCount())
	for i := uint32(0); i < f.GetCount(); i++ {
		binary.Read(buf, f.byteOrder_, &v[i].numer)
		binary.Read(buf, f.byteOrder_, &v[i].denom)
	}
	result += fmt.Sprint(v)
	return result
}

type ExifSByteField struct {
	*ExifCommonField
}
func (f *ExifSByteField) String() string {
	result := f.ExifCommonField.String()
	buf := bytes.NewBuffer(f.GetData())
	v := make([]int8, f.GetCount())
	binary.Read(buf, f.byteOrder_, v)
	result += fmt.Sprint(v)
	return result
}

type ExifUndefinedField struct {
	*ExifCommonField
}
func (f *ExifUndefinedField) String() string {
	result := f.ExifCommonField.String()
	buf := bytes.NewBuffer(f.GetData())
	v := make([]HexByte, f.GetCount())
	binary.Read(buf, f.byteOrder_, v)
	result += fmt.Sprint(v)
	return result
}

type ExifSShortField struct {
	*ExifCommonField
}
func (f *ExifSShortField) String() string {
	result := f.ExifCommonField.String()
	buf := bytes.NewBuffer(f.GetData())
	v := make([]int16, f.GetCount())
	binary.Read(buf, f.byteOrder_, v)
	result += fmt.Sprint(v)
	return result
}

type ExifSLongField struct {
	*ExifCommonField
}
func (f *ExifSLongField) String() string {
	result := f.ExifCommonField.String()
	buf := bytes.NewBuffer(f.GetData())
	v := make([]int32, f.GetCount())
	binary.Read(buf, f.byteOrder_, v)
	result += fmt.Sprint(v)
	return result
}

type ExifSRationalField struct {
	*ExifCommonField
}
func (f *ExifSRationalField) String() string {
	result := f.ExifCommonField.String()
	buf := bytes.NewBuffer(f.GetData())
	v := make([]SRational, f.GetCount())
	for i := uint32(0); i < f.GetCount(); i++ {
		binary.Read(buf, f.byteOrder_, &v[i].numer)
		binary.Read(buf, f.byteOrder_, &v[i].denom)
	}
	result += fmt.Sprint(v)
	return result
}

type ExifSingleFloatField struct {
	*ExifCommonField
}
func (f *ExifSingleFloatField) String() string {
	result := f.ExifCommonField.String()
	buf := bytes.NewBuffer(f.GetData())
	v := make([]float32, f.GetCount())
	binary.Read(buf, f.byteOrder_, v)
	result += fmt.Sprint(v)
	return result
}

type ExifDoubleFloatField struct {
	*ExifCommonField
}
func (f *ExifDoubleFloatField) String() string {
	result := f.ExifCommonField.String()
	buf := bytes.NewBuffer(f.GetData())
	v := make([]float64, f.GetCount())
	binary.Read(buf, f.byteOrder_, v)
	result += fmt.Sprint(v)
	return result
}

type ExifIFDField struct {
	*ExifCommonField
}
func (f *ExifIFDField) String() string {
	result := f.ExifCommonField.String()
	buf := bytes.NewBuffer(f.GetData())
	v := make([]HexByte, f.GetCount())
	binary.Read(buf, f.byteOrder_, v)
	result += fmt.Sprint(v)
	return result
}


type HexByte byte
func (b HexByte) String() string {
	return fmt.Sprintf("%02X", byte(b))
}

type URational struct {
	numer uint32 // 分子
	denom uint32 // 分母
}
func (ur URational) String() string {
	return fmt.Sprintf("%d/%d", ur.numer, ur.denom)
}

type SRational struct {
	numer int32 // 分子
	denom int32 // 分母
}
func (sr SRational) String() string {
	return fmt.Sprintf("%d/%d", sr.numer, sr.denom)
}
