package ib

import (
	"errors"
	"fmt"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"time"
	"utils"
)

/**
* Константы запуска в пакетном режиме
 */
const (
	MODE_DESIGNER   = "DESIGNER"       // Режим конфигуратора
	MODE_CREATE     = "CREATEINFOBASE" // Режим создания новой ИБ
	MODE_ENTERPRISE = "ENTERPRISE"     // Режим предприятия
)

// Объект для работы с информационной базой 1С v8
type InformationBase struct {
	pathTo1C           string               // путь к файлу 1cv8.exe
	connectionSettings IBConnectionSettings // параметры подключения к информационной базе
}

func CreateInformationBase(pathTo1C string, connectionSettings IBConnectionSettings) (ib *InformationBase) {
	return &InformationBase{pathTo1C: pathTo1C, connectionSettings: connectionSettings}
}

/**
* Загрузить конфигурацию из cf файла
*
* @param полный путь к cf файлу
*
* @result код результата и информация об ошибке (если есть)
 */
func (ib InformationBase) LoadCfg(pathToFile string) (resultCode int, err error) {

	utils.AddTextToLog(utils.LogLevel_INFO, "Загрузка cf:"+pathToFile)

	exist, err := utils.Exists(pathToFile)
	if err != nil {
		utils.AddTextToLog(utils.LogLevel_ERROR, err.Error())
		return
	} else if !exist {
		err = errors.New("Файл не найден: " + pathToFile)
		utils.AddTextToLog(utils.LogLevel_ERROR, err.Error())
		return
	}

	args := []string{"/LoadCfg", pathToFile}
	resultCode, err = ib.runApplication(MODE_DESIGNER, args)

	utils.AddTextToLog(utils.LogLevel_INFO, fmt.Sprintf("-Загрузка cf: %t", err == nil))
	return
}

/**
* Выполняет обновление конфигурации информационной базы 1С
*
* @param код результата, который вернуло приложение
*
* @result код результата и информация об ошибке (если есть)
 */
func (ib InformationBase) UpdateDBCfg() (resultCode int, err error) {

	utils.AddTextToLog(utils.LogLevel_INFO, "Обновление ИБ")

	args := []string{"/UpdateDBCfg"}
	resultCode, err = ib.runApplication(MODE_DESIGNER, args)

	utils.AddTextToLog(utils.LogLevel_INFO, fmt.Sprintf("-Обновление ИБ: %t", err == nil))
	return
}

/**
* Сохранить конфигурацию в файл
*
* @param код результата, который вернуло приложение
* @param путь к файлу результата
*
* @result код результата и информация об ошибке (если есть)
 */
func (ib InformationBase) DumpCfg(pathToFile string) (resultCode int, err error) {

	utils.AddTextToLog(utils.LogLevel_INFO, "Сохранение cf: "+pathToFile)

	err = utils.RemoveIfExist(pathToFile)
	if err != nil {
		utils.AddTextToLog(utils.LogLevel_ERROR, "Ошибка удаления файла: "+pathToFile+":"+err.Error())
		return
	}

	args := []string{"/DumpCfg", pathToFile}
	resultCode, err = ib.runApplication(MODE_DESIGNER, args)

	utils.AddTextToLog(utils.LogLevel_INFO, fmt.Sprintf("-Сохранение cf: %t", err == nil))
	return
}

/**
* Выгрузка свойств объектов метаданных конфигурации в XML-файлы
*
* @param код результата, который вернуло приложение
* @param каталог, в который будет выгружена конфигурация
*
* @result код результата и информация об ошибке (если есть)
 */
func (ib InformationBase) DumpConfigToFiles(pathToDirectory string) (resultCode int, err error) {

	utils.AddTextToLog(utils.LogLevel_INFO, "Выгрузка конфигурации в xml: "+pathToDirectory)

	err = utils.RemoveIfExist(pathToDirectory)
	if err != nil {
		utils.AddTextToLog(utils.LogLevel_ERROR, "Ошибка удаления каталога: "+pathToDirectory+":"+err.Error())
		return
	}

	args := []string{"/DumpConfigToFiles", pathToDirectory}
	resultCode, err = ib.runApplication(MODE_DESIGNER, args)

	utils.AddTextToLog(utils.LogLevel_INFO, fmt.Sprintf("-Выгрузка конфигурации в xml: %t", err == nil))
	return
}

