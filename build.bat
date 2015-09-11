@echo off
rem -*- coding:OEM -*-
Setlocal EnableDelayedExpansion


SET "CURRENT_DIR=%CD%"
SET "GOPATH=%CURRENT_DIR%"
SET "GOBIN=%GOPATH%\bin"

cd src

if "%1" == "zlib" (
    cd zlibwrapper
    call build.bat
    cd ..
)

for /F "tokens=*" %%f in ('go env') DO (
    SET "_source=%%f"
    SET "_result=!_source:set =!"
    rem echo !_result!
    SET "!_result!"
)

go env

rem work (Dynamic link. it's not good)
echo v1. build %cd% with CGO_ENABLED=1 go install ...
SET "CGO_ENABLED=1"
go install

rem don't work (error: utils\compressor.go:4:2: no buildable Go source files in .....\ConfRobber\src\zlibwrapper)
echo v2. build %cd% with CGO_ENABLED=0 go install ...
SET "CGO_ENABLED=0"
go install

rem don't work (error: not link libzlibstatic.a)
echo v3. build %cd% with CGO_ENABLED=1 go build ...
SET "CGO_ENABLED=1"
go build --ldflags "-linkmode internal" -a -v -x