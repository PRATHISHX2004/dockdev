@echo off
echo Initializing Go Linux build...

REM Set GOOS and GOARCH
set GOOS=linux
set GOARCH=amd64

REM Create dist folder if it doesn't exist
if not exist "..\dist" (
    mkdir "..\dist"
)

REM Build the binary into ../dist
go build -o ..\dist\dockdev ./cmd

REM Check result
if exist ..\dist\dockdev (
    echo Build successful: ..\dist\dockdev
) else (
    echo Build failed! No output binary found.
    exit /b 1
)
