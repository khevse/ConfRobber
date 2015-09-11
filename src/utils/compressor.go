package utils

import (
	"zlibwrapper"
)

// Выполняет сжатие данных с помощью библиотеки zlib
func ZlibCompress(sourceData []byte, resultData *[]byte) {

	var vectorSourceData zlibwrapper.BinaryData
	convertBytesToVector(sourceData, &vectorSourceData)
	vectorResultData := zlibwrapper.Compress(vectorSourceData)

	convertVectorToBytes(vectorResultData, resultData)

	zlibwrapper.DeleteBinaryData(vectorSourceData)
	zlibwrapper.DeleteBinaryData(vectorResultData)
}

// Выполняет распоковку данных с помощью библиотеки zlib
func ZlibUncompress(sourceData []byte, resultData *[]byte) {

	var vectorSourceData zlibwrapper.BinaryData
	convertBytesToVector(sourceData, &vectorSourceData)
	vectorResultData := zlibwrapper.Uncompress(vectorSourceData)

	convertVectorToBytes(vectorResultData, resultData)

	zlibwrapper.DeleteBinaryData(vectorSourceData)
	zlibwrapper.DeleteBinaryData(vectorResultData)
}

// Конвертация колекции байт в вектор
func convertBytesToVector(sourceData []byte, vector *zlibwrapper.BinaryData) {

	*vector = zlibwrapper.NewBinaryData()
	(*vector).Reserve(int64(len(sourceData)))
	for _, v := range sourceData {
		(*vector).Add(v)
	}
}

// Конвертация вектора в набор байт
func convertVectorToBytes(vector zlibwrapper.BinaryData, resultData *[]byte) {

	vectorSize := vector.Size()
	*resultData = make([]byte, vectorSize)

	for i := int64(0); i < vectorSize; i++ {
		(*resultData)[i] = vector.Get(int(i))
	}
}
