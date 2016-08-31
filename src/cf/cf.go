package cf

import (
	"bytes"
	"errors"
	"os"
	"runtime"
	"utils"
)

const (
	emptyValue       = int32(0x7fffffff) // Разделитель адресов в заголовке
	defaultBlockSize = 512               // Размер блока по умолчанию
)

var (
	groupBlocksFlag            = []byte{0xFF, 0xFF, 0xFF, 0x7F} // маркер группы
	addresTOCControlCharacters = []utils.Coordinates{
		{0, len(beginHeaderMarker)},
		{len(beginHeaderMarker) + 8, 1},
		{len(beginHeaderMarker) + 8 + 9, 1},
		{len(beginHeaderMarker) + 8 + 9 + 9 + 1, len(beginHeaderMarker)}}
)

// Позиции заголовков областей содержащих информацию об одном блоке данных
// прочитанные из оглавления
type addresInTOC struct {
	AttrsPos int32 // позиция заголовка атрибутов блока
	DataPos  int32 // позиция заголовка данных блока
}

// Заголовки двух областей содержащих информацию об одном блоке данных
type headersPair struct {
	attrs header // заголовок атрибутов блока
	data  header // заголовок данных блока
}

// Возвращает количество байт одного из значений адреса таблицы оглавления:
// - 4 байта (Позиция заголовка атрибутов)
// - 4 байта (Позиция заголовка данных)
// - 4 байта (Разделитель адресов - emptyValue)
func (a addresInTOC) getElementSize() int {
	return len(utils.Int32ToBytes(a.AttrsPos)) // = 4
}

// Возвращает размер одного адреса в таблице оглавления в байтах
func (a addresInTOC) getSize() int {
	return a.getElementSize() * 3 // Позиция заголовка атрибутов + Позиция заголовка данных + Разделитель адресов = 16
}

// Объект с данными файла конфигурационного файла
type ConfCf struct {
	blocksList []block
}

// Инициализация объекта данными из файла
func (c *ConfCf) InitFromFile(fileData []byte) error {

	utils.AddTextToLog(utils.LogLevel_INFO, "Начало инициализации объекта данными конфигурационного файла")

	if len(fileData) <= len(groupBlocksFlag) {
		return errors.New("Не верный размер файла конфигурационного файла")
	}

	if !bytes.Equal(fileData[:len(groupBlocksFlag)], groupBlocksFlag) {
		return errors.New("Данные не соотвествуют файлу конфигурационному файлу")
	}

	bloksHeadersPairs, err := readTOC(fileData)
	if err != nil {
		return err
	} else if len(bloksHeadersPairs) == 0 {
		return errors.New("Ошибка чтения заголовка.")
	}

	funcForInitBlocks := func(numStream int, processBlocksTotal int, countBlocks int, blockInitResultChanal chan bool) {
		for i := processBlocksTotal; i < countBlocks && i-processBlocksTotal != numStream; i++ {
			pair := bloksHeadersPairs[i]
			go c.blocksList[i].Init(fileData, pair.attrs, pair.data, blockInitResultChanal)
		}
	}

	c.initBlocks(len(bloksHeadersPairs), funcForInitBlocks)

	utils.AddTextToLog(utils.LogLevel_INFO, "-Окончание инициализации объекта данными конфигурационного файла")

	return nil
}

// Инициализация объекта данными из каталого содержащего распакованный конфигурационный блок
func (c ConfCf) InitFromCatalog(pathToDir string) (err error) {

	utils.AddTextToLog(utils.LogLevel_INFO, "Начало инициализации объекта данными конфигурационного файла из каталога")

	var fileInfos []os.FileInfo

	if fileInfos, err = utils.ReadFilesInDir(pathToDir); err != nil {
		return err
	}

	funcForInitBlocks := func(numStream int, processBlocksTotal int, countBlocks int, blockInitResultChanal chan bool) {
		for i := processBlocksTotal; i < countBlocks && i-processBlocksTotal != numStream; i++ {
			fi := fileInfos[i]
			go c.blocksList[i].InitFromFiles(pathToDir, fi, blockInitResultChanal)
		}
	}

	c.initBlocks(len(fileInfos), funcForInitBlocks)

	utils.AddTextToLog(utils.LogLevel_INFO, "-Окончание инициализации объекта данными конфигурационного файла из каталога")

	return nil
}

