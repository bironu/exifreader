package exif

type ExifManager struct {
	Header    *ExifHeader
	Directory map[string]*ExifDirectory
	payload []byte
}

const IFD_0th_NAME = "0th"
const IFD_1st_NAME = "1st"
const IFD_Exif_NAME = "Exif"
const IFD_Interoperability_NAME = "Interoperability"
const IFD_GPS_NAME = "GPS"

func CreateExifManager(payload []byte) *ExifManager {
	var e ExifManager
	e.payload = payload
	e.Header = NewExifHeader(payload)
	e.Directory = make(map[string]*ExifDirectory, 5)

	byteOrder := e.Header.byteOrder

	ifd0th := NewExifDirectory(payload, e.Header.ifdOffset, byteOrder)
	e.Directory[IFD_0th_NAME] = ifd0th

	ifdExifField := ifd0th.findField(IFD_EXIF_TAG_NUMBER)
	if ifdExifField != nil {
		ifdExif := NewExifDirectory(payload, byteOrder.Uint32(ifdExifField.GetData()), byteOrder)
		e.Directory[IFD_Exif_NAME] = ifdExif

		ifdInteroperabilityField := ifdExif.findField(IFD_INTEROPERABILITY_TAG_NUMBER)
		if ifdInteroperabilityField != nil {
			ifdInteroperability := NewExifDirectory(payload, byteOrder.Uint32(ifdInteroperabilityField.GetData()), byteOrder)
			e.Directory[IFD_Interoperability_NAME] = ifdInteroperability
		}
	}

	ifdGpsField := ifd0th.findField(IFD_GPS_TAG_NUMBER)
	if ifdGpsField != nil {
		ifdGps := NewExifDirectory(payload, byteOrder.Uint32(ifdGpsField.GetData()), byteOrder)
		e.Directory[IFD_GPS_NAME] = ifdGps
	}

	if ifd0th.offsetNextIfd != 0 {
		ifd1st := NewExifDirectory(payload, ifd0th.offsetNextIfd, byteOrder)
		e.Directory[IFD_1st_NAME] = ifd1st
	}
	return &e
}

