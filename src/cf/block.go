package cf

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"
	"utils"
)

const (
	blockType_unknown  = 0 // неизвестный тип блока
	blockType_simply   = 1 // простой тип блока
	blockType_multiple = 2 // составной тип блока
)

// Объект вложенного блока
type subBlock struct {
	Attrs attributes
	Data  []byte
}

// Объект с данными блока
type block struct {
	blockType uint       // тип блока
	attrs     attributes // атрибуты блока
	data      []byte     // данные блока (если блок не составной)
	subBlocks []subBlock // подчиненные блоки (если это составной блок)
}

// Инициализировать блок данными из файла
// @param данные файла
// @param заголовок атрибутов блока
// @param заголовок данных блока
// @param канал в который будет записана инфорамция что блок инициализирован
func (b *block) Init(sourceData []byte, attrsHeader header, dataHeader header, c chan<- bool) {

	utils.AddTextToLog(utils.LogLevel_TRACE, "Начало инициализации блока")

	var blockData []byte
	getAttrsAndData(sourceData, attrsHeader, dataHeader, &b.attrs, &blockData)
	sourceLen := len(blockData)
	if sourceLen > 0 {
		utils.ZlibUncompress(blockData, &blockData)
	}

	if sourceLen > 0 && len(blockData) == 0 {
		panic("Ошибка распаковки блока: " + b.attrs.name)
	}

	headersTOC, err := findTOC(blockData)

	if err == nil &&
		len(headersTOC) > 0 &&
		bytes.Equal(groupBlocksFlag, blockData[:len(groupBlocksFlag)]) {

		b.blockType = blockType_multiple

		subBlocksHeadersPairs, err := readTOC(blockData)
		if err != nil {
			panic("Ошибка разбора составного блока: " + err.Error())
		}

		b.subBlocks = make([]subBlock, len(subBlocksHeadersPairs))

		for i, p := range subBlocksHeadersPairs {
			getAttrsAndData(blockData, p.attrs, p.data, &b.subBlocks[i].Attrs, &b.subBlocks[i].Data)
		}
	} else {
		b.data = blockData
		b.blockType = blockType_simply
	}

	utils.AddTextToLog(utils.LogLevel_TRACE, "Окончание инициализации блока")
	c <- true
}

// Инициализировать блок из данных сохраненных блоков в файлы
func (b *block) InitFromFiles(pathToDir string, fileInfo os.FileInfo, c chan<- bool) {

	utils.AddTextToLog(utils.LogLevel_TRACE, "Начало инициализации блока")

	var err error
	b.attrs.name = fileInfo.Name()
	currentTime := time.Now().UnixNano() / int64(time.Millisecond)
	fullPathToFile := path.Join(pathToDir, b.attrs.name)

	funcGenerateErr := func(errorNum int) {
		err = errors.New(fmt.Sprintf("Ошибка инициализации блока из даных файла №%d: %s", errorNum, err.Error()))
		utils.AddTextToLog(utils.LogLevel_ERROR, err.Error())
		panic(err.Error())
	}

	if fileInfo.IsDir() {
		b.blockType = blockType_multiple

		var files []os.FileInfo
		if files, err = utils.ReadFilesInDir(fullPathToFile); err != nil {
			funcGenerateErr(1)
		}

		b.subBlocks = make([]subBlock, len(files))
		for i, fi := range files {
			b.subBlocks[i].Attrs.name = fi.Name()
			b.subBlocks[i].Attrs.creationDate = currentTime
			b.subBlocks[i].Attrs.modificationDate = b.subBlocks[i].Attrs.creationDate

			pathToSubBlock := path.Join(fullPathToFile, fi.Name())
			if b.subBlocks[i].Data, err = ioutil.ReadFile(pathToSubBlock); err != nil {
				funcGenerateErr(2)
			}
		}

	} else {
		b.blockType = blockType_simply

		if b.data, err = ioutil.ReadFile(fullPathToFile); err != nil {
			funcGenerateErr(3)
		}

	}

	b.attrs.creationDate = currentTime
	b.attrs.modificationDate = b.attrs.creationDate
	if b.IsForm() {
		b.attrs.groupType = groupTypeForm
	} else if b.IsModule() {
		b.attrs.groupType = groupTypeModule
	} else {
		b.attrs.groupType = groupTypeNoModule
	}

	utils.AddTextToLog(utils.LogLevel_TRACE, "Окончание инициализации блока")
	c <- true
}

