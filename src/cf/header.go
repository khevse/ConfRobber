package cf

import (
	"bytes"
	"errors"
	"utils"
)

const (
	valueHeaderSize = 27 // Размер полезных данных заголовка
	fullHeaderSize  = 31 // Размер заголовка вместе с маркерами
)

var (
	// Маркеры начала и окончания заголовка
	beginHeaderMarker = []byte{'\r', '\n'}
	endHeaderMarker   = []byte{'\r', '\n'}
	emptyHeaderValue  = []byte{'7', 'f', 'f', 'f', 'f', 'f', 'f', 'f', ' '}
	space             = []byte{0x20}

	attrsControlCharacters = []utils.Coordinates{
		{7, 1},             // последний символ в значении даты модификации
		{7 + 8, 1},         // последний символ в значении даты создания
		{7 + 8 + 4, 1},     // последний символ в значении типа блока
		{7 + 8 + 4 + 2, 1}} // второй символ в имени блока

	headerControlCharacters = []utils.Coordinates{
		{0, len(beginHeaderMarker)},
		{len(beginHeaderMarker) + 8, 1},
		{len(beginHeaderMarker) + 8 + 9, 1},
		{len(beginHeaderMarker) + 8 + 9 + 9, 1},
		{len(beginHeaderMarker) + 8 + 9 + 9 + 1, len(endHeaderMarker)}}
)

// В текущей позиции начинаются данные заголовка
// Пример: "\r\n000000a3 00000200 7fffffff \r\n" (всегда 31 символ)
//
// @param - данные
// @param - текущая позиция, в которой ожидаем увидеть данные заголовка
func isHeader(data []byte, startPos int) bool {

	values := utils.GetValues(data, startPos, headerControlCharacters)
	if len(values) != len(headerControlCharacters) {
		return false
	}

	return (bytes.Equal(values[0], beginHeaderMarker) &&
		bytes.Equal(values[1], space) &&
		bytes.Equal(values[2], space) &&
		bytes.Equal(values[3], space) &&
		bytes.Equal(values[4], endHeaderMarker))
}

// Возвращает заголовок области данных
// @param - полезный размер области
// @param - полный размер области
func getHeaderForCf(valuableRegionSize int, fullRegionSize int) []byte {

	delim := []byte{' '}
	fullSize := utils.Int32ToHexToBytes(int32(fullRegionSize))
	valuableSize := utils.Int32ToHexToBytes(int32(valuableRegionSize))

	data := []byte{}
	utils.AddToSlice(&data,
		beginHeaderMarker,
		valuableSize,
		delim,
		fullSize,
		delim,
		emptyHeaderValue,
		endHeaderMarker)
	return data
}

// Создать заголовок описывающий область файла
// @param - все данные файла
// @param - поизиция, с которой выполнять чтение
//
// @result - заголовок и информация об ошибке если таковая имеется
func createHeader(data []byte, headerPosition int) (h *header, err error) {

	pos := headerPosition + len(beginHeaderMarker)
	values := bytes.Fields(data[pos : pos+valueHeaderSize-1]) // "000000ac 00000200 00000200" => [000000ac, 00000200, 00000200]

	err = nil
	if len(values) != 3 {
		err = errors.New("Ошибка чтения значений заголовка.")
	}

	if err == nil {
		nextHeaderPos := int(utils.HexToInt(values[2]))
		if bytes.Equal(utils.Int32ToBytes(int32(nextHeaderPos)), utils.Int32ToBytes(emptyValue)) {
			nextHeaderPos = 0
		}

		h = &header{regionPosition: headerPosition + fullHeaderSize,
			valuableRegionSize: int(utils.HexToInt(values[0])), // 000000ac -> 172
			totalRegionSize:    int(utils.HexToInt(values[1])), // 00000200 -> 512
			nextHeaderPosition: nextHeaderPos}                  // 00000200 -> 512
	}

	return
}

// Заголовок области, содержащий координаты и размеры в конфигурационном файле 1С
//   - пример заголовка с маркерами:    "\r\n000000ac 00000200 7fffffff \r\n"
//   - пример заголовка без маркеров: "000000ac 00000200 00000200 "
//  , где:
//   - первое значение: "000000ac" - размер полезных данных области
//   - второе значение: "00000200" - полный размер данных области
//   - третье значение: "00000200" - позиция продолжения данных области. Если содержит значение "7fffffff" - продолжение отсутствует
type header struct {
	valuableRegionSize int // размер полезных данных области
	totalRegionSize    int // полный размер области
	regionPosition     int // позиция с которой начинаются данные
	nextHeaderPosition int // позиция заголовка продолжения данных
}

// Возвращает позицию области в файле
func (h header) getDataPosition() int {
	return h.regionPosition
}

// Возвращает размер полезных данных области
func (h header) getDataSize() (size int) {

	if h.valuableRegionSize > h.totalRegionSize || h.valuableRegionSize == 0 {
		size = h.totalRegionSize
	} else {
		size = h.valuableRegionSize
	}
	return
}

// Возвращает позицию продолжения данных
func (h header) getNextHeaderPosition() int {
	return h.nextHeaderPosition
}

// Заголовок принадлежит области с данными атрибутов блока
// @param - все данные конфигурационного файла 1С
func (h header) isAttributesData(data []byte) bool {

	if h.regionPosition == 0 {
		return false
	}

	const minBlockSize = 21 // второй символ в имени блока
	if len(data)-h.regionPosition <= minBlockSize {
		return false
	}

	values := utils.GetValues(data, h.regionPosition, attrsControlCharacters)

	return (bytes.IndexByte(values[0], 0x00) == 0 &&
		bytes.IndexByte(values[1], 0x00) == 0 &&
		bytes.IndexByte(values[2], 0x00) == 0 &&
		bytes.IndexByte(values[3], 0x00) == 0)
}
