#include "zlibwrapper.hpp"
#include "zlib.h"


/**
 * Распаковать данные полученные из cf файла в контейнер
 *
 * @param - исходные данные, которые будем распаковывать
 * @param - размер исходных данных
 *
 * @result - true - если распаковка данных выполнена без ошибок то результат будет не пустой
 */
BinaryData uncompress(const BinaryData &sourceData) {

	BinaryData processedData;

	z_stream strm;
	strm.zalloc = Z_NULL;
	strm.zfree = Z_NULL;
	strm.opaque = Z_NULL;
	strm.avail_in = 0;
	strm.next_in = Z_NULL;
	int ret = inflateInit2(&strm, -MAX_WBITS);
	if ( ret != Z_OK ) {
		return processedData;
	}

	BYTE *sourceDataNotConst = const_cast<BYTE*>( sourceData.data() );
	strm.avail_in = sourceData.size();
	strm.next_in = sourceDataNotConst;

	const int bufferLen = 5120;
	BYTE buffer[bufferLen] = { 0 };

	do {
		strm.avail_out = bufferLen;
		strm.next_out = buffer;
		ret = inflate(&strm, Z_NO_FLUSH);
		if ( ret == Z_STREAM_ERROR) {
			processedData.clear();
			break;
		}

		const size_t inflateLen = bufferLen - strm.avail_out;
		const size_t currentLen = processedData.size();

		processedData.resize(currentLen + inflateLen);
		memcpy(processedData.data() + currentLen, buffer, inflateLen);

		memset(buffer, 0, bufferLen);

	} while ( strm.avail_out == 0 );

	(void) inflateEnd(&strm);


	return processedData;
}

/**
 * Упаковать данные в формат пригодный для cf файла
 *
 * @param - исходные данные, которые подлежат сжатию
 * @param - размер исходных данных
 *
 * @result - true - если упаковка данных выполнена без ошибок
 */
BinaryData compress(const BinaryData &sourceData) {

	BinaryData processedData;

	z_stream strm;
	strm.zalloc = Z_NULL;
	strm.zfree  = Z_NULL;
	strm.opaque = Z_NULL;

	int ret = deflateInit2(&strm, Z_BEST_COMPRESSION, Z_DEFLATED, -MAX_WBITS, 8, Z_DEFAULT_STRATEGY);

	if (ret != Z_OK) {
		return processedData;
	}

	const int bufferLen = 16384;
	BYTE buffer[bufferLen] = { 0 };

	BYTE *sourceDataNotConst = const_cast<BYTE*>( sourceData.data() );
	size_t readDataSize    = 0;
	size_t balanceDataSize = sourceData.size();

	while ( balanceDataSize > 0 ) {
		const size_t currentDataSize = balanceDataSize > bufferLen ? bufferLen : balanceDataSize;

		strm.avail_in = currentDataSize;
		strm.next_in  = sourceDataNotConst + readDataSize;

		readDataSize += currentDataSize;
		balanceDataSize -= currentDataSize;

		do {
			strm.avail_out = bufferLen;
			strm.next_out  = buffer;

			ret = deflate(&strm, balanceDataSize == 0 ? Z_FINISH : Z_NO_FLUSH);

			const size_t deflateLen = bufferLen - strm.avail_out;
			const size_t currentLen = processedData.size();

			processedData.resize(currentLen + deflateLen);
			memcpy(processedData.data() + currentLen, buffer, deflateLen);
		} while (strm.avail_out == 0 && ret != Z_STREAM_END);
	}

	(void)deflateEnd(&strm);

	return processedData;
}


