@ECHO OFF
SETLOCAL
chcp 866>nul
CLS

CD ..

SET BUILD_PATH=.\release\pomadorik
SET GO111MODULE=auto
SET GOMODCACHE=%CD%\packages
SET FLAGS="-w -s -H=windowsgui" 

rm -rf release\

REM download all packages that project needs
go mod tidy

ECHO Building release...

rem build with console hide showing
go build -ldflags %FLAGS% -o %BUILD_PATH%.exe

mkdir release\icon
mkdir release\sounds

COPY "icon\app-icon.png" "release\icon\app-icon.png"
COPY "icon\app-icon.ico" "release\icon\app-icon.ico"
COPY "sounds\timer.mp3" "release\sounds\timer.mp3"

ECHO Release app was built successfully!
