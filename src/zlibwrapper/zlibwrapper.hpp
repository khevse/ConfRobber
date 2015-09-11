
#ifndef __zlibwrapper_h__
#define __zlibwrapper_h__

#include <_mingw.h>
#include <stdbool.h>
#include <vector>

typedef unsigned char BYTE;
typedef std::vector<BYTE> BinaryData;


/**
 * –аспаковать данные полученные из cf файла в контейнер
 * @result - true - если распаковка данных выполнена без ошибок
 */
BinaryData uncompress(const BinaryData &sourceData);

/**
 * ”паковать данные в формат пригодный дл€ cf файла
 * @result - true - если упаковка данных выполнена без ошибок
 */
BinaryData compress(const BinaryData &sourceData);

#endif // __zlibwrapper_h__
