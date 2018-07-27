package core

// Profile is a given runtime profile
type Profile struct {
	ID           string
	Orchestrator string
	Secrets      string
}

// DefaultProfile returns a  default profile
func DefaultProfile() *Profile {
	return &Profile{
		ID:           "default",
		Orchestrator: "docker",
	}
}
