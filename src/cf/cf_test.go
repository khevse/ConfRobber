package cf

import (
	"bytes"
	"flag"
	"fmt"
	"ib"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"utils"
	"zlibwrapper"
)

const countBloks = 2

var (
	currentDir   string
	pathToTarget string
)

func getTestData_Attrs() attributes {

	attrs := attributes{modificationDate: 635668862670000,
		creationDate: 635668862670000 + 1,
		groupType:    groupTypeNoModule,
		name:         "name"}

	return attrs
}

func getTestData_SimpleBlockData() string {
	return "hello"
}

func getTestData_AttrsCF() []byte {
	return getTestData_Attrs().getData()
}

func getTestData_BlockCF() []byte {

	data, _ := zlibwrapper.Compress([]byte(getTestData_SimpleBlockData()))
	return data
}

func getTestData_HeaderAndAttrs() []byte {

	blockAttrs := getTestData_AttrsCF()

	data := []byte{}
	utils.AddToSlice(&data,
		getHeaderForCf(len(blockAttrs), len(blockAttrs)),
		blockAttrs)
	return data
}

func getTestData_HeaderAndBlockData() []byte {

	blockData := getTestData_BlockCF()

	data := []byte{}
	utils.AddToSlice(&data,
		getHeaderForCf(len(blockData), len(blockData)),
		blockData)

	return data
}

func getTestData_File() (fileData []byte, addresses []addresInTOC) {

	addresses = []addresInTOC{}
	fileData = []byte{}

	for i := 0; i < countBloks; i++ {
		pos := utils.AddToSlice(&fileData,
			getTestData_HeaderAndAttrs(),
			getTestData_HeaderAndBlockData())

		addresses = append(addresses,
			addresInTOC{AttrsPos: int32(pos[0]), DataPos: int32(pos[1])})
	}

	fileDataLenWithoutTOC := len(fileData)
	addTableOfContent(&fileData, addresses)
	TOCLen := len(fileData) - fileDataLenWithoutTOC

	for i, _ := range addresses {
		addresses[i].AttrsPos += int32(TOCLen)
		addresses[i].DataPos += int32(TOCLen)
	}

	return
}

func TestMain(m *testing.M) {

	currentDir, _ = utils.GetPathToCurrentDir()

	// В наименовании каталога специально присутствует пробел, чтобы протестировать
	// работу с передачей параметров в командной строке при открытии 1cv8.exe
	pathToTarget = path.Join(utils.GetParentDir(currentDir), "target dir")

	utils.RemoveIfExist(pathToTarget)
	os.Mkdir(pathToTarget, os.ModeDir)
	utils.InitLogger(pathToTarget, utils.LogLevel_TRACE)

	flag.Parse()
	os.Exit(m.Run())

	utils.RemoveIfExist(pathToTarget)

}

func TestParseHeader(t *testing.T) {

	getMsg := func(field string, standart int, val int) string {
		return fmt.Sprintf("Ошибка чтения '%s' %d != %d", field, standart, val)
	}

	headerAndAttrs := getTestData_HeaderAndAttrs()
	attrsCF := getTestData_AttrsCF()

	const headerPosition = 0

	h, err := createHeader(headerAndAttrs, headerPosition)
	if err != nil {
		t.Error("Ошибка разбора заголовка", err.Error())
	}

	if h.getDataSize() != len(attrsCF) {
		t.Error(getMsg("размер полезных данных блока", len(attrsCF), h.getDataSize()))
	}

	if h.totalRegionSize != len(attrsCF) {
		t.Error(getMsg("размер блока", len(attrsCF), h.totalRegionSize))
	}

	if h.getDataPosition() != headerPosition+fullHeaderSize {
		t.Error(getMsg("позиция блока в файле", headerPosition+fullHeaderSize, h.getDataPosition()))
	}
}

func TestIsHeader(t *testing.T) {

	const headerPosition = 0

	headerAndAttrs := getTestData_HeaderAndAttrs()

	if isHeader(headerAndAttrs, headerPosition+1) != false {
		t.Error("IsHeader: is not header", headerAndAttrs)
	}

	if isHeader(headerAndAttrs, headerPosition) != true {
		t.Error("IsHeader: is header", headerAndAttrs[headerPosition:])
	}
}

func TestIsAttributesData(t *testing.T) {

	const headerPosition = 0

	headerAndAttrs := getTestData_HeaderAndAttrs()
	h, err := createHeader(headerAndAttrs, headerPosition)
	if err != nil {
		t.Error("Ошибка определения принадлежности заголовка к данным атрибутов блока: ", err.Error())
	}

	if h.isAttributesData(headerAndAttrs) != true {
		t.Error("Ошибка определения принадлежности заголовка к данным атрибутов блока")
	}
}

