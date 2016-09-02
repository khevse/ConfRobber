package utils

import (
	"bytes"
	"flag"
	"os"
	"path"
	"testing"
)

var (
	pathToTarget string
)

func TestMain(m *testing.M) {

	currentDir, _ := GetPathToCurrentDir()
	pathToTarget = path.Join(GetParentDir(currentDir), "target dir")

	RemoveIfExist(pathToTarget)
	os.Mkdir(pathToTarget, os.ModeDir)
	InitLogger(pathToTarget, LogLevel_TRACE)

	flag.Parse()
	os.Exit(m.Run())

	RemoveIfExist(pathToTarget)
}

func TestEqual(m *testing.T) {

	a := []byte{'0', '1', '2'}
	b := a[:]
	c := a[1:]

	if bytes.Equal(a, b) != true {
		m.Error(a, " != ", b)
	}

	if bytes.Equal(a, c) == true {
		m.Error("Equal: ", a, " == ", c)
	}
}

func TestGetValues(t *testing.T) {

	data := []byte{1, 2, 3, 4, 5, 6, 7}
	standart := [][]byte{
		[]byte{3, 4},
		[]byte{5},
		[]byte{7}}

	positions := []Coordinates{
		{0, 2}, // 3,4
		{2, 1}, // 5
		{4, 1}} // 7

	list := GetValues(data, 2, positions)

	for i, val := range list {

		if bytes.Equal(val, standart[i]) != true {
			t.Error("GetValues: ", standart[i], " != ", val)
		}
	}
}

func TestHexToInt(t *testing.T) {

	data := []byte{'0', '0', '0', '0', '0', '0', 'a', 'c'}
	val := HexToInt(data)

	if val != 172 {
		t.Error("HexToInt: ", data, " = ", 172, ", got ", val)
	}
}

func TestGetParentDir(t *testing.T) {

	standart := "c:"
	fullPath := standart + string(os.PathSeparator) + "f"

	v := GetParentDir(fullPath)

	if v != standart {
		t.Error("GetParentDir error", v, " != ", standart)
	}
}

func TestUtf(t *testing.T) {

	s := "abc"
	u8 := []rune(s)
	u16 := Utf8ToUtf16(u8)

	standart := []byte{'a', 0x00, 'b', 0x00, 'c', 0x00}
	if !bytes.Equal(standart, u16) {
		t.Error("Utf8ToUtf16 error", standart, " != ", u16)
	}

	res := Utf16ToUtf8(u16)
	if res != s {
		t.Error("Utf16ToUtf8", s, "!=", res)
	}
}

func TestAddToSlice(t *testing.T) {

	a := []byte{0x01, 0x02}
	b := []byte{0x03, 0x04}

	res := []byte{0x00}
	positions := AddToSlice(&res, a, b)

	standart := []byte{0x00}
	standart = append(standart, a...)
	standart = append(standart, b...)

	if !bytes.Equal(res, standart) {
		t.Error("AddToSlice", standart, "!=", res)
	}

	if positions[0] != 1 {
		t.Error("AddToSlice positions[0] ", 1, "!=", positions[0])
	}

	if positions[1] != 3 {
		t.Error("AddToSlice positions[0] ", 2, "!=", positions[1])
	}
}
