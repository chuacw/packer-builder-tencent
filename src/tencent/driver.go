package tencent

import (
	"time"
)

type Driver interface {
	CreateImage(config Config) (bool, CVMError, CVMInstanceInfo)
	DeleteImage(imageId string) error
	WaitForImageCreation(imageId string, timeout time.Duration) error
	WaitForImageDeletion(imageId string, timeout time.Duration) error
	WaitForImageState(imageId string, desiredState string, timeout time.Duration) error
}
