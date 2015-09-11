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

rem работает, но линковка динамическая
echo build %cd% with CGO_ENABLED=1...
SET "CGO_ENABLED=1"
go install

rem не работает
echo build %cd% with CGO_ENABLED=0...
SET "CGO_ENABLED=0"
go install