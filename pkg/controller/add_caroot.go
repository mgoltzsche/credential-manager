package controller

import (
	"github.com/mgoltzsche/image-registry-operator/pkg/controller/caroot"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, caroot.Add)
}
