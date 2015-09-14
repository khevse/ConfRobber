@echo off
rem -*- coding:OEM -*-
Setlocal EnableDelayedExpansion

cd ..
cd ..
SET "GOPATH=%CD%"
cd src
cd zlibwrapper

SET "CURRENT_DIR=%CD%"
SET "BUILD_DIR_NAME=build"
SET "BUILD_DIR=%CURRENT_DIR%\build"

rem 32 or 64
SET "SYS_TYPE=64"

if "%MINGW_DIR%" == "" (
    SET "MINGW_DIR=C:/C++/MinGW/mingw64"
    SET "PATH=%PATH%;%GOBIN%;%MINGW_DIR%"
)

rem Remove old files
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

for /F "tokens=*" %%f in ('go env') DO (
    SET "_source=%%f"
    SET "_result=!_source:set =!"
    rem echo !_result!
    SET "!_result!"
)

if "%PATH_TO_PROJECT_LIBS%" == "" (
    SET "PATH_TO_PROJECT_LIBS=%GOPATH%\pkg\%GOOS%_%GOARCH%"
)

SET "PATH=%PATH%;%PATH_TO_PROJECT_LIBS%"
rem echo Add to file zlibwrapper.go new string
call replace_text "zlibwrapper.go" "#define intgo swig_intgo" "#cgo windows LDFLAGS: -L. -L"%PATH_TO_PROJECT_LIBS%" -mwindows -lzlibwrapper -lzlib -lgcc -lstdc++"
rem call replace_text "zlibwrapper.go" "#define intgo swig_intgo" "#cgo windows CFLAGS: -fno-stack-check -fno-stack-protector -mno-stack-arg-probe"

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

rem Remove temp files
for %%i in (bin, %BUILD_DIR_NAME%, include, lib, share) do (
    rmdir /s /q "%CURRENT_DIR%\%%i"
)

for %%i in (*.exe, libzlib.dll.a) do (
    del /f /s /q "%CURRENT_DIR%\%%i"
)

rem pause