
#ifndef __zlibwrapper_h__
#define __zlibwrapper_h__

#include <_mingw.h>
#include <stdbool.h>
#include <vector>

typedef unsigned char BYTE;
typedef std::vector<BYTE> BinaryData;


/**
 * ����������� ������ ���������� �� cf ����� � ���������
 * @result - true - ���� ���������� ������ ��������� ��� ������
 */
BinaryData uncompress(const BinaryData &sourceData);

/**
 * ��������� ������ � ������ ��������� ��� cf �����
 * @result - true - ���� �������� ������ ��������� ��� ������
 */
BinaryData compress(const BinaryData &sourceData);

#endif // __zlibwrapper_h__
