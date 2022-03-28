REM @ECHO OFF
SETLOCAL
SET GOPATH=K:\Development\Packer\packer-builder-tencent

go get github.com/hashicorp/go-multierror
go get github.com/hashicorp/packer/common
go get github.com/hashicorp/packer/helper/communicator
go get github.com/hashicorp/packer/helper/config
go get github.com/hashicorp/packer/helper/multistep
go get github.com/hashicorp/packer/packer
go get github.com/mitchellh/gox
go get github.com/mitchellh/multistep
go get golang.org/x/crypto/ssh
go get golang.org/x/sys/windows

ENDLOCAL
