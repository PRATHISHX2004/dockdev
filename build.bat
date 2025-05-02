@echo off
echo Initializing Go Linux build...

REM
set GOOS=linux
set GOARCH=amd64

REM
go build -o dockdev ./cmd

IF EXIST dockdev (
    echo Build successful: dockdev created.
) ELSE (
    echo Build failed. No output binary found.
)