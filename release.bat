@ECHO OFF
SETLOCAL
chcp 866>nul
CLS


SET BUILD_PATH=.\pomadorik_release
SET GO111MODULE=auto
SET GOMODCACHE=%CD%\packages

SET FLAGS="-w -s -H=windowsgui" 

REM download all packages that project needs
go mod tidy

ECHO Building app...

rem build with console hide showing
go build -ldflags %FLAGS% -o %BUILD_PATH%.exe

ECHO App was built successfully!

rem start in dev mode (default)
@REM %BUILD_PATH%.exe