// Возвращает данные блока включая вложенные блоки.
// Если тип блока простой, то коллекция будет содержать только один элемент
func (b block) GetData() []subBlock {

	countSubBlocks := 0
	if b.blockType == blockType_simply {
		countSubBlocks = 1
	} else {
		countSubBlocks = len(b.subBlocks)
	}

	blocks := make([]subBlock, countSubBlocks)

	if b.blockType == blockType_simply {
		blocks[0].Attrs = b.attrs
		blocks[0].Data = b.data
	} else {

		for i, v := range b.subBlocks {
			blocks[i].Attrs = v.Attrs
			blocks[i].Data = v.Data
		}
	}

	return blocks
}

func (b block) GetDataForConfigFile() (attrsForCf []byte, dataForCf []byte) {

	if b.blockType == blockType_simply {
		dataForCf = b.data
	} else {

		dataForCf = []byte{}
		addresses := []addresInTOC{}
		for _, sb := range b.subBlocks {

			headerAndAttrs := []byte{}
			subBlockAttrs := sb.Attrs.getData()
			utils.AddToSlice(&headerAndAttrs,
				getHeaderForCf(len(subBlockAttrs), len(subBlockAttrs)),
				subBlockAttrs)

			headerAndData := []byte{}
			utils.AddToSlice(&headerAndData,
				getHeaderForCf(len(sb.Data), len(sb.Data)),
				sb.Data)

			pos := utils.AddToSlice(&dataForCf, headerAndAttrs, headerAndData)

			addresses = append(addresses,
				addresInTOC{AttrsPos: int32(pos[0]), DataPos: int32(pos[1])})
		}

		addTableOfContent(&dataForCf, addresses)

	}

	if len(dataForCf) > 0 {
		utils.ZlibCompress(dataForCf, &dataForCf)
		if len(dataForCf) == 0 {
			panic("Ошибка сжатия данных блока: " + b.attrs.name)
		}

		prepareDataForConfigFile(&dataForCf)
	}

	attrsForCf = b.attrs.getData()

	return
}

// Записать данные блока в файлы
func (b block) WriteToFile(pathToDir string) {

	var pathToBlockDir string

	if b.blockType == blockType_simply {
		pathToBlockDir = pathToDir
	} else if b.blockType == blockType_multiple {
		dirName := path.Join(pathToDir, b.GetName())
		err := os.Mkdir(dirName, os.ModeDir)
		if err != nil {
			panic(fmt.Sprintf("Ошибка создания каталога для записи блока в файл '%s': %s", b.GetName(), err.Error()))
		}
		pathToBlockDir = dirName
	} else {
		panic(fmt.Sprintf("Ошибка записи блока в файл '%s': тип не определен", b.GetName()))
	}

	subBlocksList := b.GetData()
	for _, sb := range subBlocksList {
		filename := path.Join(pathToBlockDir, sb.Attrs.name)
		err := ioutil.WriteFile(filename, sb.Data, os.ModeAppend)

		if err != nil {
			panic(fmt.Sprintf("Ошибка записи данных блока в файл '%s': %s", b.GetName(), err.Error()))
		}
	}
}

// Получить наименование блока
func (b block) GetName() string {
	return b.attrs.name
}

// Блок содержит данные формы
func (b block) IsForm() bool {
	return b.checkTypeBlock("form")
}

// Блок содержит данные модуля
func (b block) IsModule() bool {
	return b.checkTypeBlock("text")
}

// Проверить тип блока
func (b block) checkTypeBlock(dataType string) bool {

	for _, sb := range b.subBlocks {

		if sb.Attrs.name == dataType {
			return true
		}
	}

	return false
}

// Получить данные атрибутов и данных блока на основании заголовков
// @param данные файла
// @param заголовок атрибутов блока
// @param заголовок данных блока
// @param атрибуты блока
// @param данные блока
func getAttrsAndData(sourceData []byte, attrsHeader header, dataHeader header, attrs *attributes, data *[]byte) {

	attrsData := getRegion(sourceData, attrsHeader)
	*attrs = *createAttrs(attrsData)

	blockData := getRegion(sourceData, dataHeader)
	*data = make([]byte, len(blockData))
	copy(*data, blockData)
}