// Получить данные блоков конфигурационного файла
func (c ConfCf) GetData() []block {
	return c.blocksList
}

// Получить данные для записи конфигурационного файла
func (c ConfCf) GetDataForConfigFile() (fileData []byte) {

	fileData = []byte{}
	addresses := []addresInTOC{}

	for _, b := range c.blocksList {
		attrsData, data := b.GetDataForConfigFile()

		headerAndAttrs := []byte{}
		utils.AddToSlice(&headerAndAttrs,
			getHeaderForCf(len(attrsData), len(attrsData)),
			attrsData)

		headerAndData := []byte{}
		utils.AddToSlice(&headerAndData,
			getHeaderForCf(len(data), len(data)),
			data)

		pos := utils.AddToSlice(&fileData, headerAndAttrs, headerAndData)
		addresses = append(addresses,
			addresInTOC{AttrsPos: int32(pos[0]), DataPos: int32(pos[1])})
	}

	addTableOfContent(&fileData, addresses)
	return
}

// Сохранить данные блоков в файлы
func (c ConfCf) SaveBlocksToFiles(pathToDir string) {

	for _, b := range c.blocksList {
		b.WriteToFile(pathToDir)
	}
}

type funcInitBlocks func(int, int, int, chan bool)

// Инициализировать данные блоков
func (c *ConfCf) initBlocks(countBlocks int, initFunc funcInitBlocks) {

	var numStream int
	if utils.GetLogLevel() == utils.LogLevel_TRACE {
		numStream = 1 // При отладке используем только один поток, т.к. иначе сложно читать лог
	} else {
		numStream = runtime.NumCPU()
	}

	blockInitResultChanal := make(chan bool, numStream)
	resultChanal := make(chan bool)

	processBlocksTotal := 0
	processBlocksInStream := 0
	c.blocksList = make([]block, countBlocks)

	go initFunc(numStream, processBlocksTotal, countBlocks, blockInitResultChanal)
	go func() {
		for {
			select {
			case currentResult := <-blockInitResultChanal:
				if currentResult != true {
					panic("обработаны не все блоки")
				}

				processBlocksTotal++
				processBlocksInStream++

				if processBlocksInStream == numStream {
					processBlocksInStream = 0
					go initFunc(numStream, processBlocksTotal, countBlocks, blockInitResultChanal)
				}

			default:
				if processBlocksTotal == countBlocks {
					resultChanal <- true
					break
				}
			}
		}
	}()
	<-resultChanal
}

// Найти заголовки блоков оглавления (оглавление может состоять из нескольских блоков)
// @param данные файла
// @result заголовки областей оглавления, описание ошибки (если есть)
func findTOC(data []byte) (headerTOC []header, err error) {

	utils.AddTextToLog(utils.LogLevel_TRACE, "Начало поиска оглавления")

	err = errors.New("Оглавление не найдено")
	var previousHeader *header

	for i := 0; i < len(data); i++ {

		if !isHeader(data, i) {
			continue
		}

		err = nil
		utils.AddTextToLog(utils.LogLevel_TRACE, "Позиция части оглавления: "+utils.IntToString(i))
		previousHeader, err = createHeader(data, i)
		break
	}

	for true {
		if err != nil {
			headerTOC = make([]header, 0)
			break
		}

		headerTOC = append(headerTOC, *previousHeader)

		if previousHeader.getNextHeaderPosition() == 0 {
			break
		} else {
			i := previousHeader.getNextHeaderPosition()
			if i > len(data) {
				break
			}

			utils.AddTextToLog(utils.LogLevel_TRACE, "Позиция части оглавления: "+utils.IntToString(i))
			previousHeader, err = createHeader(data, i)
		}
	}

	utils.AddTextToLog(utils.LogLevel_TRACE, "-Окончание поиска оглавления")

	return
}

// Возвращает пары <заголовок атрибутов блока, заголовок данных блока> на основании оглавления файла
func readTOC(data []byte) (pairs []headersPair, err error) {

	utils.AddTextToLog(utils.LogLevel_TRACE, "Начало чтения оглавления")

	pairs = make([]headersPair, 0)
	headersTOC, err := findTOC(data)
	if err == nil {
		for _, h := range headersTOC {
			err = readPartTOC(data, h, &pairs)
			if err != nil {
				pairs = make([]headersPair, 0)
				break
			}
		}
	}

	utils.AddTextToLog(utils.LogLevel_TRACE, "-Окончание чтения оглавления. Найдено блоков: "+utils.IntToString(len(pairs)))

	return
}

