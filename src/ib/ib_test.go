package ib

import (
	"flag"
	"os"
	"path"
	"testing"
	"utils"
)

var pathToTarget string

func TestMain(m *testing.M) {

	currentDir, _ := utils.GetPathToCurrentDir()

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

func TestInformationBase(t *testing.T) {

	pathTo1C := path.Join("C:", "Program Files (x86)", "1cv8", "8.3.5.1517", "bin", "1cv8.exe")
	pathToIb := path.Join(pathToTarget, "ib")
	userName := ""
	userPwd := ""

	cn := CreateIBConnectionSettings(pathToIb, userName, userPwd, "", "", "")
	ib := CreateInformationBase(pathTo1C, *cn)
	_, err := ib.Create("")
	if err != nil {
		t.Error("InformationBase Create", err.Error())
	}
}
