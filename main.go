package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"

	mp4 "github.com/abema/go-mp4"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	output, err := os.Create("output.mp4")
	// require.NoError(t, err)
	checkError(err)
	defer output.Close()

	w := mp4.NewWriter(output)

	// start ftyp
	_, err = w.StartBox(&mp4.BoxInfo{Type: mp4.BoxTypeFtyp()})
	checkError(err)

	// require.NoError(t, err)
	// assert.Equal(t, uint64(0), bi.Offset)
	// assert.Equal(t, uint64(8), bi.Size)

	ftyp := &mp4.Ftyp{
		MajorBrand:   [4]byte{'m', 'p', '4', '2'},
		MinorVersion: 0x1,
		CompatibleBrands: []mp4.CompatibleBrandElem{
			{CompatibleBrand: [4]byte{'m', 'p', '4', '1'}},
			{CompatibleBrand: [4]byte{'m', 'p', '4', '2'}},
		},
	}
	_, err = mp4.Marshal(w, ftyp, mp4.Context{})
	checkError(err)

	// require.NoError(t, err)

	// end ftyp
	_, err = w.EndBox()
	checkError(err)

	// require.NoError(t, err)
	// assert.Equal(t, uint64(0), bi.Offset)
	// assert.Equal(t, uint64(24), bi.Size)

	// start moov
	_, err = w.StartBox(&mp4.BoxInfo{Type: mp4.BoxTypeMoov()})
	checkError(err)

	// require.NoError(t, err)
	// assert.Equal(t, uint64(24), bi.Offset)
	// assert.Equal(t, uint64(8), bi.Size)

	// copy
	err = w.CopyBox(bytes.NewReader([]byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x0a,
		'u', 'd', 't', 'a',
		0x01, 0x02, 0x03, 0x04,
		0x05, 0x06, 0x07, 0x08,
	}), &mp4.BoxInfo{Offset: 6, Size: 15})
	checkError(err)

	// require.NoError(t, err)

	// start trak
	_, err = w.StartBox(&mp4.BoxInfo{Type: mp4.BoxTypeTrak()})
	checkError(err)

	// require.NoError(t, err)
	// assert.Equal(t, uint64(47), bi.Offset)
	// assert.Equal(t, uint64(8), bi.Size)

	// start tkhd
	_, err = w.StartBox(&mp4.BoxInfo{Type: mp4.BoxTypeTkhd()})
	checkError(err)

	// require.NoError(t, err)
	// assert.Equal(t, uint64(55), bi.Offset)
	// assert.Equal(t, uint64(8), bi.Size)

	_, err = mp4.Marshal(w, &mp4.Tkhd{
		CreationTimeV0:     1,
		ModificationTimeV0: 2,
		TrackID:            3,
		DurationV0:         4,
		Layer:              5,
		AlternateGroup:     6,
		Volume:             7,
		Width:              8,
		Height:             9,
	}, mp4.Context{})
	checkError(err)

	// require.NoError(t, err)

	// end tkhd
	_, err = w.EndBox()
	checkError(err)

	// require.NoError(t, err)
	// assert.Equal(t, uint64(55), bi.Offset)
	// assert.Equal(t, uint64(92), bi.Size)

	// end trak
	_, err = w.EndBox()
	checkError(err)

	// require.NoError(t, err)
	// assert.Equal(t, uint64(47), bi.Offset)
	// assert.Equal(t, uint64(100), bi.Size)

	// end moov
	_, err = w.EndBox()
	checkError(err)

	// require.NoError(t, err)
	// assert.Equal(t, uint64(24), bi.Offset)
	// assert.Equal(t, uint64(123), bi.Size)

	// update ftyp
	_, err = w.Seek(8, io.SeekStart)
	checkError(err)

	// require.NoError(t, err)
	// assert.Equal(t, int64(8), n)
	ftyp.CompatibleBrands[1].CompatibleBrand = [4]byte{'E', 'F', 'G', 'H'}
	_, err = mp4.Marshal(w, ftyp, mp4.Context{})
	checkError(err)

	// require.NoError(t, err)

	_, err = output.Seek(0, io.SeekStart)
	checkError(err)

	// require.NoError(t, err)
	_, err = ioutil.ReadAll(output)
	checkError(err)

	// require.NoError(t, err)
	// assert.Equal(t, []byte{
	// 	// ftyp
	// 	0x00, 0x00, 0x00, 0x18, // size
	// 	'f', 't', 'y', 'p', // type
	// 	'a', 'b', 'e', 'm', // major brand
	// 	0x12, 0x34, 0x56, 0x78, // minor version
	// 	'a', 'b', 'c', 'd', // compatible brand
	// 	'E', 'F', 'G', 'H', // compatible brand
	// 	// moov
	// 	0x00, 0x00, 0x00, 0x7b, // size
	// 	'm', 'o', 'o', 'v', // type
	// 	// udta (copy)
	// 	0x00, 0x00, 0x00, 0x0a,
	// 	'u', 'd', 't', 'a',
	// 	0x01, 0x02, 0x03, 0x04,
	// 	0x05, 0x06, 0x07,
	// 	// trak
	// 	0x00, 0x00, 0x00, 0x64, // size
	// 	't', 'r', 'a', 'k', // type
	// 	// tkhd
	// 	0x00, 0x00, 0x00, 0x5c, // size
	// 	't', 'k', 'h', 'd', // type
	// 	0,                // version
	// 	0x00, 0x00, 0x00, // flags
	// 	0x00, 0x00, 0x00, 0x01, // creation time
	// 	0x00, 0x00, 0x00, 0x02, // modification time
	// 	0x00, 0x00, 0x00, 0x03, // track ID
	// 	0x00, 0x00, 0x00, 0x00, // reserved
	// 	0x00, 0x00, 0x00, 0x04, // duration
	// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // reserved
	// 	0x00, 0x05, // layer
	// 	0x00, 0x06, // alternate group
	// 	0x00, 0x07, // volume
	// 	0x00, 0x00, // reserved
	// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	// 	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // matrix
	// 	0x00, 0x00, 0x00, 0x08, // width
	// 	0x00, 0x00, 0x00, 0x09, // height
	// }, bin)
}
