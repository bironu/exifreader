package main

import (
	"fmt"
	"os"
	"bironu/exifreader/exif"
	"bytes"
	"encoding/binary"
	"io"
)

const MARKER_SIZE = 2
const LENGTH_SIZE = 2
const MARKER_APP0 = 0xFFE0
const MARKER_APP1 = 0xFFE1
const MARKER_APP15 = 0xFFEF
const MARKER_SOI = 0xFFD8


func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ExifReader.exe <jpg_file>")
		return
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		// Openエラー処理
		fmt.Printf("%s file open failed.\n", os.Args[1])
		panic(err)
	}
	defer file.Close()

	// JPEGファイル判定
	{
		var soi uint16
		err = binary.Read(file, binary.BigEndian, &soi)
		if err != nil || soi != MARKER_SOI {
			fmt.Printf("%s is not jpeg file. soi = %04X\n", os.Args[1], soi)
			return
		}
	}

	var manager *exif.ExifManager
	// Appセグメント読み込み
	for {
		// セグメントマーカー読み込み
		var marker uint16
		err = binary.Read(file, binary.BigEndian, &marker)
		if err != nil || marker < MARKER_APP0 || MARKER_APP15 < marker {
			// Appセグメント以外が来たら抜ける
			break
		}

		var length uint16
		err = binary.Read(file, binary.BigEndian, &length)
		if err != nil {
			fmt.Printf("segment length read failed. marker = %04X\n", marker)
			return
		}
		length -= 2 // lengthの分2byteをあらかじめ差っ引く

		if marker == MARKER_APP1 {
			// App1セグメントなら
			segmentBuf := make([]byte, length)
			err = binary.Read(file, binary.BigEndian, segmentBuf)
			if err != nil {
				fmt.Printf("%s file read error. err = %v", os.Args[1], err)
				return
			}

			exifIdentifiers := [6]byte{0x45, 0x78, 0x69, 0x66, 0x00, 0x00} // "Exif\0\0"
			identifierBuf := segmentBuf[:len(exifIdentifiers)]
			if bytes.Equal(exifIdentifiers[0:], identifierBuf) {
				// Exifデータだったら読み込んでループ抜ける
				manager = exif.CreateExifManager(segmentBuf[len(exifIdentifiers):])
				break
			}
		} else {
			// 読み込む必要のないセグメントなら
			_, err := file.Seek(int64(length), io.SeekCurrent) // 現在地からセグメント分飛ばし
			if err != nil {
				fmt.Printf("%s file seek error. err = %v", os.Args[1], err)
				return
			}
		}
	}

	if manager != nil {
		printExifInfo(exif.IFD_0th_NAME, manager.Directory[exif.IFD_0th_NAME])
		printExifInfo(exif.IFD_Exif_NAME, manager.Directory[exif.IFD_Exif_NAME])
		printExifInfo(exif.IFD_Interoperability_NAME, manager.Directory[exif.IFD_Interoperability_NAME])
		printExifInfo(exif.IFD_GPS_NAME, manager.Directory[exif.IFD_GPS_NAME])
		printExifInfo(exif.IFD_1st_NAME, manager.Directory[exif.IFD_1st_NAME])
	} else {
		fmt.Println("exif segment not found!")
	}
}

func printExifInfo(name string, d *exif.ExifDirectory) {
	if d == nil {
		return
	}

	fmt.Printf("[%s IFD]\n", name)
	for _, v := range d.GetFieldList() {
		fmt.Println(v)
	}
	fmt.Println()
}