/**
* Загрузка свойств объектов метаданных конфигурации
*
* @param код результата, который вернуло приложение
* @param каталог, содержащий XML-файлы конфигурации
*
* @result код результата и информация об ошибке (если есть)
 */
func (ib InformationBase) LoadConfigFromFiles(pathToDirectory string) (resultCode int, err error) {

	utils.AddTextToLog(utils.LogLevel_INFO, "Загрузка конфигурации из xml: "+pathToDirectory)

	exist, err := utils.Exists(pathToDirectory)
	if err != nil {
		utils.AddTextToLog(utils.LogLevel_ERROR, err.Error())
		return
	} else if !exist {
		err = errors.New("Каталог не найден: " + pathToDirectory)
		utils.AddTextToLog(utils.LogLevel_ERROR, err.Error())
		return
	}

	args := []string{"/LoadConfigFromFiles", pathToDirectory}
	resultCode, err = ib.runApplication(MODE_DESIGNER, args)

	utils.AddTextToLog(utils.LogLevel_INFO, fmt.Sprintf("-Загрузка конфигурации из xml: %t", err == nil))
	return
}

/**
* Выгрузка свойств объектов метаданных конфигурации (модули и шаблоны)
*
* @param каталог расположения файлов свойств
* @param типы которые выгружаем:
*                   Module — признак необходимости выгрузки модулей;
                    Template — признак необходимости выгрузки макетов;
                    Help — признак необходимости выгрузки справочной информации;
                    AllWritable — признак выгрузки свойств только доступных для записи объектов;
                    Picture — признак выгрузки общих картинок;
                    Right — признак выгрузки прав.
*
* @result код результата и информация об ошибке (если есть)
*/
func (ib InformationBase) DumpConfigFiles(pathToDirectory string, types ...string) (resultCode int, err error) {

	utils.AddTextToLog(utils.LogLevel_INFO, "Выгрузка конфигурации: "+pathToDirectory)

	err = utils.RemoveIfExist(pathToDirectory)
	if err != nil {
		utils.AddTextToLog(utils.LogLevel_ERROR, "Ошибка удаления каталога: "+pathToDirectory+":"+err.Error())
		return
	}

	// /DumpConfigFiles "{0}" -Module -Template
	args := []string{"/DumpConfigFiles", pathToDirectory}
	for _, v := range types {
		args = append(args, "-"+v)
	}
	resultCode, err = ib.runApplication(MODE_DESIGNER, args)

	utils.AddTextToLog(utils.LogLevel_INFO, fmt.Sprintf("-Выгрузка конфигурации: %t", err == nil))
	return
}

/**
* Создать новую информационную базу
*
* @param код результата, который вернуло приложение
* @param шаблон конфигурации на основании которого необходимо создать информационную базу (может быть не указан)
*
* @result код результата и информация об ошибке (если есть)
 */
func (ib InformationBase) Create(pathToTemplate string) (resultCode int, err error) {

	utils.AddTextToLog(utils.LogLevel_INFO, "Создание новой ИБ: "+pathToTemplate)

	if len(pathToTemplate) > 0 {
		var exist bool
		exist, err = utils.Exists(pathToTemplate)
		if err != nil {
			utils.AddTextToLog(utils.LogLevel_ERROR, err.Error())
			return
		} else if !exist {
			err = errors.New("Файл не найден: " + pathToTemplate)
			utils.AddTextToLog(utils.LogLevel_ERROR, err.Error())
			return
		}
	}

	err = utils.RemoveIfExist(ib.connectionSettings.GetIbPath())
	if err != nil {
		utils.AddTextToLog(utils.LogLevel_ERROR, "Ошибка удаления каталога: "+ib.connectionSettings.GetIbPath()+":"+err.Error())
		return
	}

	/*
		1С предлагает для создания новой конфигурации на основе существующего
		cf файла использование параметра /UseTemplate, но по какой-то причине
		на некоторых версия платформы он работает не корректно.

		Поэтому делаем в два этапа:
		1. создаем новую конфигурацию
		2. загружаем в нее нужный cf файл
		3. обновляем конфигурацию информационной базы
	*/

	args := []string{}
	resultCode, err = ib.runApplication(MODE_CREATE, args)

	if err == nil && len(pathToTemplate) > 0 {
		resultCode, err = ib.LoadCfg(pathToTemplate)
	}

	if err == nil {
		resultCode, err = ib.UpdateDBCfg()
	}

	utils.AddTextToLog(utils.LogLevel_INFO, fmt.Sprintf("-Создание новой ИБ: %t", err == nil))
	return
}

