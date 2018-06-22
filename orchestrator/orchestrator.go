package orchestrator

import (
	"fmt"

	"github.com/euforia/thrap/thrapb"
)

// Config holds the config used to init the orchestrator
type Config struct {
	Provider string
	Conf     map[string]interface{}
}

// Orchestrator implements an application/project deployment orchestrator
type Orchestrator interface {
	// Init is called to initialize the orchestrator with the given config
	Init(config map[string]interface{}) error

	// Deploy should deploy the stack returning the response, deploy object
	// based on the orchestrator or an error
	Deploy(stack *thrapb.Stack, opts DeployOptions) (interface{}, interface{}, error)
}

// New returns a new orchestrator based on the given config
func New(conf *Config) (Orchestrator, error) {
	var (
		orch Orchestrator
		err  error
	)

	switch conf.Provider {
	case "nomad":
		orch = &nomadOrchestrator{}

	default:
		err = fmt.Errorf("unsupported orchestrator: '%s'", conf.Provider)

	}

	if err == nil {
		err = orch.Init(conf.Conf)
	}

	return orch, err
}
