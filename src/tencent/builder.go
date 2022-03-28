package tencent

import (
	"context"
	"errors"
	"log"

	//  "github.com/hashicorp/packer/common"
	//  "github.com/hashicorp/packer/helper/communicator"

	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	//  "github.com/hashicorp/packer/packer/plugin"
)

const (
	// BuilderID to identify this plugin
	BuilderID = "Eximchain.tencent"
)

// Builder structure for the Builder plugin
type Builder struct {
	config  *Config
	runner  multistep.Runner
	context context.Context
	cancel  context.CancelFunc
}

// NewBuilder creates a new instance of the builder
func NewBuilder() *Builder {
	ctx, cancel := context.WithCancel(context.Background())
	return &Builder{
		context: ctx,
		cancel:  cancel,
	}
}

// Prepare decodes the configuration file
func (b *Builder) Prepare(raws ...interface{}) ([]string, error) {

	// errs := &multierror.Error{}

	c, warnings, errs := NewConfig(raws...) // calls config.go's NewConfig
	if errs != nil {
		return warnings, errs
	}

	b.config = c

	if b.config.PackerDebug {
		log.Printf("Config: %+v\n", c)
	}

	// err := config.Decode(&b.config, &config.DecodeOpts{
	// 	Interpolate:        true,
	// 	InterpolateContext: &b.config.ctx,
	// }, raws...)

	// if err != nil {
	// 	errs = multierror.Append(errs, err)
	// }

	// if b.config.PackerDebug {
	// 	log.Println("***** Debug enabled ***")
	// }

	return nil, nil
}

// Run runs the Builder plugin
func (b *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	config := b.config

	if config.PackerDebug || CloudAPIDebug {
		log.Println("In Run")
		log.Printf("Config: %+v\n", config)
	}

	driver := NewTencentDriver(ui, config)

	state := new(multistep.BasicStateBag)
	state.Put("config", b.config)
	state.Put("debug", b.config.PackerDebug)
	state.Put("driver", driver)
	// state.Put("image", "SupposedToBeImageId") // return actual ID here
	state.Put("hook", hook)
	state.Put("ui", ui)

	// Create the VM, and verify the VM is running

	steps := []multistep.Step{
		//   &StepCreateSourceMachine{},
		//   &communicator.StepConnect{
		//   Config: &config.Comm,
		//   Host:   commHost,
		//   SSHConfig: sshConfig(
		//     b.config.Comm.SSHAgentAuth,
		//     b.config.Comm.SSHUsername,
		//     b.config.Comm.SSHPrivateKey,
		//     b.config.Comm.SSHPassword),
		//   },
		//   &common.StepProvision{},
		//   &StepStopMachine{},
		&StepCreateImage{},
	}

	b.runner = common.NewRunnerWithPauseFn(steps, b.config.PackerConfig, ui, state)
	b.runner.Run(state)

	// // If there was an error, return that
	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}

	// // If there is no image, just return
	// if _, ok := state.GetOk("image"); !ok {
	//   return nil, nil
	// }
	// If we were interrupted or cancelled, then just exit.
	if _, ok := state.GetOk(multistep.StateCancelled); ok {
		return nil, errors.New("Build was cancelled.")
	}

	if _, ok := state.GetOk(multistep.StateHalted); ok {
		return nil, errors.New("Build was halted.")
	}

	if config.PackerDebug || CloudAPIDebug {
		log.Println("Generating artifact")
	}

	artifact := &Artifact{
		ImageID:        state.Get("image").(string),
		BuilderIDValue: BuilderID}

	if config.PackerDebug || CloudAPIDebug {
		log.Println("Artifact generated...")
	}

	return artifact, nil

}

// Cancel cancels a possibly running Builder. This should block until
// the builder actually cancels and cleans up after itself.
func (b *Builder) Cancel() {
	if b.runner != nil {
		log.Println("Cancelling the step runner...")
		b.runner.Cancel()
	}
}