// Возвращает пары <заголовок атрибутов блока, заголовок данных блока> на основании ЧАСТИ оглавления файла
func readPartTOC(data []byte, pathHeader header, pairs *[]headersPair) (err error) {

	utils.AddTextToLog(utils.LogLevel_TRACE, "Начало чтения части оглавления")

	dataTOC := getRegion(data, pathHeader)
	addres := addresInTOC{}
	addresPath := addres.getElementSize()
	bytesInElement := addres.getSize()
	elements := len(dataTOC) / bytesInElement

	utils.AddTextToLog(utils.LogLevel_TRACE, "Резмер данных оглавления: "+utils.IntToString(len(dataTOC)))
	utils.AddTextToLog(utils.LogLevel_TRACE, "Количество элементов в оглавлении: "+utils.IntToString(elements))

	for i := 0; i < elements; i++ {
		startPos := i * bytesInElement

		blockAttrs := dataTOC[startPos : startPos+addresPath]
		blockData := dataTOC[startPos+addresPath : startPos+addresPath+addresPath]

		attrsPos := int(utils.BytesToInt32(blockAttrs))
		dataPos := int(utils.BytesToInt32(blockData))

		if !isHeader(data, attrsPos) || !isHeader(data, dataPos) {
			continue
		}

		attrsHeader, err := createHeader(data, attrsPos)
		if err != nil {
			break
		}
		dataHeader, err := createHeader(data, dataPos)
		if err != nil {
			break
		}

		if !attrsHeader.isAttributesData(data) || dataHeader.getDataSize() == 0 {
			continue
		}

		var newPair headersPair
		newPair.attrs = *attrsHeader
		newPair.data = *dataHeader

		*pairs = append(*pairs, newPair)
	}

	utils.AddTextToLog(utils.LogLevel_TRACE, "-Окончание чтения части оглавления")
	return
}

// Добавляет оглавление группового блока
func addTableOfContent(fileData *[]byte, addresses []addresInTOC) {

	utils.AddTextToLog(utils.LogLevel_TRACE, "Начало добавления оглавления")

	addres := addresInTOC{}
	delimetr := utils.Int32ToBytes(emptyValue)
	countBlocks := len(addresses)
	valuableTOCSize := addres.getSize() * countBlocks
	fullTOCSize := valuableTOCSize

	if fullTOCSize < defaultBlockSize {
		fullTOCSize = defaultBlockSize
	}

	data := []byte{}
	utils.AddToSlice(&data,
		groupBlocksFlag,                              // маркер группы
		utils.Int32ToBytes(defaultBlockSize),         // размер блока по умолчанию
		[]byte{0x00, 0x00, 0x00, 0x00},               // unknown
		[]byte{0x00, 0x00, 0x00, 0x00},               // unknown
		getHeaderForCf(valuableTOCSize, fullTOCSize)) // Заголовок области оглавления

	beginFileDataLen := int32(len(data) + fullTOCSize)

	for _, v := range addresses {
		utils.AddToSlice(&data,
			utils.Int32ToBytes(v.AttrsPos+beginFileDataLen), // позиция заголовка атрибутов блока
			utils.Int32ToBytes(v.DataPos+beginFileDataLen),  // позиция заголовка данных блока
			delimetr) // макер разделителя адресов
	}

	emptyData := []byte{}
	if valuableTOCSize < fullTOCSize {
		emptyData = make([]byte, fullTOCSize-valuableTOCSize)
	}

	data = append(data, emptyData...)
	*fileData = append(data, (*fileData)...)

	utils.AddTextToLog(utils.LogLevel_TRACE, "-Окончание добавления оглавления")
}

// Возвращает область данных по данным заголовка
func getRegion(data []byte, h header) []byte {
	return data[h.getDataPosition() : h.getDataPosition()+h.getDataSize()]
}

// Дополнить данные блока до размера блока по умолчанию
func prepareDataForConfigFile(blockData *[]byte) {

	dataValueSize := len(*blockData)
	totalDataSize := dataValueSize
	if dataValueSize < defaultBlockSize {
		totalDataSize = defaultBlockSize
	}

	if dataValueSize < totalDataSize {
		newData := make([]byte, totalDataSize-dataValueSize)
		*blockData = append(*blockData, newData...)
	}
}
