Setlocal EnableDelayedExpansion
@echo off
rem -*- coding:OEM -*-

cd ..
cd ..
SET "GOPATH=%CD%"
cd src
cd zlibwrapper

SET "CURRENT_DIR=%CD%"
SET "BUILD_DIR_NAME=build"
SET "BUILD_DIR=%CURRENT_DIR%\build"

rem 32 или 64
SET "SYS_TYPE=32"

SET "MINGW_DIR=C:/C++/MinGW/mingw64"
SET "PATH=%PATH%;%GOBIN%;%MINGW_DIR%"


rem Удаление не нужных файлов
for %%i in (bin, %BUILD_DIR_NAME%, include, lib, share) do (
    rmdir /s /q "%CURRENT_DIR%\%%i"
)

for %%i in (*.exe, *.c, *.cxx, *.a, *.dll, *.6, *.log, *.go) do (
    del /f /s /q "%CURRENT_DIR%\%%i"
)

if "%1" == "clean" (
    exit
)

 
SWIG -go -intgosize %SYS_TYPE%  -c++ -cgo zlibwrapper.i

rem echo Add to file zlibwrapper.go new string
call replace_text "zlibwrapper.go" "#define intgo swig_intgo" "#cgo LDFLAGS: -L. -lzlibstatic -lzlibwrapper"

mkdir %BUILD_DIR_NAME%
cd %BUILD_DIR_NAME%

cmake -G "MinGW Makefiles" -DCMAKE_BUILD_TYPE=RELEASE -DBUILD_SHARED_LIBS=0 ../
mingw32-make install
if "%ERRORLEVEL%" == "0" (
    cd ..
    call go_install1.5
) Else (
    cd ..
    exit
)

if not "%1" == "release" (
    if not "%1" == "" (
        exit
    )
)

rem Удаление не нужных файлов
for %%i in (bin, %BUILD_DIR_NAME%, include, lib, share) do (
    rmdir /s /q "%CURRENT_DIR%\%%i"
)

for %%i in (*.exe, libzlib.dll, libzlib.dll.a) do (
    del /f /s /q "%CURRENT_DIR%\%%i"
)

rem pause