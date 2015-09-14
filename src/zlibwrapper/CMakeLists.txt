cmake_minimum_required(VERSION 2.8.4)

set(PROJECT zlibwrapper)

project(${PROJECT})

set(PROJECT_SOURCE_DIR ${${PROJECT}_SOURCE_DIR})
set(${PROJECT}_BINARY_DIR "${PROJECT_SOURCE_DIR}")
set(CMAKE_INSTALL_PREFIX  "${PROJECT_SOURCE_DIR}")

# Каталог в который будут компилироваться исполняемые файлы проекта
set (EXECUTABLE_OUTPUT_PATH "${PROJECT_SOURCE_DIR}")
# Каталог в который будут компилироваться библиотеки проекта
set (LIBRARY_OUTPUT_PATH "${PROJECT_SOURCE_DIR}")

set (CMAKE_VERBOSE_MAKEFILE ON)
set(SYSTEM WIN$ENV{SYS_TYPE})

set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -std=c++11 -Wall -Wextra -fmessage-length=0 -D${SYSTEM} -fPIC -m64 -pthread -mwindows -static -static-libgcc -static-libstdc++")
set(CGO_ENABLED "0")

set(MINGW_DIR "$ENV{MINGW_DIR}")
set(CMAKE_LIBRARY_PATH ${CMAKE_LIBRARY_PATH} "${MINGW_DIR}/lib;${CMAKE_LIBRARY_PATH};${LIBRARY_OUTPUT_PATH}")
set(CMAKE_INCLUDE_PATH ${CMAKE_INCLUDE_PATH} "${CMAKE_INCLUDE_PATH};${MINGW_DIR}/include")

# Расположение заголовочных файлов
INCLUDE_DIRECTORIES(
    .
    zlib
    )

# Добавление подпроекта, указывается имя дирректории
add_subdirectory(zlib)

# Библиотеки, с которыми будем осуществлять линковку
set(LIBS zlib)

set(SOURCE_FILES 
	zlibwrapper_wrap.cxx
    zlibwrapper.cpp)

# Создание основной библиотеки файла
add_library(${PROJECT} STATIC ${SOURCE_FILES} zlibwrapper.hpp)

# Линковка программы с библиотеками
target_link_libraries(${PROJECT} ${LIBS})
