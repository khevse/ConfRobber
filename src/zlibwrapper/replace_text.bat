@echo off
rem -*- coding:OEM -*-
Setlocal EnableDelayedExpansion

SET "PARAM1=%1"
SET "PARAM2=%2"
SET "PARAM3=%3"
SET "PARAM4=%4"

rem Delete quotes
SET "SRC=%PARAM1:~1,-1%"
SET "FIND_STR=%PARAM2:~1,-1%"
SET "NEW_STR=%PARAM3:~1,-1%"

SET "NEW_STR=%NEW_STR:^lt;=^<%"
SET "NEW_STR=%NEW_STR:^qt;=^>%"

rem echo param1=%PARAM1%; param2=%PARAM2%; param3=%PARAM3%
echo param1=%SRC%; param2=%FIND_STR%; param3=%NEW_STR%

for /f "usebackq tokens=*" %%a in ("%SRC%") do (
    if "%%a"=="" (
        echo.>>tmp
    ) else (
        if "%%a"=="%FIND_STR%" (
            echo %NEW_STR%>>tmp
            if "%PARAM4%" == "" (
                echo %FIND_STR%>>tmp
            )
        ) else (
            echo %%a>>tmp
        )
    )
)

copy /y tmp %SRC% 1>nul&&del /f /q tmp>nul

go fmt