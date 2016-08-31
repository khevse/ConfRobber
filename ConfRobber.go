package main

import (
	"cf"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"utils"
)

func main() {

	defer func() {
		str := recover()
		if str == nil {
			str = "Ok"
		} else if utils.LogIsInit() {
			utils.LogInfo.Println(str)
		}
		fmt.Println(str)
	}()

	var operationType string

	args := os.Args[1:]
	countParams := len(args)

	if countParams > 0 {
		operationType = args[0]
	}

	errorText := checkArgc(countParams, operationType)
	if len(errorText) > 0 {
		panic(errorText)
	}

	if operationType == "-P" {
		unpackToDir(args[1], args[2])
	} else if operationType == "-B" {
		buildCf(args[1], args[2])
	} else if operationType == "-help" {
		fmt.Println(checkArgc(countParams, operationType))
	}
}

// Распаковать конфигурационный файл в каталог
func unpackToDir(pathToCf string, pathToTarget string) {

	createTargetDir(pathToTarget)

	utils.InitLogger(pathToTarget, utils.LogLevel_INFO)
	utils.AddTextToLog(utils.LogLevel_INFO, "Начало")

	var err error
	var fileData []byte
	if fileData, err = ioutil.ReadFile(pathToCf); err != nil {
		panic(err.Error())
	}

	pathToUnpackDir := createTargetDir(path.Join(pathToTarget, "unpack"))

	var objectCf cf.ConfCf
	if err = objectCf.InitFromFile(fileData); err != nil {
		panic(err.Error())
	}

	objectCf.SaveBlocksToFiles(pathToUnpackDir)

	utils.AddTextToLog(utils.LogLevel_INFO, "-Завершение")
}

// Создать новый конфигурационный файл на основании данных файлов содержащихся в каталоге
func buildCf(pathToDirWithSourceData string, pathToCf string) {

	pathToTarget := createTargetDir(utils.GetParentDir(pathToCf))

	utils.InitLogger(pathToTarget, utils.LogLevel_INFO)
	utils.AddTextToLog(utils.LogLevel_INFO, "Начало")

	var err error
	var objectCf cf.ConfCf
	if err = objectCf.InitFromCatalog(pathToDirWithSourceData); err != nil {
		panic(err.Error())
	}

	fileData := objectCf.GetDataForConfigFile()
	if err = ioutil.WriteFile(pathToCf, fileData, os.ModeAppend); err != nil {
		panic(err.Error())
	}

	utils.AddTextToLog(utils.LogLevel_INFO, "-Завершение")
}

// Создать каталог в который будут записаны результаты
func createTargetDir(pathToTarget string) string {

	err := os.MkdirAll(pathToTarget, os.ModeDir)
	checkError(err)

	return pathToTarget
}

// Проверить наличие ошибок
func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// Проверить аргументы вызова
func checkArgc(countParams int, operationType string) (errorText string) {

	if operationType == "-P" {
		if countParams != 3 {
			lines := []string{"Operation type: unpack the configuration file (*.cf)",
				"Options:",
				" 1. Operation type - '-P'",
				" 2. Path to the source file *.cf",
				" 3. Path to the target directory"}

			errorText = strings.Join(lines, "\n")
		}
	} else if operationType == "-B" {
		if countParams != 3 {
			lines := []string{"Operation type: build the configuration file (*.cf)",
				"Options:",
				" 1. Operation type - '-B'",
				" 2. Path to the directory with the source files",
				" 3. Path to the configuration file (*.cf)"}

			errorText = strings.Join(lines, "\n")
		}
	} else if operationType == "-help" {
		errorText = fmt.Sprintf("%s\n\n%s\n",
			checkArgc(0, "-P"), checkArgc(0, "-B"))
	} else {
		errorText = fmt.Sprintf("Option is not valid. Available options: \n%s\n\n%s\n",
			checkArgc(0, "-P"), checkArgc(0, "-B"))
	}

	return
}
