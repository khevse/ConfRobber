@echo off
rem -*- coding:OEM -*-
Setlocal EnableDelayedExpansion

SET "CURRENT_DIR=%CD%"
SET "GOPATH=%CURRENT_DIR%"
SET "GOBIN=%GOPATH%\bin"
SET "MINGW_DIR=C:/C++/MinGW/mingw64"
SET "CGO_ENABLED=1"
SET "PATH_TO_PROJECT_LIBS=%GOPATH%\pkg\%GOOS%_%GOARCH%"
SET "GOGCCFLAGS=%GOGCCFLAGS% -L"%PATH_TO_PROJECT_LIBS%" -I"%PATH_TO_PROJECT_LIBS%""
SET "PATH=%PATH%;%GOBIN%;%MINGW_DIR%;%PATH_TO_PROJECT_LIBS%"

cd src

for /F "tokens=*" %%f in ('go env') DO (
    SET "_source=%%f"
    SET "_result=!_source:set =!"
    rem echo !_result!
    SET "!_result!"
)

go env

cd zlibwrapper
call build_C_lib.bat
cd ..
go build --ldflags "-linkmode internal" -a -v -x

cd ..

cd %GOBIN%
mkdir result

move ".\..\src\src.exe" "%GOBIN%\result\ConfRobber.exe"
copy ".\..\src\zlibwrapper\libzlib.dll" "%GOBIN%\result\libzlib.dll"
copy "%MINGW_DIR%\bin\libstdc++-6.dll" "%GOBIN%\result\libstdc++-6.dll"
copy "%MINGW_DIR%\bin\libgcc_s_seh-1.dll" "%GOBIN%\result\libgcc_s_seh-1.dll"
