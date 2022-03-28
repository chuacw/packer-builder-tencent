REM @ECHO OFF
SETLOCAL
SET GOPATH=D:\Development\Packer\packer-builder-tencent

REM go build tencent
REM go install tencent
REM go build ./...
SET DESTDIR=K:\Development\Packer\
SET BASENAME=packer-builder-tencent

del bin\%BASENAME%.exe
del bin\%BASENAME%.linux
del %DESTDIR%\%BASENAME%.exe
del %DESTDIR%\%BASENAME%.linux

SET GOOS=windows
go clean
go build -o bin/%BASENAME%.exe ./cmd/clone
MOVE /Y bin\%BASENAME%.exe %DESTDIR%

SET GOOS=linux
go clean
go build -o bin/%BASENAME%.linux ./cmd/clone
MOVE /Y bin\%BASENAME%.linux %DESTDIR%


ENDLOCAL