func TestAttributesInit(t *testing.T) {

	const headerPosition = 0

	attrs := getTestData_Attrs()
	headerAndAttrs := getTestData_HeaderAndAttrs()

	h, err := createHeader(headerAndAttrs, headerPosition)
	if err != nil {
		t.Error("Ошибка чтения заголовка атрибутов блока: ", err.Error())
	}

	attrsData := getRegion(headerAndAttrs, *h)
	a := createAttrs(attrsData)

	if a.name != attrs.name {
		t.Error("Ошибка чтения наименования блока: ", attrs.name, " != ", a.name)
	}

	if a.groupType != attrs.groupType {
		t.Error("Ошибка чтения типа блока: ", attrs.groupType, " != ", a.groupType)
	}

	if a.creationDate != attrs.creationDate {
		t.Error("Ошибка чтения даты создания блока: ", attrs.creationDate, " != ", a.creationDate)
	}

	if a.modificationDate != attrs.modificationDate {
		t.Error("Ошибка чтения даты модификации блока: ", attrs.modificationDate, " != ", a.modificationDate)
	}
}

func TestReadFile(t *testing.T) {

	fileData, _ := getTestData_File()

	blocksHeadersPairs, err := readTOC(fileData)
	if err != nil {
		t.Error("Ошибка чтения адресов из оглавления: ", err.Error())
	} else if len(blocksHeadersPairs) != countBloks {
		t.Error("Ошибка чтения адресов из оглавления: ", countBloks, " != ", len(blocksHeadersPairs))
	}

	bloksList := make([]block, countBloks)

	standartAttrs := getTestData_Attrs()
	standartData := []byte(getTestData_SimpleBlockData())

	countProcessedBlocks := 0

	for i, p := range blocksHeadersPairs {
		countProcessedBlocks += 1

		bloksList[i] = block{}
		bloksList[i].Init(fileData, p.attrs, p.data, make(chan bool, 1))

		subBlocks := bloksList[i].GetData()
		if len(subBlocks) != 1 {
			t.Error("Ошибка получения подчиненных блоков", 1, " != ", len(subBlocks))
		}

		for _, sb := range subBlocks {

			if sb.Attrs.name != standartAttrs.name {
				t.Error("Ошибка определения имени подчиненного блока", standartAttrs.name, " != ", sb.Attrs.name)
			} else if !bytes.Equal(sb.Data, standartData) {
				t.Error("Ошибка определения данных подчиненного блока", sb.Data, " != ", standartData)
			}
		}
	}

	if countProcessedBlocks != countBloks {
		t.Error("Обработаны не все блоки", countBloks, " != ", countProcessedBlocks)
	}
}

func TestInitCfObject(t *testing.T) {

	fileData, _ := getTestData_File()

	var cf ConfCf
	err := cf.InitFromFile(fileData)

	if err != nil {
		t.Error("В cf файле ошибка:", err.Error())
	}

	standartAttrs := getTestData_Attrs()
	standartData := []byte(getTestData_SimpleBlockData())
	countProcessedBlocks := 0

	blocksList := cf.GetData()
	for _, b := range blocksList {
		countProcessedBlocks += 1

		subBlocks := b.GetData()
		if len(subBlocks) != 1 {
			t.Error("В cf файле ошибка: не найдено подчиненных блоков", 1, " != ", len(subBlocks))
		}

		for _, sb := range subBlocks {

			if sb.Attrs.name != standartAttrs.name {
				t.Error("В cf файле ошибка определения имени подчиненного блока", standartAttrs.name, " != ", sb.Attrs.name)
			} else if !bytes.Equal(sb.Data, standartData) {
				t.Error("В cf файле ошибка определения данных подчиненного блока", sb.Data, " != ", standartData)
			}
		}
	}

	if countProcessedBlocks != countBloks {
		t.Error("В cf файле найдены не все блоки", countBloks, " != ", countProcessedBlocks)
	}
}

func TestInit(t *testing.T) {

	var err error
	dirWithTestData := path.Join(utils.GetParentDir(currentDir), "test_data")
	pathToTestCf := path.Join(dirWithTestData, "original.cf")

	fileData, err := ioutil.ReadFile(pathToTestCf)
	if err != nil {
		t.Error(fmt.Sprintf("Ошибка чтения файла '%s': %s", pathToTestCf, err.Error()))
	}

	var objectCf ConfCf

	if err = objectCf.InitFromFile(fileData); err != nil {
		t.Error(err.Error())
	}

	pathToUnpackDir := path.Join(pathToTarget, "unpack")
	if err = os.Mkdir(pathToUnpackDir, os.ModeDir); err != nil {
		t.Error(err.Error())
	}

	objectCf.SaveBlocksToFiles(pathToUnpackDir)

	if err = objectCf.InitFromCatalog(pathToUnpackDir); err != nil {
		t.Error(err.Error())
	}

	newFileName := path.Join(pathToTarget, "assembly.cf")
	newFileData := objectCf.GetDataForConfigFile()

	if err = ioutil.WriteFile(newFileName, newFileData, os.ModeAppend); err != nil {
		t.Error(err.Error())
	}

	pathTo1C := path.Join("C:", "Program Files (x86)", "1cv8", "8.3.5.1517", "bin", "1cv8.exe")
	pathToIb := path.Join(pathToTarget, "ib")
	userName := ""
	userPwd := ""

	cn := ib.CreateIBConnectionSettings(pathToIb, userName, userPwd, "", "", "")
	ib := ib.CreateInformationBase(pathTo1C, *cn)
	_, err = ib.Create(newFileName)
	if err != nil {
		t.Error("InformationBase Create", err.Error())
	}

}
