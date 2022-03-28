REM @ECHO OFF
SETLOCAL
SET GOPATH=D:\Development\Packer\packer-builder-tencent

SET DESTDIR=D:\Development\Packer\
SET BASENAME=packer-builder-tencent

SET GOOS=windows
go clean
go build -o ./hello.exe .


ENDLOCAL
