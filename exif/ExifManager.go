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
	e.Header = CreateExifHeader(payload)
	e.Directory = make(map[string]*ExifDirectory, 5)

	byteOrder := e.Header.byteOrder

	ifd0th := CreateExifDirectory(payload, e.Header.ifdOffset, byteOrder)
	e.Directory[IFD_0th_NAME] = ifd0th

	ifdExifField := ifd0th.findField(0x8769)
	if ifdExifField != nil {
		ifdExif := CreateExifDirectory(payload, byteOrder.Uint32(ifdExifField.GetData()), byteOrder)
		e.Directory[IFD_Exif_NAME] = ifdExif

		ifdInteroperabilityField := ifdExif.findField(0xA005)
		if ifdInteroperabilityField != nil {
			ifdInteroperability := CreateExifDirectory(payload, byteOrder.Uint32(ifdInteroperabilityField.GetData()), byteOrder)
			e.Directory[IFD_Interoperability_NAME] = ifdInteroperability
		}
	}

	ifdGpsField := ifd0th.findField(0x8825)
	if ifdGpsField != nil {
		ifdGps := CreateExifDirectory(payload, byteOrder.Uint32(ifdGpsField.GetData()), byteOrder)
		e.Directory[IFD_GPS_NAME] = ifdGps
	}

	if ifd0th.offsetNextIfd != 0 {
		ifd1st := CreateExifDirectory(payload, ifd0th.offsetNextIfd, byteOrder)
		e.Directory[IFD_1st_NAME] = ifd1st
	}
	return &e
}

