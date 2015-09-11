package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// Проверить наличие файла или каталога
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// Создает файл если если он отсутствует
func CreateIfNotExist(path string) (bool, error) {
	has, _ := Exists(path)
	if !has {
		file, err := os.Create(path)
		if err != nil {
			return false, err
		}
		defer file.Close()
	}

	return true, nil
}

// Возвращает каталог родитель
func GetParentDir(pathToFile string) string {

	parts := strings.Split(pathToFile, string(os.PathSeparator))

	newPath := ""
	for i, v := range parts {
		if i == len(parts)-1 {
			break
		}

		if len(newPath) == 0 {
			newPath = v
		} else {
			newPath += string(os.PathSeparator) + v
		}
	}

	return newPath
}

// Возвращает путь к исполняемому бинарному файлу
func GetPathToBinaryFile() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

// Возвращает путь к текущему каталогу
func GetPathToCurrentDir() (pathToDir string, err error) {

	pathToDir, err = os.Getwd()
	if err == nil {
		pathToDir = GetParentDir(pathToDir)
	} else {
		pathToDir = ""
	}

	return
}

func RemoveIfExist(pathToFile string) (err error) {

	exist, err := Exists(pathToFile)
	if err != nil {
		return
	}

	if exist {
		err = os.RemoveAll(pathToFile)
	}

	return
}

func ReadFilesInDir(pathToDir string) (fileInfos []os.FileInfo, err error) {

	dir, err := os.Open(pathToDir)
	if err != nil {
		return
	}
	defer dir.Close()

	fileInfos, err = dir.Readdir(-1)
	return
}
