@echo off

set APP_NAME=cli-app.exe
set GOOS=windows
set GOARCH=amd64
set OUTPUT_DIR=.\bin

if not exist %OUTPUT_DIR% (
    mkdir %OUTPUT_DIR%
)

echo Building: %GOOS%/%GOARCH%...
go build -o %OUTPUT_DIR%\%APP_NAME% -buildvcs=false .\cmd\cli

if exist "%OUTPUT_DIR%\%APP_NAME%" (
    echo Path: %OUTPUT_DIR%\%APP_NAME%
) else (
    echo Building error.
    exit /b 1
)