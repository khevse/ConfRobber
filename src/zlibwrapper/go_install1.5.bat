@echo off
rem -*- coding:cp1251 -*-
Setlocal EnableDelayedExpansion

for /F "tokens=*" %%f in ('go env') DO (
    SET "_source=%%f"
    SET "_result=!_source:set =!"
    rem echo !_result!
    SET "!_result!"
)

go env

go install
if not "%ERRORLEVEL%" == "0" (
    echo exit code "%ERRORLEVEL%"
    exit
)

echo "Installing go shared lib..."

SET "PATH_TO_MINGW_LIBS=%MinGW%\x86_64-w64-mingw32\lib"
SET "PATH_TO_GOLANG_LIBS=%GOROOT%\pkg\%GOOS%_%GOARCH%"
SET "PATH_TO_PROJECT_LIBS=%GOPATH%\pkg\%GOOS%_%GOARCH%"

for %%i in (libzlibstatic.a, libzlibwrapper.a) do (
    
    echo copy_file ".\%%i" to "%PATH_TO_MINGW_LIBS%\%%i"
    copy ".\%%i" "%PATH_TO_MINGW_LIBS%\%%i"
    
    echo copy_file ".\%%i" to "%PATH_TO_GOLANG_LIBS%\%%i"
    copy ".\%%i" "%PATH_TO_GOLANG_LIBS%\%%i"

    echo copy_file ".\%%i" to "%PATH_TO_PROJECT_LIBS%\%%i"
    copy ".\%%i" "%PATH_TO_PROJECT_LIBS%\%%i"
)