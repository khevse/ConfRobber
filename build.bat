@echo off
rem -*- coding:OEM -*-
Setlocal EnableDelayedExpansion


SET "CURRENT_DIR=%CD%"
SET "GOPATH=%CURRENT_DIR%"
SET "GOBIN=%GOPATH%\bin"
SET "MINGW_DIR=C:/C++/MinGW/mingw64"
SET "PATH=%PATH%;%GOBIN%;%MINGW_DIR%"

cd src

if "%1" == "zlib" (
    cd zlibwrapper
    call build_C_lib.bat
    cd ..
)

for /F "tokens=*" %%f in ('go env') DO (
    SET "_source=%%f"
    SET "_result=!_source:set =!"
    rem echo !_result!
    SET "!_result!"
)

go env

if "%1" == "v1" (
rem work (Dynamic link. it's not good)
echo v1. build %cd% with CGO_ENABLED=1 go install ...
SET "CGO_ENABLED=1"
go install
)

if "%1" == "v2" (
rem don't work (error: utils\compressor.go:4:2: no buildable Go source files in .....\ConfRobber\src\zlibwrapper)
echo v2. build %cd% with CGO_ENABLED=0 go install ...
SET "CGO_ENABLED=0"
go install
)

if "%1" == "" (
rem don't work (error: not link libzlibstatic.a)
echo v3. build %cd% with CGO_ENABLED=1 go build ...
SET "CGO_ENABLED=1"
SET "PATH_TO_PROJECT_LIBS=%GOPATH%\pkg\%GOOS%_%GOARCH%"
SET "GOGCCFLAGS=%GOGCCFLAGS% -L"%PATH_TO_PROJECT_LIBS%" -I"%PATH_TO_PROJECT_LIBS%" -lzlibstatic"

cd zlibwrapper
call build_C_lib.bat
cd ..

go build --ldflags "-linkmode internal" -a -v -x 
)
