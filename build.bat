@echo off
rem -*- coding:OEM -*-
Setlocal EnableDelayedExpansion

SET "CURRENT_DIR=%CD%"
SET "GOPATH=%CURRENT_DIR%"
SET "GOBIN=%GOPATH%/bin"

go clean

cmd /c "cd /d ""%CURRENT_DIR%/src/zlibwrapper"" && build"

go build -a -v -x

if not %errorlevel% == 0 (
    echo exit code "%errorlevel%"
    exit %errorlevel%
)