/**
* Запустить приложение в пакетном режиме
*
* @param код результата, который вернуло приложение
* @param режим открытия информационной базы. Например: конфигуратор или предприятие
* @param аргументы запуска (может быть пустым)
*
* @result код результата и информация об ошибке (если есть)
 */
func (ib InformationBase) runApplication(mode string, args []string) (resultCode int, err error) {

	fullargs := []string{mode}
	ib.addAccessString(mode, &fullargs)
	for _, v := range args {

		if len(v) > 0 {
			fullargs = append(fullargs, v)
		}
	}
	addPathToLog(mode, &fullargs)

	cmd := exec.Command(ib.pathTo1C, fullargs...)
	utils.AddTextToLog(utils.LogLevel_INFO, "Командная строка: "+strings.Join(cmd.Args, " "))

	if err = cmd.Run(); err == nil {
		waitStatus := cmd.ProcessState.Sys().(syscall.WaitStatus)
		resultCode = waitStatus.ExitStatus()
	} else {
		utils.AddTextToLog(utils.LogLevel_ERROR, err.Error())

		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus := exitError.Sys().(syscall.WaitStatus)
			resultCode = waitStatus.ExitStatus()
		}
	}

	utils.AddTextToLog(utils.LogLevel_INFO, fmt.Sprintf("Код результата: %d", resultCode))
	return
}

/**
* Добавляет в параметры вызова приложения 1cv8.exe аргументы с параметрами доступа к информационной базе в пакетном режиме
*
* @param режим для которого необходио получить параметры доступа к ИБ
* @param аргументы вызова приложения 1cv8.exe
 */
func (ib InformationBase) addAccessString(mode string, args *[]string) {

	if mode == MODE_DESIGNER || mode == MODE_ENTERPRISE {
		// Для серверного режима вместо '/F' нужно указывать '/S'
		*args = append(*args, "/F", ib.connectionSettings.GetIbPath())

		if len(ib.connectionSettings.GetIbUserName()) > 0 {
			*args = append(*args, "/N", ib.connectionSettings.GetIbUserName())
		}

		if len(ib.connectionSettings.GetIbUserPwd()) > 0 {
			*args = append(*args, "/P", ib.connectionSettings.GetIbUserPwd())
		}

	} else if mode == MODE_CREATE {
		*args = append(*args, fmt.Sprintf("File='%s'", ib.connectionSettings.GetIbPath()))
	}

	if mode == MODE_DESIGNER && len(ib.connectionSettings.GetStoragePath()) > 0 {
		*args = append(*args, "/ConfigurationRepositoryF", ib.connectionSettings.GetStoragePath())
		*args = append(*args, "/ConfigurationRepositoryN", ib.connectionSettings.GetStorageUserName())

		if len(ib.connectionSettings.GetStorageUserPwd()) > 0 {
			*args = append(*args, "/ConfigurationRepositoryP", ib.connectionSettings.GetStorageUserPwd())
		}
	}
}

/**
* Добавляет путь к файлу логирования в аргументы вызовы
*
* @param режим, в котором осуществляется подключение к ИБ
* @param аргументы вызова приложения 1cv8.exe
 */
func addPathToLog(mode string, args *[]string) {

	var operationName string
	if mode == MODE_DESIGNER {
		for _, v := range *args {

			if v[0] == '/' && len(v) > 2 { // кроме праметров: "/F", "/P", "/N" и т.п.
				operationName = strings.Replace(v, "/", "", 1) // "/DumpCfg" => "DumpCfg"
				break
			}
		}
	} else if mode == MODE_CREATE {
		operationName = "Create"
	} else if mode == MODE_ENTERPRISE {
		operationName = "enterprise"
	}

	t := time.Now()
	timeStamp := fmt.Sprintf("%d-%02d-%02d_%02d-%02d-%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	path := path.Join(utils.GetPathToLogDir(), fmt.Sprintf("_%s_%s.log", timeStamp, operationName))
	*args = append(*args, "/Out", path)
}
