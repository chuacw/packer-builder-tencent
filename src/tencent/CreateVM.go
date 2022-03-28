package tencent

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

// StepCreateImage creates an image with the specified attributes

type StepCreateImage struct{}
type StepDeleteImage struct{}

func (s *StepCreateImage) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)

	if config.PackerDebug || CloudAPIDebug {
		log.Println("Creating image...")
	}

	err, cvmError, instanceInfo := driver.CreateImage(*config)
	if err {
		state.Put("error", fmt.Errorf("Problem creating image: %s", cvmError.Message))
		return multistep.ActionHalt
	}

	ui.Say(fmt.Sprintf("Instance ID: %s", instanceInfo.InstanceId))

	// ui.Say("Waiting for image to become available...")
	// err = driver.WaitForImageCreation(imageId, 10*time.Minute)
	// if err != nil {
	// 	state.Put("error", fmt.Errorf("Problem waiting for image to become available: %s", err))
	// 	return multistep.ActionHalt
	// }

	state.Put("image", instanceInfo.InstanceId)

	return multistep.ActionContinue
}

func (s *StepCreateImage) Cleanup(state multistep.StateBag) {
	// No cleanup
}

func (s *StepDeleteImage) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(Config)
	driver := state.Get("driver").(Driver)
	ui := state.Get("ui").(packer.Ui)

	if config.InstanceName != "" {

	}
	if driver != nil {

	}

	ui.Say("Deleting image...")
	return multistep.ActionHalt

	ui.Say("This is where we begin!!!")

	// imageId, err := driver.CreateImage(config)
	// if err != nil {
	// 	state.Put("error", fmt.Errorf("Problem creating image from machine: %s", err))
	// 	return multistep.ActionHalt
	// }

	// ui.Say("Waiting for image to become available...")
	// err = driver.WaitForImageCreation(imageId, 10*time.Minute)
	// if err != nil {
	// 	state.Put("error", fmt.Errorf("Problem waiting for image to become available: %s", err))
	// 	return multistep.ActionHalt
	// }

	// state.Put("image", imageId)

	return multistep.ActionContinue
}

func (s *StepDeleteImage) Cleanup(state multistep.StateBag) {
	// No cleanup
}
