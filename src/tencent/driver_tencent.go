package tencent

import (
	"errors"
	"log"
	"time"

	"github.com/hashicorp/packer/packer"
)

type TencentDriver struct {
	Ui     packer.Ui
	config Config
}

func (driver TencentDriver) CreateImage(config1 Config) (bool, CVMError, CVMInstanceInfo) {
	if config1.PackerDebug || CloudAPIDebug {
		log.Printf("Creating image: %+v\n", config1)
	}
	// driver.Ui.Say(fmt.Sprintf("Creating image: %+v", config1))
	if config1.PackerDebug || CloudAPIDebug {
		log.Println("Calling CreateVM")
	}
	cvmError, instanceInfo := CreateVM(&config1)
	if config1.PackerDebug || CloudAPIDebug {
		log.Println("Succeeded in calling CreateVM")
	}
	var err bool = false
	if cvmError.Code != "" {
		if config1.PackerDebug || CloudAPIDebug {
			log.Printf("CreateVM %+v", cvmError)
		}
		err = true
	}
	return err, cvmError, instanceInfo
}
func (driver TencentDriver) DeleteImage(imageId string) error {
	log.Printf("deleting image: %s\n", imageId)
	return errors.New("Deleting image error")
}
func (driver TencentDriver) WaitForImageCreation(imageId string, timeout time.Duration) error {
	log.Printf("Wait for image creation: %s\n", imageId)
	return errors.New("Wait for image creation error")
}
func (driver TencentDriver) WaitForImageDeletion(imageId string, timeout time.Duration) error {
	log.Printf("Wait for image deletion: %s\n", imageId)
	return errors.New("Image deletion error")
}

func (driver TencentDriver) WaitForImageState(imageId string, state string, timeout time.Duration) error {
	log.Printf("Wait For image state: %s\n", imageId)
	return errors.New("Wait for Image State error")
}

func NewTencentDriver(ui packer.Ui, config *Config) *TencentDriver {
	driver := new(TencentDriver)
	driver.Ui = ui
	driver.config = *config
	return driver
}
