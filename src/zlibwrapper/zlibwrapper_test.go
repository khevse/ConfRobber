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

func TestMultiThread(t *testing.T) {

	processCount := 200

	resulChanel := make(chan bool)
	processChanel := make(chan bool, processCount)

	source := []byte("hello")

	for i := 0; i < processCount; i++ {
		go func(cn chan bool) {
			for p := 0; p < 1000; p++ {
				compressData, _ := Compress(source)
				decompressData, _ := Decompress(compressData)

				if !bytes.Equal(decompressData, source) {
					t.Errorf("'%v' != '%v'", source, decompressData)
				}
			}
			cn <- true
		}(processChanel)
	}

	go func() {
		for {
			select {
			case <-processChanel:
				{
					processCount--
					if processCount == 0 {
						resulChanel <- true
					}
				}
			}
		}
	}()

	<-resulChanel
}
