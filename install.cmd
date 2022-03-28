REM @ECHO OFF
SETLOCAL
SET GOPATH=D:\Development\Packer\tencent

SET GOOS=windows
SET GOARCH=amd64

go install tencent

ENDLOCAL
