@ECHO OFF
SETLOCAL
chcp 866>nul
CLS

CD ..

SET BUILD_PATH=.\pomadorik

@REM SET GOPATH=%CD%
@REM SET PATH=%PATH%;%GOPATH%\BIN;
SET GO111MODULE=auto
SET GOMODCACHE=%CD%\packages


REM download all packages that project needs
go mod tidy

ECHO Building service...

@REM go run main.go

go build -o %BUILD_PATH%.exe

rem start in dev mode (default)
%BUILD_PATH%.exe