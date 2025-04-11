@echo off

set BIN_PATH=.\bin\cli-app.exe

if not exist %BIN_PATH% (
    echo File not found. Build app (scripts\build.bat).
    exit /b 1
)

echo Running
%BIN_PATH%