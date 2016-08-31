package zlibwrapper

import (
	"bytes"
	"testing"
)

func TestCompress(t *testing.T) {

	source := []byte("hello")
	val, err := Compress(source)
	if err != nil {
		t.Error(err.Error())
	}

	etalon := []byte{202, 72, 205, 201, 201, 7}
	if !bytes.Equal(val, etalon) {
		t.Errorf("'%v' != '%v'", etalon, val)
	}
}

func TestDecompress(t *testing.T) {

	source := []byte{202, 72, 205, 201, 201, 7}
	val, err := Decompress(source)
	if err != nil {
		t.Error(err.Error())
	}

	etalon := []byte("hello")
	if !bytes.Equal(val, etalon) {
		t.Errorf("'%v' != '%v'", etalon, val)
	}
}

func TestCompressAndDecompress(t *testing.T) {

	source := []byte("hello")

	compressData, err := Compress(source)
	if err != nil {
		t.Error(err.Error())
	}

	decompressData, err := Decompress(compressData)
	if err != nil {
		t.Error(err.Error())
	}

	if !bytes.Equal(decompressData, source) {
		t.Errorf("'%v' != '%v'", source, decompressData)
	}
}
