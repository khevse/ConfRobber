package utils

type Coordinates struct {
	Position int
	Length   int
}

// Возвращает колекцию блоков полученных из указанных позиций
//
// @param данные из которых будет читать части
// @param начальная позиция чтения данных
// @param координаты блоков которые нужно получить
func GetValues(data []byte, startPos int, positions []Coordinates) [][]byte {

	values := [][]byte{}

	for _, v := range positions {
		valLen := v.Length
		startDataPos := startPos + v.Position
		endDataPos := startDataPos + valLen

		if endDataPos <= len(data) {
			currentVal := make([]byte, valLen)
			copy(currentVal, data[startDataPos:endDataPos])

			values = append(values, []byte{})
			i := len(values) - 1
			values[i] = append(values[i], currentVal...)

		} else {
			break
		}
	}

	return values
}

// Добавляет в коллекцию несколько наборов данных и
// возвращает поизиции, в которые были добавлены блоки
func AddToSlice(s *[]byte, args ...[]byte) []int {

	dataLen := 0
	for _, v := range args {
		dataLen += len(v)
	}

	positions := make([]int, len(args))
	sourceDataLen := len(*s)

	startPos := 0
	newData := make([]byte, dataLen)

	for i, v := range args {
		valLen := len(v)
		copy(newData[startPos:startPos+valLen], v)

		positions[i] = startPos + sourceDataLen
		startPos += valLen
	}

	*s = append(*s, newData...)

	return positions
}
