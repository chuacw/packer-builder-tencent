package tencent

import (
	"fmt"
	"log"
)

// Artifact is an artifact implementation that contains built Tencent images.
type Artifact struct {
	// ImageID is the image ID of the artifact
	ImageID string

	// BuilderIDValue is the unique ID for the builder that created this Image
	BuilderIDValue string
}

// BuilderId must be returned
func (a *Artifact) BuilderId() string {
	return a.BuilderIDValue
}

// Files returns the files used/required by the artifact
func (*Artifact) Files() []string {
	// alternative implementation - return []string{}
	return nil
}

// Id returns the ID of the artifact
func (a *Artifact) Id() string {
	return a.ImageID
}

// String returns the ID of the artifact
func (a *Artifact) String() string {
	//TODO(chuacw): Return the proper image id string
	return fmt.Sprintf("Image was created: %s", a.ImageID)
}

// State returns the state specified in the artifact
func (a *Artifact) State(name string) interface{} {
	//TODO(chuacw): Figure out how to make this work
	return nil
}

// Destroy destroys the image specified by the artifact
func (a *Artifact) Destroy() error {
	log.Printf("Deleting image ID (%s)", a.ImageID)

	// requires Driver to be declared in the Artifact struct
	//	err := a.Driver.DeleteImage(a.ImageID)
	//	if err != nil {
	//		return err
	//	}

	return nil
}
