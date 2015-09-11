package utils

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

const (
	LogLevel_TRACE   = 1
	LogLevel_INFO    = 2
	LogLevel_WARNING = 3
	LogLevel_ERROR   = 4
)

var (
	LogInfo      *log.Logger // Объект с помощью, которого будет изменяться файл логирования
	logLevel     *int        // Уровень логирования. Если сообщение ниже текущего уровня, то оно пропускается
	logPathToDir *string     // путь к каталогу в котором расположен лог
)

func InitLogger(pathToDir string, level int) {

	logLevel = new(int)
	*logLevel = level

	logPathToDir = new(string)
	*logPathToDir = createLogDir(pathToDir)

	pathToLog := createLogFile()

	file, err := os.OpenFile(pathToLog, log.Ldate|log.LUTC|log.Lmicroseconds, 0666)
	checkError(fmt.Sprintf("Ошибка открытия файла логирования '%s'", pathToLog), err)

	multi := io.MultiWriter(file, os.Stdout)
	LogInfo = log.New(multi, "", log.Ldate|log.LUTC|log.Lmicroseconds)
}

func LogIsInit() bool {
	return logPathToDir != nil
}

func GetPathToLogDir() string {
	return *logPathToDir
}

func GetLogLevel() int {
	return *logLevel
}

func AddTextToLog(level int, text string) {

	// Пропускаем всё, что меньше текущего уровня логирования
	if level < *logLevel {
		return
	}

	var prefixText string
	switch level {
	case LogLevel_WARNING:
		prefixText = "INFO"
	case LogLevel_ERROR:
		prefixText = "ERROR"
	case LogLevel_TRACE:
		prefixText = "TRACE"
	default:
		prefixText = "INFO"
	}

	LogInfo.Printf("%s: %s", prefixText, text)
}

func AddErrorToLog(err error) {
	AddTextToLog(LogLevel_ERROR, err.Error())
}

func checkError(prefix string, err error) {

	if err != nil {
		log.Fatalln(prefix + ":" + err.Error())
	}
}

func createLogDir(pathToDir string) (pathToLogDir string) {

	var err error
	var exist bool

	if exist, err = Exists(pathToDir); !exist && err == nil {
		err = errors.New("Не найден")
	}
	checkError(fmt.Sprintf("Ошибка проверки наличия каталога '%s'", pathToDir), err)

	pathToLogDir = path.Join(pathToDir, "Logs")
	exist, err = Exists(pathToLogDir)
	checkError(fmt.Sprintf("Ошибка проверки наличия каталога для хранения логов '%s'", pathToLogDir), err)

	if !exist {
		err = os.Mkdir(pathToLogDir, os.ModeDir)
		checkError(fmt.Sprintf("Ошибка проверки файла '%s'", pathToLogDir), err)
	}

	return
}

func createLogFile() (pathToLog string) {

	var err error
	var exist bool

	pathToLog = path.Join(GetPathToLogDir(), "main.log")
	exist, err = Exists(pathToLog)
	checkError(fmt.Sprintf("Ошибка проверки файла '%s'", pathToLog), err)

	if exist {
		err = os.Remove(pathToLog)
		checkError(fmt.Sprintf("Ошибка удаления файла '%s'", pathToLog), err)
	}

	file, err := os.Create(pathToLog)
	checkError(fmt.Sprintf("Ошибка создания файла '%s'", pathToLog), err)
	defer file.Close()

	return
}
