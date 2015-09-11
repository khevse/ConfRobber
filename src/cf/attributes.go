package cf

import (
	"utils"
)

const (

	// Типы блока
	groupTypeModule   = 740   // Модуль ( Последовательный групповой блок: атрибуты, данные, атрибуты, данные, ...)
	groupTypeForm     = 689   // Форма ( Зеркальный групповой блок: атрибуты1, атрибуты2, данные2, данные1 )
	groupTypeConfig   = 686   // Заголовок конфигурационного файла
	groupTypeNoModule = 84846 // Без модуля (Последовательный групповой блок: атрибуты, данные, атрибуты, данные, ...)

)

// Создает объект с данными атрибутов блока
func createAttrs(data []byte) *attributes {

	a := attributes{}
	a.creationDate = utils.BytesToInt64(data[:8])
	a.modificationDate = utils.BytesToInt64(data[8:16])
	a.groupType = int32(utils.BytesToInt32(data[16:20]))

	name_utf16 := data[20:] // n.a.m.e... ,где '.' == 0x00

	for i, symbol := range name_utf16 {
		next_symbol := byte(0x00)
		if i+1 <= len(name_utf16) {
			next_symbol = name_utf16[i+1:][0]
		}

		if symbol == 0x00 && next_symbol == 0x00 {
			break
		}

		if i%2 == 0 {
			a.name += string(symbol)
		}
	}

	return &a
}

// Атрибуты блока
type attributes struct {
	creationDate     int64  // дата создания
	modificationDate int64  // дата модификации
	groupType        int32  // тип блока
	name             string // наименование блока
}

// Возвращает данные атрибутов для конфигурационного файла
func (a attributes) getData() []byte {

	var unknowBlock []byte
	if len(a.name) == 0 {
		unknowBlock = utils.Int32ToBytes(emptyValue)
	} else {
		unknowBlock = []byte{0x00, 0x00, 0x00, 0x00}
	}

	data := []byte{}
	utils.AddToSlice(&data,
		utils.Int64ToBytes(a.creationDate),
		utils.Int64ToBytes(a.modificationDate),
		utils.Int32ToBytes(a.groupType),
		utils.Utf8ToUtf16([]rune(a.name)),
		unknowBlock)

	return data
}
